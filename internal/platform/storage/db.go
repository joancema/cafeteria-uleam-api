// Package storage abre la base, migra los modelos de todos los slices y siembra
// datos. Es transversal: importa los slices (sus modelos) en UNA sola direccion;
// ningun slice importa este paquete, asi que no hay ciclo.
package storage

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"cafeteria-uleam-api/internal/auth"
	"cafeteria-uleam-api/internal/categoria"
	"cafeteria-uleam-api/internal/producto"
)

// Abrir conecta GORM, migra los modelos y devuelve la conexion + una funcion de
// cierre para el graceful shutdown.
func Abrir(rutaDB string) (*gorm.DB, func() error, error) {
	gdb, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("abrir GORM: %w", err)
	}
	if err := gdb.AutoMigrate(&producto.Producto{}, &categoria.Categoria{}, &auth.Usuario{}); err != nil {
		return nil, nil, fmt.Errorf("AutoMigrate: %w", err)
	}
	cerrar := func() error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return gdb, cerrar, nil
}

// Sembrar inserta datos iniciales solo si aun no hay categorias.
func Sembrar(gdb *gorm.DB) {
	var n int64
	gdb.Model(&categoria.Categoria{}).Count(&n)
	if n > 0 {
		return
	}
	categorias := []categoria.Categoria{
		{ID: 1, Nombre: "Bebidas calientes", Descripcion: "Cafes, tes e infusiones"},
		{ID: 2, Nombre: "Alimentos solidos", Descripcion: "Sandwiches y comida lista"},
	}
	gdb.Create(&categorias)
	productos := []producto.Producto{
		{ID: 1, Nombre: "Cafe americano", Precio: 1.25, Stock: 50, CategoriaID: 1},
		{ID: 2, Nombre: "Sandwich vegetariano", Precio: 2.50, Stock: 20, CategoriaID: 2},
	}
	gdb.Create(&productos)
}
