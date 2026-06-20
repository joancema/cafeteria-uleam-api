// Package storage gestiona el almacenamiento en memoria de la cafetería.
//
// El tipo Memoria mantiene en un solo lugar todos los datos del dominio:
// Productos y Categorías.
package storage

import (
	"sync"

	"cafeteria-uleam-api/internal/models"
)

// Memoria es un almacén unificado de la cafetería.
type Memoria struct {
	productos     []models.Producto
	nextProductID int

	categorias      []models.Categoria
	nextCategoriaID int

	mu sync.Mutex
}

// NuevaMemoria crea un almacén vacío y listo para usar.
func NuevaMemoria() *Memoria {
	return &Memoria{
		productos:       []models.Producto{},
		nextProductID:   1,
		categorias:      []models.Categoria{},
		nextCategoriaID: 1,
	}
}

// =========================================================
// PRODUCTOS
// =========================================================

// SeedProductos carga productos iniciales en memoria.
func (m *Memoria) SeedProductos() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.productos = []models.Producto{
		{ID: 1, Nombre: "Café americano", Precio: 1.25, Stock: 50, CategoriaID: 1},
		{ID: 2, Nombre: "Capuccino", Precio: 1.75, Stock: 40, CategoriaID: 1},
		{ID: 3, Nombre: "Sandwich vegetariano", Precio: 2.50, Stock: 20, CategoriaID: 2},
		{ID: 4, Nombre: "Croissant de jamón", Precio: 1.80, Stock: 25, CategoriaID: 2},
		{ID: 5, Nombre: "Jugo de naranja", Precio: 1.50, Stock: 30, CategoriaID: 3},
		{ID: 6, Nombre: "Galleta de avena", Precio: 0.75, Stock: 60, CategoriaID: 4},
	}
	m.nextProductID = 7
}

// ListarProductos devuelve todos los productos en memoria.
func (m *Memoria) ListarProductos() []models.Producto {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.Producto, len(m.productos))
	copy(copia, m.productos)
	return copia
}

// BuscarProductoPorID devuelve el producto con el ID dado (patrón comma-ok).
func (m *Memoria) BuscarProductoPorID(id int) (models.Producto, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.productos {
		if p.ID == id {
			return p, true
		}
	}
	return models.Producto{}, false
}

// CrearProducto agrega un producto nuevo y devuelve el producto con ID asignado.
func (m *Memoria) CrearProducto(p models.Producto) models.Producto {
	m.mu.Lock()
	defer m.mu.Unlock()

	p.ID = m.nextProductID
	m.nextProductID++
	m.productos = append(m.productos, p)
	return p
}

// ActualizarProducto reemplaza el producto con el ID dado.
func (m *Memoria) ActualizarProducto(id int, datos models.Producto) (models.Producto, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.productos {
		if p.ID == id {
			datos.ID = id
			m.productos[i] = datos
			return datos, true
		}
	}
	return models.Producto{}, false
}

// BorrarProducto elimina el producto con el ID dado.
func (m *Memoria) BorrarProducto(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.productos {
		if p.ID == id {
			m.productos = append(m.productos[:i], m.productos[i+1:]...)
			return true
		}
	}
	return false
}

// =========================================================
// CATEGORIAS
// =========================================================

// SeedCategorias carga categorías iniciales que coinciden con CategoriaID de los productos pre-cargados.
func (m *Memoria) SeedCategorias() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.categorias = []models.Categoria{
		{ID: 1, Nombre: "Bebidas calientes", Descripcion: "Cafés, tés e infusiones servidos calientes"},
		{ID: 2, Nombre: "Alimentos sólidos", Descripcion: "Sandwiches, panes y comida lista para llevar"},
		{ID: 3, Nombre: "Bebidas frías", Descripcion: "Jugos, refrescos y bebidas frías en general"},
		{ID: 4, Nombre: "Snacks y galletas", Descripcion: "Productos de panadería y snacks empacados"},
	}
	m.nextCategoriaID = 5
}

// ListarCategorias devuelve todas las categorías en memoria.
func (m *Memoria) ListarCategorias() []models.Categoria {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.Categoria, len(m.categorias))
	copy(copia, m.categorias)
	return copia
}

// BuscarCategoriaPorID devuelve la categoría con el ID dado (patrón comma-ok).
func (m *Memoria) BuscarCategoriaPorID(id int) (models.Categoria, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.categorias {
		if c.ID == id {
			return c, true
		}
	}
	return models.Categoria{}, false
}

// CrearCategoria agrega una categoría nueva y devuelve la categoría con ID asignado.
func (m *Memoria) CrearCategoria(c models.Categoria) models.Categoria {
	m.mu.Lock()
	defer m.mu.Unlock()

	c.ID = m.nextCategoriaID
	m.nextCategoriaID++
	m.categorias = append(m.categorias, c)
	return c
}

// ActualizarCategoria reemplaza la categoría con el ID dado.
func (m *Memoria) ActualizarCategoria(id int, datos models.Categoria) (models.Categoria, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, c := range m.categorias {
		if c.ID == id {
			datos.ID = id
			m.categorias[i] = datos
			return datos, true
		}
	}
	return models.Categoria{}, false
}

// BorrarCategoria elimina la categoría con el ID dado.
func (m *Memoria) BorrarCategoria(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, c := range m.categorias {
		if c.ID == id {
			m.categorias = append(m.categorias[:i], m.categorias[i+1:]...)
			return true
		}
	}
	return false
}
