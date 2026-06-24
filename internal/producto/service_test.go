package producto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/producto"
)

// repoMock es un doble de producto.Repository (5 metodos, sin sufijo de entidad).
type repoMock struct{ mock.Mock }

func (m *repoMock) Listar() []producto.Producto {
	return m.Called().Get(0).([]producto.Producto)
}
func (m *repoMock) BuscarPorID(id int) (producto.Producto, bool) {
	a := m.Called(id)
	return a.Get(0).(producto.Producto), a.Bool(1)
}
func (m *repoMock) Crear(p producto.Producto) producto.Producto {
	return m.Called(p).Get(0).(producto.Producto)
}
func (m *repoMock) Actualizar(id int, datos producto.Producto) (producto.Producto, bool) {
	a := m.Called(id, datos)
	return a.Get(0).(producto.Producto), a.Bool(1)
}
func (m *repoMock) Borrar(id int) bool {
	return m.Called(id).Bool(0)
}

var _ producto.Repository = (*repoMock)(nil)

func TestService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       producto.Producto
		errEsperado   error
		debePersistir bool
	}{
		{"nombre vacio rechazado", producto.Producto{Nombre: "   ", Precio: 1.25}, producto.ErrNombreVacio, false},
		{"precio negativo rechazado", producto.Producto{Nombre: "Cafe", Precio: -1}, producto.ErrPrecioNegativo, false},
		{"producto valido se persiste", producto.Producto{Nombre: "Cafe", Precio: 1.25, Stock: 10, CategoriaID: 1}, nil, true},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(repoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 42
				repo.On("Crear", c.entrada).Return(guardado)
			}
			svc := producto.NuevoService(repo)

			creado, err := svc.Crear(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "Crear")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 42, creado.ID)
				repo.AssertCalled(t, "Crear", c.entrada)
			}
		})
	}
}

func TestService_Obtener(t *testing.T) {
	t.Run("no existe -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("BuscarPorID", 999).Return(producto.Producto{}, false)
		_, err := producto.NuevoService(repo).Obtener(999)
		require.ErrorIs(t, err, producto.ErrNoEncontrado)
	})
}
