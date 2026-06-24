package producto

import "github.com/go-chi/chi/v5"

// Rutas monta los endpoints de productos en el router que reciba. Se usa con
// r.Route("/productos", h.Rutas) desde main: cada slice publica SUS rutas.
func (h *Handler) Rutas(r chi.Router) {
	r.Get("/", h.Listar)
	r.Post("/", h.Crear)
	r.Get("/{id}", h.Obtener)
	r.Put("/{id}", h.Actualizar)
	r.Delete("/{id}", h.Borrar)
}
