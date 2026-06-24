package auth

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	secretoPorDefecto  = "cafeteria-uleam-secreto-solo-dev"
	duracionPorDefecto = 24 * time.Hour
)

type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

// Service concentra hashing (bcrypt) y JWT. Secreto y duracion son configurables
// por Options (mismo patron que en la version por capas).
type Service struct {
	repo     Repository
	secreto  []byte
	duracion time.Duration
}

type Option func(*Service)

func WithSecreto(secreto []byte) Option {
	return func(s *Service) {
		if len(secreto) > 0 {
			s.secreto = secreto
		}
	}
}

func WithDuracion(d time.Duration) Option {
	return func(s *Service) {
		if d > 0 {
			s.duracion = d
		}
	}
}

func NuevoService(repo Repository, opts ...Option) *Service {
	s := &Service{
		repo:     repo,
		secreto:  []byte(secretoPorDefecto),
		duracion: duracionPorDefecto,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) Registrar(email, password string) (Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return Usuario{}, ErrCredencialesInvalidas
	}
	if _, existe := s.repo.BuscarPorEmail(email); existe {
		return Usuario{}, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Usuario{}, err
	}
	return s.repo.Crear(Usuario{Email: email, PasswordHash: string(hash)})
}

func (s *Service) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, existe := s.repo.BuscarPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}
	return s.generarToken(u)
}

func (s *Service) generarToken(u Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duracion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secreto)
}

// ValidarToken satisface middleware.Validador: el slice expone exactamente el
// metodo que el middleware necesita, sin que el middleware lo importe.
func (s *Service) ValidarToken(tokenStr string) (int, error) {
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
