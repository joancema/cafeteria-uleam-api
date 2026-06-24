// Package handlers contiene los handlers HTTP de la API de cafeteria.
package handlers

import "cafeteria-uleam-api/internal/service"

// Server agrupa los servicios de los que dependen los handlers.
//
// Guarda la capa de servicio (no el storage directo): los handlers quedan
// delgados: decodifican el request, llaman al servicio y traducen el resultado
// a HTTP.
type Server struct {
	Productos  *service.ProductoService
	Categorias *service.CategoriaService
	Auth       *service.AuthService
}

// Deps agrupa las dependencias requeridas para construir un Server.
//
// Antes NewServer recibia un parametro posicional por servicio; agregar una
// entidad obligaba a cambiar la firma Y todos los call-sites, y dos parametros
// del mismo tipo eran faciles de intercambiar por error. Con un struct de
// dependencias y campos NOMBRADOS, agregar una entidad es agregar un campo:
// nada mas se rompe y desaparece el riesgo de intercambiar argumentos.
type Deps struct {
	Productos  *service.ProductoService
	Categorias *service.CategoriaService
	Auth       *service.AuthService
}

// NewServer construye un Server a partir de sus dependencias.
func NewServer(d Deps) *Server {
	return &Server{
		Productos:  d.Productos,
		Categorias: d.Categorias,
		Auth:       d.Auth,
	}
}
