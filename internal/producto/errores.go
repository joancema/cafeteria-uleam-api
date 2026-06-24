package producto

import "errors"

// Errores de dominio de la vertical de productos. Cada slice define los suyos:
// el handler de este mismo paquete los traduce a HTTP (ver handler.go).
var (
	ErrNombreVacio    = errors.New("el campo nombre es obligatorio")
	ErrPrecioNegativo = errors.New("el precio no puede ser negativo")
	ErrNoEncontrado   = errors.New("producto no encontrado")
)
