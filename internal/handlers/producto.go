package handlers

import (
	"encoding/json"
	"net/http"

	"cafeteria-uleam-api/internal/models"
)

// ListarProductos atiende GET /api/v1/productos.
func (s *Server) ListarProductos(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Productos.Listar())
}

// ObtenerProducto atiende GET /api/v1/productos/{id}.
func (s *Server) ObtenerProducto(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	producto, err := s.Productos.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, producto)
}

// CrearProducto atiende POST /api/v1/productos.
func (s *Server) CrearProducto(w http.ResponseWriter, r *http.Request) {
	var nuevo models.Producto
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creado, err := s.Productos.Crear(nuevo)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarProducto atiende PUT /api/v1/productos/{id}.
func (s *Server) ActualizarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	var datos models.Producto
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	actualizado, err := s.Productos.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

// BorrarProducto atiende DELETE /api/v1/productos/{id}.
func (s *Server) BorrarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	if err := s.Productos.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
