// Package producto es el NÚCLEO de dominio de productos. No importa GORM, ni chi,
// ni nada de infraestructura: solo Go estándar. Las dependencias apuntan HACIA
// aquí, nunca al revés (regla de la arquitectura hexagonal).
package producto

import "errors"

type Producto struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Nombre      string  `json:"nombre" gorm:"not null"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	CategoriaID int     `json:"categoria_id"`
}

var (
	ErrNombreVacio    = errors.New("el campo nombre es obligatorio")
	ErrPrecioNegativo = errors.New("el precio no puede ser negativo")
	ErrNoEncontrado   = errors.New("producto no encontrado")
)
