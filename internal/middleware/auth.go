package middleware

import (
	"context"
	"net/http"
	"strings"

	"cafeteria-uleam-api/internal/service"
)

// claveContexto es un tipo privado para la clave del context y evitar colisiones.
type claveContexto string

// ClaveUsuarioID es la clave bajo la que se guarda el ID del usuario autenticado.
const ClaveUsuarioID claveContexto = "usuarioID"

// Auth construye un middleware que exige un JWT valido en el header
// Authorization: Bearer <token>. Delega la validacion al AuthService: el
// middleware NO sabe de firmas ni de claims, solo de HTTP.
func Auth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(siguiente http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encabezado := r.Header.Get("Authorization")
			partes := strings.SplitN(encabezado, " ", 2)
			if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
				responderNoAutorizado(w)
				return
			}

			usuarioID, err := auth.ValidarToken(partes[1])
			if err != nil {
				responderNoAutorizado(w)
				return
			}

			ctx := context.WithValue(r.Context(), ClaveUsuarioID, usuarioID)
			siguiente.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"token ausente o invalido"}`))
}
