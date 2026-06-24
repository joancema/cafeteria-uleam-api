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

// categoriaRepoMock: doble de storage.CategoriaRepository (5 metodos).
type categoriaRepoMock struct {
	mock.Mock
}

func (m *categoriaRepoMock) ListarCategorias() []models.Categoria {
	return m.Called().Get(0).([]models.Categoria)
}
func (m *categoriaRepoMock) BuscarCategoriaPorID(id int) (models.Categoria, bool) {
	a := m.Called(id)
	return a.Get(0).(models.Categoria), a.Bool(1)
}
func (m *categoriaRepoMock) CrearCategoria(c models.Categoria) models.Categoria {
	return m.Called(c).Get(0).(models.Categoria)
}
func (m *categoriaRepoMock) ActualizarCategoria(id int, datos models.Categoria) (models.Categoria, bool) {
	a := m.Called(id, datos)
	return a.Get(0).(models.Categoria), a.Bool(1)
}
func (m *categoriaRepoMock) BorrarCategoria(id int) bool {
	return m.Called(id).Bool(0)
}

var _ storage.CategoriaRepository = (*categoriaRepoMock)(nil)

// La categoria tiene su propia regla: nombre obligatorio.
func TestCategoriaService_Crear(t *testing.T) {
	t.Run("nombre vacio rechazado", func(t *testing.T) {
		repo := new(categoriaRepoMock)
		_, err := service.NuevoCategoriaService(repo).Crear(models.Categoria{Nombre: "  "})
		require.ErrorIs(t, err, service.ErrNombreVacio)
		repo.AssertNotCalled(t, "CrearCategoria")
	})
	t.Run("valida se persiste", func(t *testing.T) {
		repo := new(categoriaRepoMock)
		entrada := models.Categoria{Nombre: "Postres", Descripcion: "Dulces"}
		guardada := entrada
		guardada.ID = 7
		repo.On("CrearCategoria", entrada).Return(guardada)
		c, err := service.NuevoCategoriaService(repo).Crear(entrada)
		require.NoError(t, err)
		assert.Equal(t, 7, c.ID)
	})
}

func TestCategoriaService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(categoriaRepoMock)
	repo.On("BuscarCategoriaPorID", 999).Return(models.Categoria{}, false)
	_, err := service.NuevoCategoriaService(repo).Obtener(999)
	require.ErrorIs(t, err, service.ErrNoEncontrado)
}

func TestCategoriaService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(categoriaRepoMock)
	repo.On("BorrarCategoria", 999).Return(false)
	require.ErrorIs(t, service.NuevoCategoriaService(repo).Borrar(999), service.ErrNoEncontrado)
}
