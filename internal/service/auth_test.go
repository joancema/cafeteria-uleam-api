package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

// usuarioRepoFake: repositorio de usuarios EN MEMORIA para los tests.
// En produccion el AuthService recibe el repositorio GORM; aqui le inyectamos
// este doble. El servicio no nota la diferencia: depende de la interfaz.
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

// --- Registrar ---

func TestAuthService_Registrar(t *testing.T) {
	t.Run("campos vacios -> credenciales invalidas", func(t *testing.T) {
		svc := service.NuevoAuthService(nuevoUsuarioRepoFake())
		_, err := svc.Registrar("  ", "")
		require.ErrorIs(t, err, service.ErrCredencialesInvalidas)
	})
	t.Run("email duplicado -> ErrEmailEnUso", func(t *testing.T) {
		svc := service.NuevoAuthService(nuevoUsuarioRepoFake())
		_, err := svc.Registrar("ana@uleam.edu.ec", "secreta123")
		require.NoError(t, err)
		_, err = svc.Registrar("ana@uleam.edu.ec", "otra456")
		require.ErrorIs(t, err, service.ErrEmailEnUso)
	})
	t.Run("valido guarda hash, no la contrasena en claro", func(t *testing.T) {
		svc := service.NuevoAuthService(nuevoUsuarioRepoFake())
		u, err := svc.Registrar("ana@uleam.edu.ec", "secreta123")
		require.NoError(t, err)
		assert.NotEmpty(t, u.PasswordHash)
		assert.NotEqual(t, "secreta123", u.PasswordHash, "el hash bcrypt jamas debe ser la contrasena en claro")
	})
}

// --- Login ---

func TestAuthService_Login(t *testing.T) {
	prepararUsuario := func() *service.AuthService {
		svc := service.NuevoAuthService(nuevoUsuarioRepoFake())
		_, _ = svc.Registrar("ana@uleam.edu.ec", "secreta123")
		return svc
	}
	t.Run("usuario inexistente -> credenciales invalidas", func(t *testing.T) {
		svc := service.NuevoAuthService(nuevoUsuarioRepoFake())
		_, err := svc.Login("nadie@uleam.edu.ec", "x")
		require.ErrorIs(t, err, service.ErrCredencialesInvalidas)
	})
	t.Run("contrasena incorrecta -> credenciales invalidas", func(t *testing.T) {
		_, err := prepararUsuario().Login("ana@uleam.edu.ec", "incorrecta")
		require.ErrorIs(t, err, service.ErrCredencialesInvalidas)
	})
	t.Run("correcto devuelve token", func(t *testing.T) {
		token, err := prepararUsuario().Login("ana@uleam.edu.ec", "secreta123")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

// --- ValidarToken ---

func TestAuthService_ValidarToken(t *testing.T) {
	t.Run("token basura -> error", func(t *testing.T) {
		svc := service.NuevoAuthService(nuevoUsuarioRepoFake())
		_, err := svc.ValidarToken("esto.no.es-un-jwt")
		require.Error(t, err)
	})
}

// TestAuthService_RoundTrip recorre el flujo completo sin base de datos:
// registrar (bcrypt) -> login (firma el JWT) -> validar (lo verifica).
func TestAuthService_RoundTrip(t *testing.T) {
	svc := service.NuevoAuthService(nuevoUsuarioRepoFake())

	creado, err := svc.Registrar("docente@uleam.edu.ec", "secreta123")
	require.NoError(t, err)

	token, err := svc.Login("docente@uleam.edu.ec", "secreta123")
	require.NoError(t, err)

	uid, err := svc.ValidarToken(token)
	require.NoError(t, err)
	assert.Equal(t, creado.ID, uid, "el token debe portar el ID del usuario")
}
