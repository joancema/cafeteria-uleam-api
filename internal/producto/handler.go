package producto

import (
	"encoding/json"
	"errors"
	"net/http"

	"cafeteria-uleam-api/internal/platform/web"
)

// Handler traduce HTTP <-> Service para la vertical de productos.
type Handler struct {
	svc *Service
}

func NuevoHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Listar(w http.ResponseWriter, _ *http.Request) {
	web.RespondJSON(w, http.StatusOK, h.svc.Listar())
}

func (h *Handler) Obtener(w http.ResponseWriter, r *http.Request) {
	id, err := web.IDDeURL(r)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	p, err := h.svc.Obtener(id)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusOK, p)
}

func (h *Handler) Crear(w http.ResponseWriter, r *http.Request) {
	var nuevo Producto
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		web.RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creado, err := h.svc.Crear(nuevo)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusCreated, creado)
}

func (h *Handler) Actualizar(w http.ResponseWriter, r *http.Request) {
	id, err := web.IDDeURL(r)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	var datos Producto
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		web.RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	actualizado, err := h.svc.Actualizar(id, datos)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusOK, actualizado)
}

func (h *Handler) Borrar(w http.ResponseWriter, r *http.Request) {
	id, err := web.IDDeURL(r)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusNoContent, nil)
}

// statusDeError mapea los errores de ESTE slice a codigos HTTP.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, ErrNombreVacio), errors.Is(err, ErrPrecioNegativo):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
