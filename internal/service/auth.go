package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/storage"
)

// Valores por defecto para desarrollo. En produccion se inyectan via Options
// (WithSecreto / WithDuracionToken) desde la configuracion cargada del .env.
const (
	secretoPorDefecto  = "cafeteria-uleam-secreto-solo-dev"
	duracionPorDefecto = 24 * time.Hour
)

// Claims es el contenido del JWT: el ID del usuario + los campos estandar (exp, iat).
type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

// AuthService concentra TODA la logica de autenticacion: hashing de contrasenas
// (bcrypt) y generacion/validacion de JWT. El handler y el middleware no saben
// de bcrypt ni de firmas: solo llaman a este servicio.
//
// Antes el secreto y la duracion eran una VARIABLE y una CONSTANTE globales del
// paquete. Eso impedia configurarlos por entorno y dificultaba testear con otro
// secreto. Ahora son campos del struct, con defaults seguros, configurables por
// el patron funcional Options.
type AuthService struct {
	repo     storage.UserRepository
	secreto  []byte
	duracion time.Duration
}

// AuthOption configura un AuthService en su construccion (patron Options).
type AuthOption func(*AuthService)

// WithSecreto inyecta la clave de firma del JWT (desde config/.env en produccion).
// Si recibe un secreto vacio, conserva el valor por defecto.
func WithSecreto(secreto []byte) AuthOption {
	return func(a *AuthService) {
		if len(secreto) > 0 {
			a.secreto = secreto
		}
	}
}

// WithDuracionToken inyecta la validez del token. Si recibe <= 0, conserva el default.
func WithDuracionToken(d time.Duration) AuthOption {
	return func(a *AuthService) {
		if d > 0 {
			a.duracion = d
		}
	}
}

// NuevoAuthService construye el servicio con defaults de desarrollo y aplica las
// Options recibidas. Como opts es variadico, las llamadas antiguas
// NuevoAuthService(repo) siguen compilando: simplemente reciben los defaults.
func NuevoAuthService(repo storage.UserRepository, opts ...AuthOption) *AuthService {
	s := &AuthService{
		repo:     repo,
		secreto:  []byte(secretoPorDefecto),
		duracion: duracionPorDefecto,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Registrar crea un usuario nuevo con la contrasena hasheada (bcrypt).
func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrCredencialesInvalidas
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
	})
}

// Login verifica las credenciales y, si son validas, devuelve un JWT firmado.
func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}

	return s.generarToken(u)
}

// generarToken arma y firma el JWT con el ID del usuario y la expiracion.
// Usa el secreto y la duracion del servicio (inyectables por Options).
func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duracion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secreto)
}

// ValidarToken verifica firma y expiracion, y devuelve el ID del usuario.
// Lo usa el middleware de autenticacion: el JWT vive aqui, no en el middleware.
func (s *AuthService) ValidarToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return s.secreto, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrCredencialesInvalidas
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidas
	}
	return claims.UsuarioID, nil
}
