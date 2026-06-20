package service

import (
	"strings"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/storage"
)

// ProductoService contiene la logica de negocio de productos.
//
// Depende SOLO de storage.ProductoRepository (la interfaz estrecha), no del
// Almacen completo: por eso un test podria inyectar un mock de 5 metodos.
type ProductoService struct {
	repo storage.ProductoRepository
}

func NuevoProductoService(repo storage.ProductoRepository) *ProductoService {
	return &ProductoService{repo: repo}
}

func (s *ProductoService) Listar() []models.Producto {
	return s.repo.ListarProductos()
}

func (s *ProductoService) Obtener(id int) (models.Producto, error) {
	p, ok := s.repo.BuscarProductoPorID(id)
	if !ok {
		return models.Producto{}, ErrNoEncontrado
	}
	return p, nil
}

func (s *ProductoService) Crear(p models.Producto) (models.Producto, error) {
	if err := validarProducto(p); err != nil {
		return models.Producto{}, err
	}
	return s.repo.CrearProducto(p), nil
}

func (s *ProductoService) Actualizar(id int, datos models.Producto) (models.Producto, error) {
	if err := validarProducto(datos); err != nil {
		return models.Producto{}, err
	}
	actualizado, ok := s.repo.ActualizarProducto(id, datos)
	if !ok {
		return models.Producto{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *ProductoService) Borrar(id int) error {
	if !s.repo.BorrarProducto(id) {
		return ErrNoEncontrado
	}
	return nil
}

// validarProducto centraliza las reglas de negocio que antes vivian en el handler.
func validarProducto(p models.Producto) error {
	if strings.TrimSpace(p.Nombre) == "" {
		return ErrNombreVacio
	}
	if p.Precio < 0 {
		return ErrPrecioNegativo
	}
	return nil
}
