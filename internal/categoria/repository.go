package categoria

type Repository interface {
	Listar() []Categoria
	BuscarPorID(id int) (Categoria, bool)
	Crear(c Categoria) Categoria
	Actualizar(id int, datos Categoria) (Categoria, bool)
	Borrar(id int) bool
}
