// Package interno contiene las tripas del módulo de productos. Al vivir bajo un
// directorio internal/, Go GARANTIZA que solo el módulo producto puede importarlo:
// ningún otro módulo (ni auth ni nadie) puede tocar estos tipos. La frontera del
// módulo está forzada por el lenguaje.
package interno

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// ---------- modelo + errores ----------

type Producto struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Nombre      string  `json:"nombre" gorm:"not null"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	CategoriaID int     `json:"categoria_id"`
}

var (
	ErrNombreVacio    = errors.New("el campo nombre es obligatorio")
	ErrPrecioNegativo = errors.New("el precio no puede ser negativo")
	ErrNoEncontrado   = errors.New("producto no encontrado")
)

// Migrar crea la tabla del módulo y siembra datos. CADA módulo es dueño de su
// propio esquema: la plataforma solo abre la conexión, no conoce los modelos.
func Migrar(db *gorm.DB) {
	db.AutoMigrate(&Producto{})
	var n int64
	db.Model(&Producto{}).Count(&n)
	if n == 0 {
		db.Create(&[]Producto{
			{ID: 1, Nombre: "Cafe americano", Precio: 1.25, Stock: 50, CategoriaID: 1},
			{ID: 2, Nombre: "Sandwich vegetariano", Precio: 2.50, Stock: 20, CategoriaID: 2},
		})
	}
}

// ---------- repositorio ----------

type Repository interface {
	Listar() []Producto
	BuscarPorID(id int) (Producto, bool)
	Crear(p Producto) Producto
	Actualizar(id int, datos Producto) (Producto, bool)
	Borrar(id int) bool
}

type RepoGORM struct{ db *gorm.DB }

func NuevoRepoGORM(db *gorm.DB) *RepoGORM { return &RepoGORM{db: db} }

func (a *RepoGORM) Listar() []Producto {
	var ps []Producto
	a.db.Find(&ps)
	return ps
}
func (a *RepoGORM) BuscarPorID(id int) (Producto, bool) {
	var p Producto
	if err := a.db.First(&p, id).Error; err != nil {
		return Producto{}, false
	}
	return p, true
}
func (a *RepoGORM) Crear(p Producto) Producto {
	a.db.Create(&p)
	return p
}
func (a *RepoGORM) Actualizar(id int, datos Producto) (Producto, bool) {
	var existente Producto
	if err := a.db.First(&existente, id).Error; err != nil {
		return Producto{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}
func (a *RepoGORM) Borrar(id int) bool {
	return a.db.Delete(&Producto{}, id).RowsAffected > 0
}

var _ Repository = (*RepoGORM)(nil)

// ---------- service ----------

type Service struct{ repo Repository }

func NuevoService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Listar() []Producto { return s.repo.Listar() }
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
	a, ok := s.repo.Actualizar(id, datos)
	if !ok {
		return Producto{}, ErrNoEncontrado
	}
	return a, nil
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

// ---------- handler ----------

type Handler struct{ svc *Service }

func NuevoHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Rutas(r chi.Router) {
	r.Get("/", h.listar)
	r.Post("/", h.crear)
	r.Get("/{id}", h.obtener)
	r.Put("/{id}", h.actualizar)
	r.Delete("/{id}", h.borrar)
}

func (h *Handler) listar(w http.ResponseWriter, _ *http.Request) {
	responderJSON(w, http.StatusOK, h.svc.Listar())
}
func (h *Handler) obtener(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	p, err := h.svc.Obtener(id)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusOK, p)
}
func (h *Handler) crear(w http.ResponseWriter, r *http.Request) {
	var nuevo Producto
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creado, err := h.svc.Crear(nuevo)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, creado)
}
func (h *Handler) actualizar(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	var datos Producto
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	a, err := h.svc.Actualizar(id, datos)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusOK, a)
}
func (h *Handler) borrar(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id debe ser un numero entero")
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusNoContent, nil)
}

func idDeURL(r *http.Request) (int, error) { return strconv.Atoi(chi.URLParam(r, "id")) }
func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
func responderError(w http.ResponseWriter, status int, msg string) {
	responderJSON(w, status, map[string]string{"error": msg})
}
func statusDeError(err error) int {
	switch {
	case errors.Is(err, ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, ErrNombreVacio), errors.Is(err, ErrPrecioNegativo):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
