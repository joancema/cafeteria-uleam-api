package producto

import "strings"

// servicio implementa el puerto de entrada Servicio. Es privado (minúscula): el
// exterior solo lo conoce a través de la interfaz Servicio, no del struct.
type servicio struct {
	repo Repositorio
}

// NuevoServicio inyecta el puerto de salida y DEVUELVE el puerto de entrada.
// Quien lo llame recibe una interfaz, no un struct concreto.
func NuevoServicio(repo Repositorio) Servicio {
	return &servicio{repo: repo}
}

func (s *servicio) Listar() []Producto {
	return s.repo.Listar()
}

func (s *servicio) Obtener(id int) (Producto, error) {
	p, ok := s.repo.BuscarPorID(id)
	if !ok {
		return Producto{}, ErrNoEncontrado
	}
	return p, nil
}

func (s *servicio) Crear(p Producto) (Producto, error) {
	if err := validar(p); err != nil {
		return Producto{}, err
	}
	return s.repo.Crear(p), nil
}

func (s *servicio) Actualizar(id int, datos Producto) (Producto, error) {
	if err := validar(datos); err != nil {
		return Producto{}, err
	}
	actualizado, ok := s.repo.Actualizar(id, datos)
	if !ok {
		return Producto{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *servicio) Borrar(id int) error {
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
