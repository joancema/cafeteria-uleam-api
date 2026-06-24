package auth

import (
	"time"

	"gorm.io/gorm"
)

type RepoGORM struct {
	db *gorm.DB
}

func NuevoRepoGORM(db *gorm.DB) *RepoGORM {
	return &RepoGORM{db: db}
}

func (r *RepoGORM) Crear(u Usuario) (Usuario, error) {
	u.CreadoEn = time.Now()
	if err := r.db.Create(&u).Error; err != nil {
		return Usuario{}, err
	}
	return u, nil
}

func (r *RepoGORM) BuscarPorEmail(email string) (Usuario, bool) {
	var u Usuario
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return Usuario{}, false
	}
	return u, true
}

var _ Repository = (*RepoGORM)(nil)
