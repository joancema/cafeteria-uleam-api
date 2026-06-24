package auth

type Repository interface {
	Crear(u Usuario) (Usuario, error)
	BuscarPorEmail(email string) (Usuario, bool)
}
