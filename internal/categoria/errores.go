package categoria

import "errors"

var (
	ErrNombreVacio  = errors.New("el campo nombre es obligatorio")
	ErrNoEncontrado = errors.New("categoria no encontrada")
)
