package categoria

import "gorm.io/gorm"

type RepoGORM struct {
	db *gorm.DB
}

func NuevoRepoGORM(db *gorm.DB) *RepoGORM {
	return &RepoGORM{db: db}
}

func (a *RepoGORM) Listar() []Categoria {
	var categorias []Categoria
	a.db.Find(&categorias)
	return categorias
}

func (a *RepoGORM) BuscarPorID(id int) (Categoria, bool) {
	var c Categoria
	if err := a.db.First(&c, id).Error; err != nil {
		return Categoria{}, false
	}
	return c, true
}

func (a *RepoGORM) Crear(c Categoria) Categoria {
	a.db.Create(&c)
	return c
}

func (a *RepoGORM) Actualizar(id int, datos Categoria) (Categoria, bool) {
	var existente Categoria
	if err := a.db.First(&existente, id).Error; err != nil {
		return Categoria{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *RepoGORM) Borrar(id int) bool {
	return a.db.Delete(&Categoria{}, id).RowsAffected > 0
}

var _ Repository = (*RepoGORM)(nil)
