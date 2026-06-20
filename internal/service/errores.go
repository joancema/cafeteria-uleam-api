package service

import "errors"

// Errores de dominio. El handler los traduce a codigos HTTP:
//
//	ErrNombreVacio, ErrPrecioNegativo -> 400 Bad Request
//	ErrNoEncontrado                   -> 404 Not Found
//	ErrEmailEnUso                     -> 409 Conflict
//	ErrCredencialesInvalidas          -> 401 Unauthorized
//
// El repositorio sigue en comma-ok; es el service quien reintroduce el error
// con significado de negocio.
var (
	ErrNombreVacio           = errors.New("el campo nombre es obligatorio")
	ErrPrecioNegativo        = errors.New("el precio no puede ser negativo")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrEmailEnUso            = errors.New("el email ya esta registrado")
	ErrCredencialesInvalidas = errors.New("email o contrasena incorrectos")
)
