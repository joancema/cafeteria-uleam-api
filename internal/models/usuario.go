package models

import "time"

// Usuario representa una cuenta que puede autenticarse en la API.
//
// PasswordHash NUNCA se serializa a JSON (tag "-"): el hash bcrypt jamás debe
// salir en una respuesta. Email es unico: el indice unico de GORM impide dos
// usuarios con el mismo correo.
type Usuario struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreadoEn     time.Time `json:"creado_en"`
}
