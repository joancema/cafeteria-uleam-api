package handlers

import (
	"encoding/json"
	"net/http"
)

// credenciales es el cuerpo JSON esperado en register y login.
type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Registrar atiende POST /api/v1/auth/register.
func (s *Server) Registrar(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	usuario, err := s.Auth.Registrar(c.Email, c.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, usuario)
}

// Login atiende POST /api/v1/auth/login y devuelve {"token": "..."}.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	token, err := s.Auth.Login(c.Email, c.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
