// Package auth es la FACHADA pública del módulo de autenticación.
package auth

import (
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"cafeteria-mod/internal/modules/auth/internal"
)

type Modulo struct {
	handler *interno.Handler
	svc     *interno.Service
}

func Nuevo(db *gorm.DB, secreto []byte, duracion time.Duration) *Modulo {
	interno.Migrar(db)
	repo := interno.NuevoRepoGORM(db)
	svc := interno.NuevoService(repo, secreto, duracion)
	return &Modulo{handler: interno.NuevoHandler(svc), svc: svc}
}

// RegistrarRutas publica los endpoints públicos del módulo (register/login).
func (m *Modulo) RegistrarRutas(r chi.Router) {
	r.Post("/auth/register", m.handler.Registrar)
	r.Post("/auth/login", m.handler.Login)
}

// ValidarToken es el CONTRATO PÚBLICO que otros módulos / la plataforma usan para
// autenticar. La fachada delega en sus tripas; nadie fuera del módulo ve el
// Service interno. Así otro módulo depende de este método, no de las tripas.
func (m *Modulo) ValidarToken(token string) (int, error) {
	return m.svc.ValidarToken(token)
}
