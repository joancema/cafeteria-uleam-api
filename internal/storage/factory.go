package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite" // driver database/sql "sqlite" (pure-Go) para el backend sqlc
	"github.com/glebarez/sqlite"      // driver GORM (pure-Go)
	"gorm.io/gorm"

	"cafeteria-uleam-api/internal/models"
)

// Recursos agrupa todo lo que la capa de almacenamiento expone a la aplicacion:
// el almacen de productos/categorias (segun el backend elegido), el repositorio
// de usuarios (siempre GORM) y una funcion para cerrar conexiones al apagar.
type Recursos struct {
	Almacen      Almacen
	Usuarios     UserRepository
	BackendUsado string
	Cerrar       func() error
}

// Inicializar centraliza TODO el plumbing que antes vivia suelto en main.go:
//
//  1. Abre GORM (dueno del esquema), migra y siembra.
//  2. Elige el backend de productos/categorias segun el parametro backend.
//  3. Crea el repositorio de usuarios (siempre GORM).
//  4. Expone una funcion Cerrar para el graceful shutdown.
//
// Es un patron Factory: la aplicacion pide "dame los recursos para esta config"
// y no necesita saber como se arman ni que drivers se importan.
func Inicializar(rutaDB, backend string) (*Recursos, error) {
	// 1. GORM es el DUENO DEL ESQUEMA: abre, migra y siembra.
	gdb, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("abrir GORM: %w", err)
	}
	if err := gdb.AutoMigrate(&models.Producto{}, &models.Categoria{}, &models.Usuario{}); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}
	almacenGorm := NuevoAlmacenSQLite(gdb)
	almacenGorm.SembrarSiVacio()

	// 2. Elegir el backend de productos/categorias (este switch ES el Factory).
	var almacen Almacen
	var sdb *sql.DB
	backendUsado := "gorm"
	switch backend {
	case "sqlc":
		sdb, err = sql.Open("sqlite", rutaDB)
		if err != nil {
			return nil, fmt.Errorf("abrir sql.DB para sqlc: %w", err)
		}
		almacen = NuevoAlmacenSQLC(sdb)
		backendUsado = "sqlc"
	default:
		almacen = almacenGorm
	}

	// 3. Usuarios viven SIEMPRE en GORM (decision tomada en S10).
	usuarios := NuevoUsuarioGORM(gdb)

	// 4. Cierre ordenado: primero la conexion sql.DB de sqlc (si existe), luego GORM.
	cerrar := func() error {
		if sdb != nil {
			if err := sdb.Close(); err != nil {
				return err
			}
		}
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		Almacen:      almacen,
		Usuarios:     usuarios,
		BackendUsado: backendUsado,
		Cerrar:       cerrar,
	}, nil
}
