package categoria

import "github.com/go-chi/chi/v5"

func (h *Handler) Rutas(r chi.Router) {
	r.Get("/", h.Listar)
	r.Post("/", h.Crear)
	r.Get("/{id}", h.Obtener)
	r.Put("/{id}", h.Actualizar)
	r.Delete("/{id}", h.Borrar)
}
