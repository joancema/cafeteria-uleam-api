package storage

import (
	"time"

	"gorm.io/gorm"

	"cafeteria-uleam-api/internal/models"
)

// UsuarioGORM implementa UserRepository sobre GORM.
//
// A diferencia de productos/categorias, los usuarios viven SOLO en GORM: no los
// replicamos en Memoria ni en sqlc (hacerlo seria el mismo ejercicio que ya
// hicimos). El AuthService recibe esta interfaz, nunca este tipo concreto.
type UsuarioGORM struct {
	db *gorm.DB
}

// NuevoUsuarioGORM envuelve una conexion *gorm.DB ya abierta.
func NuevoUsuarioGORM(db *gorm.DB) *UsuarioGORM {
	return &UsuarioGORM{db: db}
}

// CrearUsuario inserta un usuario. Devuelve error si el insert falla (por
// ejemplo, email duplicado por el indice unico).
func (r *UsuarioGORM) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.CreadoEn = time.Now()
	if err := r.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

// BuscarUsuarioPorEmail busca por email; comma-ok como el resto de lecturas.
func (r *UsuarioGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

// Chequeo en tiempo de compilacion: UsuarioGORM debe cumplir UserRepository.
var _ UserRepository = (*UsuarioGORM)(nil)
