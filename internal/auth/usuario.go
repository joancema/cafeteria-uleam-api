// Package auth agrupa la vertical de autenticacion: usuarios, JWT y el servicio
// que el middleware usa para validar tokens.
package auth

import "time"

type Usuario struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreadoEn     time.Time `json:"creado_en"`
}
