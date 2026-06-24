package producto

import "strings"

// Service contiene la logica de negocio de productos. Depende del Repository
// (interfaz), no del adaptador concreto.
type Service struct {
	repo Repository
}

func NuevoService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Listar() []Producto {
	return s.repo.Listar()
}

func (s *Service) Obtener(id int) (Producto, error) {
	p, ok := s.repo.BuscarPorID(id)
	if !ok {
		return Producto{}, ErrNoEncontrado
	}
	return p, nil
}

func (s *Service) Crear(p Producto) (Producto, error) {
	if err := validar(p); err != nil {
		return Producto{}, err
	}
	return s.repo.Crear(p), nil
}

func (s *Service) Actualizar(id int, datos Producto) (Producto, error) {
	if err := validar(datos); err != nil {
		return Producto{}, err
	}
	actualizado, ok := s.repo.Actualizar(id, datos)
	if !ok {
		return Producto{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *Service) Borrar(id int) error {
	if !s.repo.Borrar(id) {
		return ErrNoEncontrado
	}
	return nil
}

func validar(p Producto) error {
	if strings.TrimSpace(p.Nombre) == "" {
		return ErrNombreVacio
	}
	if p.Precio < 0 {
		return ErrPrecioNegativo
	}
	return nil
}
