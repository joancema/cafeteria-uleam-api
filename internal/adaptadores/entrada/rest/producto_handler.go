// Package rest contiene los adaptadores de ENTRADA por HTTP. Dependen del puerto
// de entrada del núcleo (producto.Servicio), nunca del struct concreto.
package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"cafeteria-hex/internal/core/producto"
)

// ProductoHandler adapta HTTP al puerto de entrada del núcleo.
type ProductoHandler struct {
	svc producto.Servicio // <- interfaz (puerto), no el struct
}

func NuevoProductoHandler(svc producto.Servicio) *ProductoHandler {
	return &ProductoHandler{svc: svc}
}

func (h *ProductoHandler) Rutas(r chi.Router) {
	r.Get("/", h.listar)
	r.Post("/", h.crear)
	r.Get("/{id}", h.obtener)
	r.Put("/{id}", h.actualizar)
	r.Delete("/{id}", h.borrar)
}

func (h *ProductoHandler) listar(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, h.svc.Listar())
}

func (h *ProductoHandler) obtener(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	p, err := h.svc.Obtener(id)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusOK, p)
}

func (h *ProductoHandler) crear(w http.ResponseWriter, r *http.Request) {
	var nuevo producto.Producto
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creado, err := h.svc.Crear(nuevo)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, creado)
}

func (h *ProductoHandler) actualizar(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	var datos producto.Producto
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	actualizado, err := h.svc.Actualizar(id, datos)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusOK, actualizado)
}

func (h *ProductoHandler) borrar(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusNoContent, nil)
}

// --- helpers HTTP locales del adaptador ---

func idDeURL(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError traduce los errores del núcleo a códigos HTTP.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, producto.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, producto.ErrNombreVacio), errors.Is(err, producto.ErrPrecioNegativo):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
