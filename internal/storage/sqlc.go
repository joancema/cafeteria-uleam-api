package storage

import (
	"context"
	"database/sql"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/storage/sqlcdb"
)

// AlmacenSQLC implementa la interfaz Almacen usando código generado por sqlc
// (SQL escrito a mano + tipado generado) sobre database/sql.
//
// Es el TERCER backend de la cafetería, hermano de Memoria y AlmacenSQLite.
// El Server y los handlers no se enteran de cuál reciben: todos cumplen Almacen.
//
// Diferencias con sqlc que el adaptador tiene que resolver:
//  1. Los métodos generados piden context.Context  -> lo inyectamos acá dentro.
//  2. sqlc devuelve sus propios structs (int64)     -> los MAPEAMOS a models (int).
//  3. sqlc devuelve (T, error)                       -> lo absorbemos a (T, bool).
type AlmacenSQLC struct {
	q *sqlcdb.Queries
}

// NuevoAlmacenSQLC envuelve una conexión *sql.DB ya abierta.
func NuevoAlmacenSQLC(db *sql.DB) *AlmacenSQLC {
	return &AlmacenSQLC{q: sqlcdb.New(db)}
}

// =========================================================
// MAPEO sqlc <-> dominio (la "capa anticorrupción")
// =========================================================

func aProductoDominio(p sqlcdb.Producto) models.Producto {
	return models.Producto{
		ID:          int(p.ID),
		Nombre:      p.Nombre,
		Precio:      p.Precio,
		Stock:       int(p.Stock),
		CategoriaID: int(p.CategoriaID),
	}
}

func aCategoriaDominio(c sqlcdb.Categoria) models.Categoria {
	return models.Categoria{
		ID:          int(c.ID),
		Nombre:      c.Nombre,
		Descripcion: c.Descripcion,
	}
}

// =========================================================
// PRODUCTOS
// =========================================================

func (a *AlmacenSQLC) ListarProductos() []models.Producto {
	filas, err := a.q.ListarProductos(context.Background())
	if err != nil {
		return nil
	}
	out := make([]models.Producto, 0, len(filas))
	for _, f := range filas {
		out = append(out, aProductoDominio(f))
	}
	return out
}

func (a *AlmacenSQLC) BuscarProductoPorID(id int) (models.Producto, bool) {
	f, err := a.q.BuscarProductoPorID(context.Background(), int64(id))
	if err != nil {
		// Absorbemos sql.ErrNoRows (y cualquier otro error) y conservamos la firma comma-ok.
		return models.Producto{}, false
	}
	return aProductoDominio(f), true
}

func (a *AlmacenSQLC) CrearProducto(p models.Producto) models.Producto {
	f, err := a.q.CrearProducto(context.Background(), sqlcdb.CrearProductoParams{
		Nombre:      p.Nombre,
		Precio:      p.Precio,
		Stock:       int64(p.Stock),
		CategoriaID: int64(p.CategoriaID),
	})
	if err != nil {
		// La interfaz no permite reportar el fallo de una creación (igual que Memoria
		// y AlmacenSQLite). Devolvemos el zero value. Ver nota en la guía docente.
		return models.Producto{}
	}
	return aProductoDominio(f)
}

func (a *AlmacenSQLC) ActualizarProducto(id int, datos models.Producto) (models.Producto, bool) {
	f, err := a.q.ActualizarProducto(context.Background(), sqlcdb.ActualizarProductoParams{
		Nombre:      datos.Nombre,
		Precio:      datos.Precio,
		Stock:       int64(datos.Stock),
		CategoriaID: int64(datos.CategoriaID),
		ID:          int64(id),
	})
	if err != nil {
		return models.Producto{}, false
	}
	return aProductoDominio(f), true
}

func (a *AlmacenSQLC) BorrarProducto(id int) bool {
	filas, err := a.q.BorrarProducto(context.Background(), int64(id))
	if err != nil {
		return false
	}
	return filas > 0
}

// =========================================================
// CATEGORIAS
// =========================================================

func (a *AlmacenSQLC) ListarCategorias() []models.Categoria {
	filas, err := a.q.ListarCategorias(context.Background())
	if err != nil {
		return nil
	}
	out := make([]models.Categoria, 0, len(filas))
	for _, f := range filas {
		out = append(out, aCategoriaDominio(f))
	}
	return out
}

func (a *AlmacenSQLC) BuscarCategoriaPorID(id int) (models.Categoria, bool) {
	f, err := a.q.BuscarCategoriaPorID(context.Background(), int64(id))
	if err != nil {
		return models.Categoria{}, false
	}
	return aCategoriaDominio(f), true
}

func (a *AlmacenSQLC) CrearCategoria(c models.Categoria) models.Categoria {
	f, err := a.q.CrearCategoria(context.Background(), sqlcdb.CrearCategoriaParams{
		Nombre:      c.Nombre,
		Descripcion: c.Descripcion,
	})
	if err != nil {
		return models.Categoria{}
	}
	return aCategoriaDominio(f)
}

func (a *AlmacenSQLC) ActualizarCategoria(id int, datos models.Categoria) (models.Categoria, bool) {
	f, err := a.q.ActualizarCategoria(context.Background(), sqlcdb.ActualizarCategoriaParams{
		Nombre:      datos.Nombre,
		Descripcion: datos.Descripcion,
		ID:          int64(id),
	})
	if err != nil {
		return models.Categoria{}, false
	}
	return aCategoriaDominio(f), true
}

func (a *AlmacenSQLC) BorrarCategoria(id int) bool {
	filas, err := a.q.BorrarCategoria(context.Background(), int64(id))
	if err != nil {
		return false
	}
	return filas > 0
}

// Chequeo en tiempo de compilación: AlmacenSQLC debe cumplir Almacen.
var _ Almacen = (*AlmacenSQLC)(nil)
