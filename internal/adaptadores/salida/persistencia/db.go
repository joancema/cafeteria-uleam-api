package persistencia

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"cafeteria-hex/internal/core/producto"
)

func Abrir(rutaDB string) (*gorm.DB, func() error, error) {
	gdb, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("abrir GORM: %w", err)
	}
	if err := gdb.AutoMigrate(&producto.Producto{}); err != nil {
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

func Sembrar(gdb *gorm.DB) {
	var n int64
	gdb.Model(&producto.Producto{}).Count(&n)
	if n > 0 {
		return
	}
	gdb.Create(&[]producto.Producto{
		{ID: 1, Nombre: "Cafe americano", Precio: 1.25, Stock: 50, CategoriaID: 1},
		{ID: 2, Nombre: "Sandwich vegetariano", Precio: 2.50, Stock: 20, CategoriaID: 2},
	})
}
