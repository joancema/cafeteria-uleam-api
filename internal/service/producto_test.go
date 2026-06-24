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

// productoRepoMock es un doble de storage.ProductoRepository (la interfaz
// estrecha de 5 metodos). Cada metodo solo registra la llamada y devuelve lo que
// el test configuro con On(...). No persiste nada.
type productoRepoMock struct {
	mock.Mock
}

func (m *productoRepoMock) ListarProductos() []models.Producto {
	return m.Called().Get(0).([]models.Producto)
}
func (m *productoRepoMock) BuscarProductoPorID(id int) (models.Producto, bool) {
	a := m.Called(id)
	return a.Get(0).(models.Producto), a.Bool(1)
}
func (m *productoRepoMock) CrearProducto(p models.Producto) models.Producto {
	return m.Called(p).Get(0).(models.Producto)
}
func (m *productoRepoMock) ActualizarProducto(id int, datos models.Producto) (models.Producto, bool) {
	a := m.Called(id, datos)
	return a.Get(0).(models.Producto), a.Bool(1)
}
func (m *productoRepoMock) BorrarProducto(id int) bool {
	return m.Called(id).Bool(0)
}

// Red de seguridad en tiempo de compilacion: el mock DEBE cumplir el contrato.
var _ storage.ProductoRepository = (*productoRepoMock)(nil)

// --- Crear: la regla de negocio (validarProducto), aislada de la base ---

func TestProductoService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Producto
		errEsperado   error // nil = exito
		debePersistir bool
	}{
		{"nombre vacio rechazado", models.Producto{Nombre: "   ", Precio: 1.25}, service.ErrNombreVacio, false},
		{"precio negativo rechazado", models.Producto{Nombre: "Cafe", Precio: -1}, service.ErrPrecioNegativo, false},
		{"producto valido se persiste", models.Producto{Nombre: "Cafe", Precio: 1.25, Stock: 10, CategoriaID: 1}, nil, true},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(productoRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 42
				repo.On("CrearProducto", c.entrada).Return(guardado)
			}
			svc := service.NuevoProductoService(repo)

			creado, err := svc.Crear(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearProducto") // la validacion corto antes
			} else {
				require.NoError(t, err)
				assert.Equal(t, 42, creado.ID)
				repo.AssertCalled(t, "CrearProducto", c.entrada)
			}
		})
	}
}

// --- Obtener: comma-ok del repo traducido a error de dominio ---

func TestProductoService_Obtener(t *testing.T) {
	t.Run("existe", func(t *testing.T) {
		repo := new(productoRepoMock)
		repo.On("BuscarProductoPorID", 1).Return(models.Producto{ID: 1, Nombre: "Cafe"}, true)
		p, err := service.NuevoProductoService(repo).Obtener(1)
		require.NoError(t, err)
		assert.Equal(t, "Cafe", p.Nombre)
	})
	t.Run("no existe -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(productoRepoMock)
		repo.On("BuscarProductoPorID", 999).Return(models.Producto{}, false)
		_, err := service.NuevoProductoService(repo).Obtener(999)
		require.ErrorIs(t, err, service.ErrNoEncontrado)
	})
}

// --- Actualizar: valida ANTES de tocar el repo, y mapea el no encontrado ---

func TestProductoService_Actualizar(t *testing.T) {
	datos := models.Producto{Nombre: "Cafe doble", Precio: 1.5}

	t.Run("valido", func(t *testing.T) {
		repo := new(productoRepoMock)
		actualizado := datos
		actualizado.ID = 1
		repo.On("ActualizarProducto", 1, datos).Return(actualizado, true)
		p, err := service.NuevoProductoService(repo).Actualizar(1, datos)
		require.NoError(t, err)
		assert.Equal(t, 1, p.ID)
	})
	t.Run("no existe -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(productoRepoMock)
		repo.On("ActualizarProducto", 999, datos).Return(models.Producto{}, false)
		_, err := service.NuevoProductoService(repo).Actualizar(999, datos)
		require.ErrorIs(t, err, service.ErrNoEncontrado)
	})
	t.Run("invalido no toca el repo", func(t *testing.T) {
		repo := new(productoRepoMock)
		_, err := service.NuevoProductoService(repo).Actualizar(1, models.Producto{Nombre: ""})
		require.ErrorIs(t, err, service.ErrNombreVacio)
		repo.AssertNotCalled(t, "ActualizarProducto")
	})
}

// --- Borrar ---

func TestProductoService_Borrar(t *testing.T) {
	t.Run("existe", func(t *testing.T) {
		repo := new(productoRepoMock)
		repo.On("BorrarProducto", 1).Return(true)
		require.NoError(t, service.NuevoProductoService(repo).Borrar(1))
	})
	t.Run("no existe -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(productoRepoMock)
		repo.On("BorrarProducto", 999).Return(false)
		require.ErrorIs(t, service.NuevoProductoService(repo).Borrar(999), service.ErrNoEncontrado)
	})
}

// --- Listar: el service solo delega ---

func TestProductoService_Listar(t *testing.T) {
	repo := new(productoRepoMock)
	repo.On("ListarProductos").Return([]models.Producto{{ID: 1}, {ID: 2}})
	lista := service.NuevoProductoService(repo).Listar()
	assert.Len(t, lista, 2)
	repo.AssertExpectations(t)
}
