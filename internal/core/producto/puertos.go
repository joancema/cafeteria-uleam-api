package producto

// Repositorio es el PUERTO DE SALIDA: lo que el núcleo necesita de la
// persistencia. Lo DEFINE el núcleo; lo IMPLEMENTA un adaptador (GORM, memoria...).
// Esta es la inversión de dependencias clave del hexagonal: la interfaz vive en
// el dominio, no en el paquete de storage.
type Repositorio interface {
	Listar() []Producto
	BuscarPorID(id int) (Producto, bool)
	Crear(p Producto) Producto
	Actualizar(id int, datos Producto) (Producto, bool)
	Borrar(id int) bool
}

// Servicio es el PUERTO DE ENTRADA: lo que el mundo exterior (un handler HTTP, un
// CLI, un test) puede pedirle al núcleo. El adaptador de entrada depende de esta
// interfaz, no del struct concreto que la implementa.
type Servicio interface {
	Listar() []Producto
	Obtener(id int) (Producto, error)
	Crear(p Producto) (Producto, error)
	Actualizar(id int, datos Producto) (Producto, error)
	Borrar(id int) error
}
