// Package middleware contiene middlewares propios de la API de cafetería.
package middleware

import "net/http"

// CORS habilita el consumo de la API desde cualquier origen (configuración
// permisiva, para desarrollo y demos en clase). Esta es la pieza que faltaba
// en el proyecto base y que los estudiantes deben re-agregar.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// El navegador manda un preflight OPTIONS antes de la petición real.
		// Lo respondemos de inmediato, sin pasar al handler.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
