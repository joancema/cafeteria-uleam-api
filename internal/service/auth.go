package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/storage"
)

// secretoJWT firma y verifica los tokens.
// TODO(S12): mover a una variable de entorno; jamas dejarlo en el codigo en produccion.
var secretoJWT = []byte("cafeteria-uleam-secreto-demo-cambiar-en-S12")

// duracionToken es la validez del token desde su emision.
const duracionToken = 24 * time.Hour

// Claims es el contenido del JWT: el ID del usuario + los campos estandar (exp, iat).
type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

// AuthService concentra TODA la logica de autenticacion: hashing de contrasenas
// (bcrypt) y generacion/validacion de JWT. El handler y el middleware no saben
// de bcrypt ni de firmas: solo llaman a este servicio. Esa es la razon de ser
// de la capa de servicio.
type AuthService struct {
	repo storage.UserRepository
}

func NuevoAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
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
func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWT)
}

// ValidarToken verifica firma y expiracion, y devuelve el ID del usuario.
// Lo usa el middleware de autenticacion: el JWT vive aqui, no en el middleware.
func (s *AuthService) ValidarToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretoJWT, nil
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
