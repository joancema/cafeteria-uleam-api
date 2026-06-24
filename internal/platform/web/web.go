// Package web reune helpers HTTP transversales: respuestas JSON y extraccion de
// parametros. NO conoce el dominio (no importa producto/categoria/auth), por eso
// los slices pueden depender de el sin crear ciclos.
package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// RespondJSON escribe data como JSON con el status indicado.
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

// RespondError escribe {"error": "..."} con el status indicado.
func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// IDDeURL extrae y valida el parametro de ruta {id} como entero.
func IDDeURL(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}
