// Package interno: tripas privadas del módulo auth (forzado por internal/ de Go).
package interno

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usuario struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreadoEn     time.Time `json:"creado_en"`
}

var (
	ErrEmailEnUso            = errors.New("el email ya esta registrado")
	ErrCredencialesInvalidas = errors.New("email o contrasena incorrectos")
)

func Migrar(db *gorm.DB) { db.AutoMigrate(&Usuario{}) }

type Repository interface {
	Crear(u Usuario) (Usuario, error)
	BuscarPorEmail(email string) (Usuario, bool)
}

type RepoGORM struct{ db *gorm.DB }

func NuevoRepoGORM(db *gorm.DB) *RepoGORM { return &RepoGORM{db: db} }

func (r *RepoGORM) Crear(u Usuario) (Usuario, error) {
	u.CreadoEn = time.Now()
	if err := r.db.Create(&u).Error; err != nil {
		return Usuario{}, err
	}
	return u, nil
}
func (r *RepoGORM) BuscarPorEmail(email string) (Usuario, bool) {
	var u Usuario
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return Usuario{}, false
	}
	return u, true
}

var _ Repository = (*RepoGORM)(nil)

type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

type Service struct {
	repo     Repository
	secreto  []byte
	duracion time.Duration
}

func NuevoService(repo Repository, secreto []byte, duracion time.Duration) *Service {
	if len(secreto) == 0 {
		secreto = []byte("cafeteria-uleam-secreto-solo-dev")
	}
	if duracion <= 0 {
		duracion = 24 * time.Hour
	}
	return &Service{repo: repo, secreto: secreto, duracion: duracion}
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
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duracion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secreto)
}

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

type Handler struct{ svc *Service }

func NuevoHandler(svc *Service) *Handler { return &Handler{svc: svc} }

type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Registrar(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	u, err := h.svc.Registrar(c.Email, c.Password)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, u)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var c credenciales
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	token, err := h.svc.Login(c.Email, c.Password)
	if err != nil {
		responderError(w, statusDeError(err), err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"token": token})
}

func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
func responderError(w http.ResponseWriter, status int, msg string) {
	responderJSON(w, status, map[string]string{"error": msg})
}
func statusDeError(err error) int {
	switch {
	case errors.Is(err, ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
