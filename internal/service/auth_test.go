package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

// usuarioRepoFake es un repositorio de usuarios EN MEMORIA para los tests.
//
// En produccion el AuthService recibe el repositorio GORM; aqui le inyectamos
// este doble. El AuthService no nota la diferencia: depende de la interfaz
// storage.UserRepository, no de GORM. Esa es la ventaja del Repository Pattern.
type usuarioRepoFake struct {
	porEmail map[string]models.Usuario
	nextID   int
}

func nuevoUsuarioRepoFake() *usuarioRepoFake {
	return &usuarioRepoFake{porEmail: map[string]models.Usuario{}, nextID: 1}
}

func (f *usuarioRepoFake) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = f.nextID
	f.nextID++
	f.porEmail[u.Email] = u
	return u, nil
}

func (f *usuarioRepoFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

var _ storage.UserRepository = (*usuarioRepoFake)(nil)

// TestAuthService_Registrar_EmailEnUso: registrar dos veces el mismo email debe
// fallar con ErrEmailEnUso (que el handler traducira a 409).
func TestAuthService_Registrar_EmailEnUso(t *testing.T) {
	svc := service.NuevoAuthService(nuevoUsuarioRepoFake())

	_, err := svc.Registrar("ana@uleam.edu.ec", "secreta123")
	require.NoError(t, err)

	_, err = svc.Registrar("ana@uleam.edu.ec", "otraClave456")
	require.ErrorIs(t, err, service.ErrEmailEnUso)
}

// TestAuthService_Login_CredencialesInvalidas: contrasena equivocada -> error.
func TestAuthService_Login_CredencialesInvalidas(t *testing.T) {
	svc := service.NuevoAuthService(nuevoUsuarioRepoFake())

	_, err := svc.Registrar("ana@uleam.edu.ec", "secreta123")
	require.NoError(t, err)

	_, err = svc.Login("ana@uleam.edu.ec", "contrasenaIncorrecta")
	require.ErrorIs(t, err, service.ErrCredencialesInvalidas)
}

// TestAuthService_TokenRoundTrip recorre el flujo completo: registrar -> login
// (genera el JWT) -> ValidarToken (lo verifica). Es el corazon de la auth y se
// prueba SIN levantar un servidor: solo el servicio.
func TestAuthService_TokenRoundTrip(t *testing.T) {
	svc := service.NuevoAuthService(nuevoUsuarioRepoFake())

	creado, err := svc.Registrar("docente@uleam.edu.ec", "secreta123")
	require.NoError(t, err)

	token, err := svc.Login("docente@uleam.edu.ec", "secreta123")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	uid, err := svc.ValidarToken(token)
	require.NoError(t, err)
	assert.Equal(t, creado.ID, uid, "el token debe portar el ID del usuario que inicio sesion")
}
