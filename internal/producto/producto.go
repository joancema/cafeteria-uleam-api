// Package producto agrupa TODA la vertical de productos: modelo, contrato de
// persistencia, adaptador GORM, logica de negocio, handler HTTP y rutas.
// Esa es la esencia del Vertical Slice: una sola carpeta por entidad, en lugar
// de repartir productos entre handlers/, service/ y storage/.
package producto

// Producto es el modelo de dominio. CategoriaID referencia una Categoria por id.
type Producto struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Nombre      string  `json:"nombre" gorm:"not null"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	CategoriaID int     `json:"categoria_id"`
}
