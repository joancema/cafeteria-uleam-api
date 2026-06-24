// Package categoria agrupa la vertical completa de categorias.
package categoria

type Categoria struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Nombre      string `json:"nombre" gorm:"not null"`
	Descripcion string `json:"descripcion"`
}
