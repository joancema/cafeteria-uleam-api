package producto

// Repository es el contrato de persistencia de productos. Vive en el slice (el
// "puerto"); la implementacion concreta (GORM) tambien vive aqui. Los metodos van
// sin sufijo "Producto" porque el paquete ya da el contexto: producto.Repository.
type Repository interface {
	Listar() []Producto
	BuscarPorID(id int) (Producto, bool)
	Crear(p Producto) Producto
	Actualizar(id int, datos Producto) (Producto, bool)
	Borrar(id int) bool
}
