// Package middleware contiene middlewares transversales.
package middleware

import (
	"context"
	"net/http"
	"strings"
)

type claveContexto string

const ClaveUsuarioID claveContexto = "usuarioID"

// Validador es lo UNICO que el middleware necesita del slice auth: validar un
// token y devolver el id del usuario. Asi el middleware NO importa el paquete
// auth (evita el acoplamiento; cualquier tipo con este metodo sirve).
type Validador interface {
	ValidarToken(token string) (int, error)
}

func Auth(v Validador) func(http.Handler) http.Handler {
	return func(siguiente http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			partes := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
				responderNoAutorizado(w)
				return
			}
			usuarioID, err := v.ValidarToken(partes[1])
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
