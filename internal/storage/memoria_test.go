// Tests del almacen en memoria con la libreria ESTANDAR (testing), sin testify.
// Es el estilo mas basico de Go y el que conviene mostrar primero: t.Fatalf
// para abortar, t.Errorf para seguir. Sin dependencias externas.
package storage

import (
	"testing"

	"cafeteria-uleam-api/internal/models"
)

func TestMemoria_CrearYBuscar(t *testing.T) {
	m := NuevaMemoria()

	creado := m.CrearProducto(models.Producto{Nombre: "Cafe americano", Precio: 1.25, Stock: 10})
	if creado.ID == 0 {
		t.Fatalf("esperaba un ID asignado, obtuve 0")
	}

	encontrado, ok := m.BuscarProductoPorID(creado.ID)
	if !ok {
		t.Fatalf("no se encontro el producto recien creado (id=%d)", creado.ID)
	}
	if encontrado.Nombre != "Cafe americano" {
		t.Errorf("nombre = %q; esperaba %q", encontrado.Nombre, "Cafe americano")
	}
}

func TestMemoria_BuscarInexistente(t *testing.T) {
	m := NuevaMemoria()

	// El patron comma-ok: ok debe ser false para un id que no existe.
	if _, ok := m.BuscarProductoPorID(999); ok {
		t.Errorf("esperaba ok=false para un id inexistente")
	}
}

func TestMemoria_ActualizarYBorrar(t *testing.T) {
	m := NuevaMemoria()
	creado := m.CrearProducto(models.Producto{Nombre: "Capuccino", Precio: 1.75})

	if _, ok := m.ActualizarProducto(creado.ID, models.Producto{Nombre: "Capuccino grande", Precio: 2.0}); !ok {
		t.Fatalf("no se pudo actualizar el producto id=%d", creado.ID)
	}

	if !m.BorrarProducto(creado.ID) {
		t.Errorf("esperaba poder borrar el producto id=%d", creado.ID)
	}
	if _, ok := m.BuscarProductoPorID(creado.ID); ok {
		t.Errorf("el producto id=%d deberia haber sido borrado", creado.ID)
	}
}
