// Package handlers contiene los handlers HTTP de la API de cafeteria.
package handlers

import "cafeteria-uleam-api/internal/service"

// Server agrupa los servicios de los que dependen los handlers.
//
// Antes guardaba un storage.Almacen directo; ahora guarda la capa de servicio,
// que es la que tiene la logica de negocio. Los handlers quedan delgados:
// decodifican el request, llaman al servicio y traducen el resultado a HTTP.
type Server struct {
	Productos  *service.ProductoService
	Categorias *service.CategoriaService
	Auth       *service.AuthService
}

// NewServer construye un Server con sus servicios ya inyectados.
func NewServer(productos *service.ProductoService, categorias *service.CategoriaService, auth *service.AuthService) *Server {
	return &Server{
		Productos:  productos,
		Categorias: categorias,
		Auth:       auth,
	}
}
