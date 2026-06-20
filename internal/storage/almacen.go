package storage

import "cafeteria-uleam-api/internal/models"

// ProductoRepository es el contrato de persistencia SOLO de productos.
type ProductoRepository interface {
	ListarProductos() []models.Producto
	BuscarProductoPorID(id int) (models.Producto, bool)
	CrearProducto(p models.Producto) models.Producto
	ActualizarProducto(id int, datos models.Producto) (models.Producto, bool)
	BorrarProducto(id int) bool
}

// CategoriaRepository es el contrato de persistencia SOLO de categorias.
type CategoriaRepository interface {
	ListarCategorias() []models.Categoria
	BuscarCategoriaPorID(id int) (models.Categoria, bool)
	CrearCategoria(c models.Categoria) models.Categoria
	ActualizarCategoria(id int, datos models.Categoria) (models.Categoria, bool)
	BorrarCategoria(id int) bool
}

// Almacen ahora es la COMPOSICION de las interfaces por entidad (embedding).
//
// Memoria, AlmacenSQLite y AlmacenSQLC ya tienen los 10 metodos, asi que siguen
// cumpliendo Almacen SIN cambiar una sola linea. La diferencia clave: ahora cada
// servicio puede depender solo de la interfaz estrecha que necesita (ISP), y un
// mock para tests implementa 5 metodos en vez de 10.
type Almacen interface {
	ProductoRepository
	CategoriaRepository
}

// UserRepository es el contrato de persistencia de usuarios.
//
// Ojo: CrearUsuario devuelve error (no comma-ok como el resto). El email
// duplicado es un fallo real que el AuthService debe poder mapear a 409; aqui
// la convencion comma-ok ya no alcanza. Es justo el punto que conecta con los
// errores de dominio de S12.
type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}

// Chequeo en tiempo de compilación: si Memoria dejara de cumplir Almacen,
// el proyecto NO compila. Red de seguridad opcional.
var _ Almacen = (*Memoria)(nil)
