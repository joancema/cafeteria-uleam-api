package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/handlers"
	"cafeteria-uleam-api/internal/middleware"
	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

// =====================================================================
// Dobles de prueba
//
// Un test de la capa HTTP NO debe tocar la base de datos real: la
// reemplazamos por estos dobles en memoria. Son SOLO para los tests
// (viven en archivos _test.go); en produccion corre GORM. La base real
// tiene su propio test en internal/storage/sqlite_test.go.
// =====================================================================

// almacenFake implementa storage.Almacen (productos + categorias) en memoria.
type almacenFake struct {
	productos  map[int]models.Producto
	categorias map[int]models.Categoria
	nextProd   int
	nextCat    int
}

func nuevoAlmacenFake() *almacenFake {
	return &almacenFake{
		productos:  map[int]models.Producto{},
		categorias: map[int]models.Categoria{},
		nextProd:   1, nextCat: 1,
	}
}

func (a *almacenFake) ListarProductos() []models.Producto {
	out := make([]models.Producto, 0, len(a.productos))
	for _, p := range a.productos {
		out = append(out, p)
	}
	return out
}
func (a *almacenFake) BuscarProductoPorID(id int) (models.Producto, bool) {
	p, ok := a.productos[id]
	return p, ok
}
func (a *almacenFake) CrearProducto(p models.Producto) models.Producto {
	p.ID = a.nextProd
	a.nextProd++
	a.productos[p.ID] = p
	return p
}
func (a *almacenFake) ActualizarProducto(id int, datos models.Producto) (models.Producto, bool) {
	if _, ok := a.productos[id]; !ok {
		return models.Producto{}, false
	}
	datos.ID = id
	a.productos[id] = datos
	return datos, true
}
func (a *almacenFake) BorrarProducto(id int) bool {
	if _, ok := a.productos[id]; !ok {
		return false
	}
	delete(a.productos, id)
	return true
}
func (a *almacenFake) ListarCategorias() []models.Categoria {
	out := make([]models.Categoria, 0, len(a.categorias))
	for _, c := range a.categorias {
		out = append(out, c)
	}
	return out
}
func (a *almacenFake) BuscarCategoriaPorID(id int) (models.Categoria, bool) {
	c, ok := a.categorias[id]
	return c, ok
}
func (a *almacenFake) CrearCategoria(c models.Categoria) models.Categoria {
	c.ID = a.nextCat
	a.nextCat++
	a.categorias[c.ID] = c
	return c
}
func (a *almacenFake) ActualizarCategoria(id int, datos models.Categoria) (models.Categoria, bool) {
	if _, ok := a.categorias[id]; !ok {
		return models.Categoria{}, false
	}
	datos.ID = id
	a.categorias[id] = datos
	return datos, true
}
func (a *almacenFake) BorrarCategoria(id int) bool {
	if _, ok := a.categorias[id]; !ok {
		return false
	}
	delete(a.categorias, id)
	return true
}

var _ storage.Almacen = (*almacenFake)(nil)

// usuarioFake implementa storage.UserRepository en memoria.
type usuarioFake struct {
	porEmail map[string]models.Usuario
	nextID   int
}

func nuevoUsuarioFake() *usuarioFake {
	return &usuarioFake{porEmail: map[string]models.Usuario{}, nextID: 1}
}
func (f *usuarioFake) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = f.nextID
	f.nextID++
	f.porEmail[u.Email] = u
	return u, nil
}
func (f *usuarioFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

var _ storage.UserRepository = (*usuarioFake)(nil)

// =====================================================================
// Router de prueba: el MISMO que main.go (mismas rutas, mismo
// middleware.Auth real), pero con los dobles en memoria.
// =====================================================================

// construirEntorno devuelve el handler listo y un producto sembrado (id 1).
func construirEntorno() (http.Handler, *almacenFake, *usuarioFake) {
	almacen := nuevoAlmacenFake()
	almacen.CrearProducto(models.Producto{Nombre: "Cafe americano", Precio: 1.25, Stock: 10, CategoriaID: 1})
	usuarios := nuevoUsuarioFake()

	srv := handlers.NewServer(handlers.Deps{
		Productos:  service.NuevoProductoService(almacen),
		Categorias: service.NuevoCategoriaService(almacen),
		Auth:       service.NuevoAuthService(usuarios),
	})
	authSvc := service.NuevoAuthService(usuarios)

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
	return r, almacen, usuarios
}

// jsonReq arma una peticion con cuerpo JSON y, si se pasa token, el header Bearer.
func jsonReq(metodo, ruta, cuerpo, token string) *http.Request {
	var body *strings.Reader
	if cuerpo == "" {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(cuerpo)
	}
	req := httptest.NewRequest(metodo, ruta, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

// tokenValido registra y loguea un usuario contra el router y devuelve su JWT.
func tokenValido(t *testing.T, h http.Handler) string {
	t.Helper()
	cred := `{"email":"docente@uleam.edu.ec","password":"secreta123"}`
	h.ServeHTTP(httptest.NewRecorder(), jsonReq(http.MethodPost, "/api/v1/auth/register", cred, ""))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, jsonReq(http.MethodPost, "/api/v1/auth/login", cred, ""))
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.NotEmpty(t, resp.Token)
	return resp.Token
}
