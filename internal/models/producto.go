// Package models define las entidades del dominio cafetería.
package models

// Producto representa un producto a la venta en la cafetería.
//
// CategoriaID referencia el ID de una Categoria por número (foreign key).
// Decisión arquitectónica: usamos ID en lugar de struct anidado. GORM usa
// ese mismo ID como clave foránea, así que el modelo casi no cambia.
type Producto struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Nombre      string  `json:"nombre" gorm:"not null"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	CategoriaID int     `json:"categoria_id"`
}
