package categoria

import "strings"

type Service struct {
	repo Repository
}

func NuevoService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Listar() []Categoria {
	return s.repo.Listar()
}

func (s *Service) Obtener(id int) (Categoria, error) {
	c, ok := s.repo.BuscarPorID(id)
	if !ok {
		return Categoria{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *Service) Crear(c Categoria) (Categoria, error) {
	if strings.TrimSpace(c.Nombre) == "" {
		return Categoria{}, ErrNombreVacio
	}
	return s.repo.Crear(c), nil
}

func (s *Service) Actualizar(id int, datos Categoria) (Categoria, error) {
	if strings.TrimSpace(datos.Nombre) == "" {
		return Categoria{}, ErrNombreVacio
	}
	actualizada, ok := s.repo.Actualizar(id, datos)
	if !ok {
		return Categoria{}, ErrNoEncontrado
	}
	return actualizada, nil
}

func (s *Service) Borrar(id int) error {
	if !s.repo.Borrar(id) {
		return ErrNoEncontrado
	}
	return nil
}
