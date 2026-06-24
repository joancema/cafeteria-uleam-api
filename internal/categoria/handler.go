package categoria

import (
	"encoding/json"
	"errors"
	"net/http"

	"cafeteria-uleam-api/internal/platform/web"
)

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
	c, err := h.svc.Obtener(id)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusOK, c)
}

func (h *Handler) Crear(w http.ResponseWriter, r *http.Request) {
	var nueva Categoria
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		web.RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creada, err := h.svc.Crear(nueva)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusCreated, creada)
}

func (h *Handler) Actualizar(w http.ResponseWriter, r *http.Request) {
	id, err := web.IDDeURL(r)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	var datos Categoria
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		web.RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	actualizada, err := h.svc.Actualizar(id, datos)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusOK, actualizada)
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

func statusDeError(err error) int {
	switch {
	case errors.Is(err, ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, ErrNombreVacio):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
