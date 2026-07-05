package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/models"
)

// ejecutar corre una peticion contra el handler y devuelve el recorder.
func ejecutar(h http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestListarProductos_OK(t *testing.T) {
	h, _, _ := construirEntorno()
	token := tokenValido(t, h)

	rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/productos", "", token))

	require.Equal(t, http.StatusOK, rec.Code)
	var lista []models.Producto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&lista))
	assert.Len(t, lista, 1) // el producto sembrado
}

func TestObtenerProducto(t *testing.T) {
	h, _, _ := construirEntorno()
	token := tokenValido(t, h)

	t.Run("existe -> 200", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/productos/1", "", token))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("no existe -> 404", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/productos/9999", "", token))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
	t.Run("id no numerico -> 400", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/productos/abc", "", token))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestCrearProducto(t *testing.T) {
	h, _, _ := construirEntorno()
	token := tokenValido(t, h)

	t.Run("valido -> 201", func(t *testing.T) {
		body := `{"nombre":"Te verde","precio":1.10,"stock":15,"categoria_id":1}`
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/productos", body, token))
		// require.Equal(t, http.StatusCreated, rec.Code)
		require.Equal(t, http.StatusTeapot, rec.Code)
		var creado models.Producto
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
		assert.NotZero(t, creado.ID)
	})
	t.Run("nombre vacio -> 400", func(t *testing.T) {
		body := `{"nombre":"   ","precio":2.0}`
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/productos", body, token))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("JSON malformado -> 400", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/productos", `{roto`, token))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestActualizarProducto(t *testing.T) {
	h, _, _ := construirEntorno()
	token := tokenValido(t, h)

	t.Run("valido -> 200", func(t *testing.T) {
		body := `{"nombre":"Cafe doble","precio":1.80,"stock":5,"categoria_id":1}`
		rec := ejecutar(h, jsonReq(http.MethodPut, "/api/v1/productos/1", body, token))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("no existe -> 404", func(t *testing.T) {
		body := `{"nombre":"X","precio":1.0}`
		rec := ejecutar(h, jsonReq(http.MethodPut, "/api/v1/productos/9999", body, token))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestBorrarProducto(t *testing.T) {
	h, _, _ := construirEntorno()
	token := tokenValido(t, h)

	t.Run("existe -> 204", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodDelete, "/api/v1/productos/1", "", token))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
	t.Run("no existe -> 404", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodDelete, "/api/v1/productos/9999", "", token))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// El corazon de la seguridad: el middleware corta ANTES del handler.
func TestRutaProtegida_SinToken(t *testing.T) {
	h, _, _ := construirEntorno()
	rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/productos", "", "")) // sin Bearer
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRutaProtegida_TokenInvalido(t *testing.T) {
	h, _, _ := construirEntorno()
	rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/productos", "", "token.falso.123"))
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
