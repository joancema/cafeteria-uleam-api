package service

import (
	"strings"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/storage"
)

// CategoriaService contiene la logica de negocio de categorias.
// Depende solo de storage.CategoriaRepository.
type CategoriaService struct {
	repo storage.CategoriaRepository
}

func NuevoCategoriaService(repo storage.CategoriaRepository) *CategoriaService {
	return &CategoriaService{repo: repo}
}

func (s *CategoriaService) Listar() []models.Categoria {
	return s.repo.ListarCategorias()
}

func (s *CategoriaService) Obtener(id int) (models.Categoria, error) {
	c, ok := s.repo.BuscarCategoriaPorID(id)
	if !ok {
		return models.Categoria{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *CategoriaService) Crear(c models.Categoria) (models.Categoria, error) {
	if strings.TrimSpace(c.Nombre) == "" {
		return models.Categoria{}, ErrNombreVacio
	}
	return s.repo.CrearCategoria(c), nil
}

func (s *CategoriaService) Actualizar(id int, datos models.Categoria) (models.Categoria, error) {
	if strings.TrimSpace(datos.Nombre) == "" {
		return models.Categoria{}, ErrNombreVacio
	}
	actualizada, ok := s.repo.ActualizarCategoria(id, datos)
	if !ok {
		return models.Categoria{}, ErrNoEncontrado
	}
	return actualizada, nil
}

func (s *CategoriaService) Borrar(id int) error {
	if !s.repo.BorrarCategoria(id) {
		return ErrNoEncontrado
	}
	return nil
}
