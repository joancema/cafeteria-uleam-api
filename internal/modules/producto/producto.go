// Package producto es la FACHADA pública del módulo de productos: lo ÚNICO que el
// resto del sistema puede usar. Las tripas viven en ./internal y Go impide que
// otros módulos las importen.
package producto

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"cafeteria-mod/internal/modules/producto/internal"
)

// Modulo agrupa el módulo ya ensamblado.
type Modulo struct {
	handler *interno.Handler
}

// Nuevo migra el esquema propio del módulo y arma sus capas internas.
func Nuevo(db *gorm.DB) *Modulo {
	interno.Migrar(db)
	repo := interno.NuevoRepoGORM(db)
	svc := interno.NuevoService(repo)
	return &Modulo{handler: interno.NuevoHandler(svc)}
}

// Rutas publica los endpoints del módulo (montar bajo /productos).
func (m *Modulo) Rutas(r chi.Router) {
	m.handler.Rutas(r)
}
