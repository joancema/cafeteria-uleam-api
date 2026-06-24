package middleware

import (
	"context"
	"net/http"
	"strings"
)

type claveContexto string

const ClaveUsuarioID claveContexto = "usuarioID"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Validador es el contrato que el middleware necesita: lo cumple la FACHADA
// pública del módulo auth (su método ValidarToken). El middleware no importa auth.
type Validador interface {
	ValidarToken(token string) (int, error)
}

func Auth(v Validador) func(http.Handler) http.Handler {
	return func(siguiente http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			partes := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
				noAutorizado(w)
				return
			}
			uid, err := v.ValidarToken(partes[1])
			if err != nil {
				noAutorizado(w)
				return
			}
			ctx := context.WithValue(r.Context(), ClaveUsuarioID, uid)
			siguiente.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func noAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"token ausente o invalido"}`))
}
