package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/handlers"
	"cafeteria-uleam-api/internal/middleware"
	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

// usuarioRepoFake: repositorio de usuarios en memoria para los tests de handler.
type usuarioRepoFake struct {
	porEmail map[string]models.Usuario
	nextID   int
}

func nuevoUsuarioRepoFake() *usuarioRepoFake {
	return &usuarioRepoFake{porEmail: map[string]models.Usuario{}, nextID: 1}
}

func (f *usuarioRepoFake) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = f.nextID
	f.nextID++
	f.porEmail[u.Email] = u
	return u, nil
}

func (f *usuarioRepoFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

// construirEntorno arma el MISMO router que main.go (mismas rutas, mismo
// middleware.Auth real) pero con almacen en memoria y repo de usuarios fake.
// Devuelve el handler listo para httptest y un token valido ya emitido.
//
// Clave pedagogica: probamos a traves del middleware REAL, no de uno simplificado.
// Si el wiring de la ruta protegida se rompe, este test se entera.
func construirEntorno(t *testing.T) (http.Handler, string) {
	t.Helper()

	almacen := storage.NuevaMemoria()
	almacen.SeedProductos()
	almacen.SeedCategorias()
	usuarios := nuevoUsuarioRepoFake()

	productoSvc := service.NuevoProductoService(almacen)
	categoriaSvc := service.NuevoCategoriaService(almacen)
	authSvc := service.NuevoAuthService(usuarios)
	srv := handlers.NewServer(productoSvc, categoriaSvc, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc)) // <- el middleware real de S10
			r.Get("/productos", srv.ListarProductos)
			r.Post("/productos", srv.CrearProducto)
			r.Get("/productos/{id}", srv.ObtenerProducto)
			r.Put("/productos/{id}", srv.ActualizarProducto)
			r.Delete("/productos/{id}", srv.BorrarProducto)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// registrarYObtenerToken hace register + login contra el propio router para
// conseguir un JWT valido, igual que lo haria un cliente real.
func registrarYObtenerToken(t *testing.T, h http.Handler) string {
	t.Helper()
	cred := `{"email":"docente@uleam.edu.ec","password":"secreta123"}`

	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(cred))
	h.ServeHTTP(httptest.NewRecorder(), reqReg)

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(cred))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqLogin)
	require.Equal(t, http.StatusOK, rec.Code, "el login deberia devolver 200")

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.NotEmpty(t, resp.Token)
	return resp.Token
}

// TestCrearProducto_Exitoso: POST con token y cuerpo valido -> 201 Created.
func TestCrearProducto_Exitoso(t *testing.T) {
	h, token := construirEntorno(t)
	body := `{"nombre":"Te verde","precio":1.10,"stock":15,"categoria_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/productos", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creado models.Producto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
	assert.NotZero(t, creado.ID)
	assert.Equal(t, "Te verde", creado.Nombre)
}

// TestObtenerProducto_NoEncontrado: id inexistente -> 404 Not Found.
func TestObtenerProducto_NoEncontrado(t *testing.T) {
	h, token := construirEntorno(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/productos/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestCrearProducto_Invalido: cuerpo que viola la regla de negocio -> 400.
func TestCrearProducto_Invalido(t *testing.T) {
	h, token := construirEntorno(t)
	body := `{"nombre":"   ","precio":2.0}` // nombre vacio

	req := httptest.NewRequest(http.MethodPost, "/api/v1/productos", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaProtegida_SinToken: sin header Authorization, el middleware corta
// antes de llegar al handler -> 401 Unauthorized.
func TestRutaProtegida_SinToken(t *testing.T) {
	h, _ := construirEntorno(t)
	body := `{"nombre":"Te verde","precio":1.10}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/productos", strings.NewReader(body))
	// A proposito: NO seteamos Authorization.
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
