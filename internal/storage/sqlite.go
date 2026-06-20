package storage

import (
	"gorm.io/gorm"

	"cafeteria-uleam-api/internal/models"
)

// AlmacenSQLite implementa la interfaz Almacen usando GORM sobre SQLite.
//
// Fíjense: los métodos tienen EXACTAMENTE las mismas firmas que los de Memoria.
// Por eso el Server y los handlers no se enteran de cuál de los dos reciben.
type AlmacenSQLite struct {
	db *gorm.DB
}

// NuevoAlmacenSQLite envuelve una conexión *gorm.DB ya abierta.
func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

// =========================================================
// PRODUCTOS
// =========================================================

func (a *AlmacenSQLite) ListarProductos() []models.Producto {
	var productos []models.Producto
	a.db.Find(&productos)
	return productos
}

func (a *AlmacenSQLite) BuscarProductoPorID(id int) (models.Producto, bool) {
	var p models.Producto
	if err := a.db.First(&p, id).Error; err != nil {
		// Absorbemos el error de la DB y conservamos la firma comma-ok.
		return models.Producto{}, false
	}
	return p, true
}

func (a *AlmacenSQLite) CrearProducto(p models.Producto) models.Producto {
	a.db.Create(&p) // GORM rellena el ID autogenerado en &p
	return p
}

func (a *AlmacenSQLite) ActualizarProducto(id int, datos models.Producto) (models.Producto, bool) {
	var existente models.Producto
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Producto{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarProducto(id int) bool {
	res := a.db.Delete(&models.Producto{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// CATEGORIAS
// =========================================================

func (a *AlmacenSQLite) ListarCategorias() []models.Categoria {
	var categorias []models.Categoria
	a.db.Find(&categorias)
	return categorias
}

func (a *AlmacenSQLite) BuscarCategoriaPorID(id int) (models.Categoria, bool) {
	var c models.Categoria
	if err := a.db.First(&c, id).Error; err != nil {
		return models.Categoria{}, false
	}
	return c, true
}

func (a *AlmacenSQLite) CrearCategoria(c models.Categoria) models.Categoria {
	a.db.Create(&c)
	return c
}

func (a *AlmacenSQLite) ActualizarCategoria(id int, datos models.Categoria) (models.Categoria, bool) {
	var existente models.Categoria
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Categoria{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarCategoria(id int) bool {
	res := a.db.Delete(&models.Categoria{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// SEEDS
// =========================================================

// SembrarSiVacio inserta datos iniciales solo si aún no hay categorías.
// Así no duplicamos datos en cada arranque del servidor.
func (a *AlmacenSQLite) SembrarSiVacio() {
	var n int64
	a.db.Model(&models.Categoria{}).Count(&n)
	if n > 0 {
		return
	}

	categorias := []models.Categoria{
		{ID: 1, Nombre: "Bebidas calientes", Descripcion: "Cafés, tés e infusiones servidos calientes"},
		{ID: 2, Nombre: "Alimentos sólidos", Descripcion: "Sandwiches, panes y comida lista para llevar"},
		{ID: 3, Nombre: "Bebidas frías", Descripcion: "Jugos, refrescos y bebidas frías en general"},
		{ID: 4, Nombre: "Snacks y galletas", Descripcion: "Productos de panadería y snacks empacados"},
	}
	a.db.Create(&categorias)

	productos := []models.Producto{
		{ID: 1, Nombre: "Café americano", Precio: 1.25, Stock: 50, CategoriaID: 1},
		{ID: 2, Nombre: "Capuccino", Precio: 1.75, Stock: 40, CategoriaID: 1},
		{ID: 3, Nombre: "Sandwich vegetariano", Precio: 2.50, Stock: 20, CategoriaID: 2},
		{ID: 4, Nombre: "Croissant de jamón", Precio: 1.80, Stock: 25, CategoriaID: 2},
		{ID: 5, Nombre: "Jugo de naranja", Precio: 1.50, Stock: 30, CategoriaID: 3},
		{ID: 6, Nombre: "Galleta de avena", Precio: 0.75, Stock: 60, CategoriaID: 4},
	}
	a.db.Create(&productos)
}

// Chequeo en tiempo de compilación: AlmacenSQLite debe cumplir Almacen.
var _ Almacen = (*AlmacenSQLite)(nil)
