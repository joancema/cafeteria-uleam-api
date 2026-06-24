package auth

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

type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Registrar(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		web.RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	u, err := h.svc.Registrar(c.Email, c.Password)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusCreated, u)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		web.RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	token, err := h.svc.Login(c.Email, c.Password)
	if err != nil {
		web.RespondError(w, statusDeError(err), err.Error())
		return
	}
	web.RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}

func statusDeError(err error) int {
	switch {
	case errors.Is(err, ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
