package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"cafeteria-uleam-api/internal/service"
)

// RespondJSON escribe data como JSON con el status HTTP indicado.
// Si data es nil (por ejemplo en un 204 No Content) no escribe cuerpo.
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

// RespondError escribe un error en formato JSON consistente: {"error": "..."}.
func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError traduce los errores de dominio del service al codigo HTTP que
// corresponde. Es el unico lugar donde el paquete handlers conoce ese mapeo.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrNombreVacio), errors.Is(err, service.ErrPrecioNegativo):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
