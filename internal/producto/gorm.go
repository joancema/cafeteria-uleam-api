package producto

import "gorm.io/gorm"

// RepoGORM implementa Repository sobre GORM. En vertical slice el adaptador de
// persistencia vive en el mismo paquete que el resto de la vertical.
type RepoGORM struct {
	db *gorm.DB
}

func NuevoRepoGORM(db *gorm.DB) *RepoGORM {
	return &RepoGORM{db: db}
}

func (a *RepoGORM) Listar() []Producto {
	var productos []Producto
	a.db.Find(&productos)
	return productos
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
