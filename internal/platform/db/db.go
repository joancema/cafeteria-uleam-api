// Package db abre la conexión. NO conoce los modelos: cada módulo migra su propio
// esquema (interno.Migrar). Eso respeta la frontera de los módulos.
package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Abrir(rutaDB string) (*gorm.DB, func() error, error) {
	gdb, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("abrir GORM: %w", err)
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
