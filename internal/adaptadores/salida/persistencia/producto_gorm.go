// Package persistencia contiene los adaptadores de SALIDA: implementan los
// puertos de salida del núcleo (aquí, producto.Repositorio) usando GORM.
package persistencia

import (
	"gorm.io/gorm"

	"cafeteria-hex/internal/core/producto"
)

// ProductoGORM es el adaptador: traduce entre el puerto del núcleo y GORM.
type ProductoGORM struct {
	db *gorm.DB
}

func NuevoProductoGORM(db *gorm.DB) *ProductoGORM {
	return &ProductoGORM{db: db}
}

func (a *ProductoGORM) Listar() []producto.Producto {
	var ps []producto.Producto
	a.db.Find(&ps)
	return ps
}

func (a *ProductoGORM) BuscarPorID(id int) (producto.Producto, bool) {
	var p producto.Producto
	if err := a.db.First(&p, id).Error; err != nil {
		return producto.Producto{}, false
	}
	return p, true
}

func (a *ProductoGORM) Crear(p producto.Producto) producto.Producto {
	a.db.Create(&p)
	return p
}

func (a *ProductoGORM) Actualizar(id int, datos producto.Producto) (producto.Producto, bool) {
	var existente producto.Producto
	if err := a.db.First(&existente, id).Error; err != nil {
		return producto.Producto{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *ProductoGORM) Borrar(id int) bool {
	return a.db.Delete(&producto.Producto{}, id).RowsAffected > 0
}

// El adaptador DEBE satisfacer el puerto del núcleo.
var _ producto.Repositorio = (*ProductoGORM)(nil)
