package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

// productoRepoMock es un doble de prueba de storage.ProductoRepository.
//
// Gracias al ISP, ProductoService depende SOLO de esta interfaz estrecha (5
// metodos), no del Almacen completo (10). Por eso el mock implementa 5 metodos
// y no toca para nada categorias. Si manana ProductoService necesitara mas,
// el compilador nos obligaria a anadirlos aqui: ese es el valor de la asercion
// de abajo.
type productoRepoMock struct {
	mock.Mock
}

func (m *productoRepoMock) ListarProductos() []models.Producto {
	args := m.Called()
	return args.Get(0).([]models.Producto)
}

func (m *productoRepoMock) BuscarProductoPorID(id int) (models.Producto, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Producto), args.Bool(1)
}

func (m *productoRepoMock) CrearProducto(p models.Producto) models.Producto {
	args := m.Called(p)
	return args.Get(0).(models.Producto)
}

func (m *productoRepoMock) ActualizarProducto(id int, datos models.Producto) (models.Producto, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Producto), args.Bool(1)
}

func (m *productoRepoMock) BorrarProducto(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// Red de seguridad en tiempo de compilacion: el mock DEBE cumplir el contrato.
var _ storage.ProductoRepository = (*productoRepoMock)(nil)

// TestProductoService_Crear comprueba la REGLA DE NEGOCIO (validarProducto) de
// forma aislada, sin base de datos. Es un test table-driven: una sola funcion
// recorre varios casos.
func TestProductoService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Producto
		errEsperado   error // nil = se espera exito
		debePersistir bool
	}{
		{
			nombre:        "nombre vacio -> ErrNombreVacio",
			entrada:       models.Producto{Nombre: "   ", Precio: 1.25},
			errEsperado:   service.ErrNombreVacio,
			debePersistir: false,
		},
		{
			nombre:        "precio negativo -> ErrPrecioNegativo",
			entrada:       models.Producto{Nombre: "Cafe americano", Precio: -1.0},
			errEsperado:   service.ErrPrecioNegativo,
			debePersistir: false,
		},
		{
			nombre:        "producto valido -> sin error y se persiste",
			entrada:       models.Producto{Nombre: "Cafe americano", Precio: 1.25, Stock: 10, CategoriaID: 1},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Preparar: un mock nuevo por caso para no arrastrar estado.
			repo := new(productoRepoMock)
			if c.debePersistir {
				// El repo devuelve el producto con un ID asignado.
				guardado := c.entrada
				guardado.ID = 42
				repo.On("CrearProducto", c.entrada).Return(guardado)
			}
			svc := service.NuevoProductoService(repo)

			// Ejecutar.
			creado, err := svc.Crear(c.entrada)

			// Verificar.
			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				// Si la validacion falla, el repo NO debe haberse tocado.
				repo.AssertNotCalled(t, "CrearProducto")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 42, creado.ID, "el service debe devolver el producto que entrego el repo")
				repo.AssertCalled(t, "CrearProducto", c.entrada)
			}
		})
	}
}

// TestProductoService_Obtener_NoEncontrado muestra como el service traduce el
// comma-ok del repositorio (false) en un error de dominio (ErrNoEncontrado).
func TestProductoService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(productoRepoMock)
	repo.On("BuscarProductoPorID", 999).Return(models.Producto{}, false)
	svc := service.NuevoProductoService(repo)

	_, err := svc.Obtener(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
