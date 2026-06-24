// Tests del repositorio REAL (GORM) contra una SQLite EN MEMORIA y desechable.
//
// A diferencia de los tests de service/handler —que usan dobles para evitar la
// base—, AQUI si tocamos una base de datos: es el unico lugar donde probamos la
// persistencia de verdad. Usamos ":memory:" para que sea rapida y sin archivos,
// y libreria testing estandar (sin testify).
//
// Nota: requiere el driver glebarez/sqlite (el mismo de produccion). Corre con
// `go test ./internal/storage`.
package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"cafeteria-uleam-api/internal/models"
)

// nuevaDBPrueba abre una SQLite en memoria, migra el esquema y la devuelve.
// SetMaxOpenConns(1) garantiza que migracion y consultas usen la MISMA conexion
// (con ":memory:" cada conexion tendria su propia base, vacia).
func nuevaDBPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir la base de prueba: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("no se pudo obtener *sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := gdb.AutoMigrate(&models.Producto{}, &models.Categoria{}, &models.Usuario{}); err != nil {
		t.Fatalf("fallo AutoMigrate: %v", err)
	}
	return gdb
}

func TestSQLite_ProductoCRUD(t *testing.T) {
	alm := NuevoAlmacenSQLite(nuevaDBPrueba(t))

	// Crear: GORM debe asignar el ID autogenerado.
	creado := alm.CrearProducto(models.Producto{Nombre: "Cafe americano", Precio: 1.25, Stock: 10})
	if creado.ID == 0 {
		t.Fatalf("esperaba un ID asignado por la base, obtuve 0")
	}

	// Buscar el recien creado.
	encontrado, ok := alm.BuscarProductoPorID(creado.ID)
	if !ok {
		t.Fatalf("no se encontro el producto id=%d", creado.ID)
	}
	if encontrado.Nombre != "Cafe americano" {
		t.Errorf("nombre = %q; esperaba %q", encontrado.Nombre, "Cafe americano")
	}

	// Actualizar.
	if _, ok := alm.ActualizarProducto(creado.ID, models.Producto{Nombre: "Cafe doble", Precio: 1.80}); !ok {
		t.Fatalf("no se pudo actualizar el producto id=%d", creado.ID)
	}

	// Borrar y confirmar que ya no esta.
	if !alm.BorrarProducto(creado.ID) {
		t.Errorf("esperaba poder borrar el producto id=%d", creado.ID)
	}
	if _, ok := alm.BuscarProductoPorID(creado.ID); ok {
		t.Errorf("el producto id=%d deberia haber sido borrado", creado.ID)
	}
}

func TestSQLite_BuscarInexistente(t *testing.T) {
	alm := NuevoAlmacenSQLite(nuevaDBPrueba(t))
	// El error de GORM (registro no encontrado) se traduce a comma-ok = false.
	if _, ok := alm.BuscarProductoPorID(999); ok {
		t.Errorf("esperaba ok=false para un id inexistente")
	}
}

// TestSQLite_UsuarioEmailUnico prueba una garantia que SOLO la base puede dar:
// el indice unico de email impide dos usuarios con el mismo correo.
func TestSQLite_UsuarioEmailUnico(t *testing.T) {
	repo := NuevoUsuarioGORM(nuevaDBPrueba(t))

	if _, err := repo.CrearUsuario(models.Usuario{Email: "ana@uleam.edu.ec", PasswordHash: "hash1"}); err != nil {
		t.Fatalf("el primer usuario deberia crearse sin error: %v", err)
	}
	if _, err := repo.CrearUsuario(models.Usuario{Email: "ana@uleam.edu.ec", PasswordHash: "hash2"}); err == nil {
		t.Errorf("esperaba error por email duplicado (indice unico), no lo hubo")
	}
}
