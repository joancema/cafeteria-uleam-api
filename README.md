# Cafetería ULEAM API — Arquitectura en capas + Autenticación JWT

API REST de la cafetería universitaria. Esta versión formaliza la **arquitectura
en 3 capas** (Handler → Service → Repository) e incorpora el **módulo de
autenticación** (registro, login y protección de rutas con JWT) como ejemplo
canónico de por qué separar capas.

## Capas

```
HTTP  ─▶  Handler  ─▶  Service  ─▶  Repository  ─▶  (GORM / sqlc / Memoria)
            (HTTP)     (negocio)     (persistencia)
```

- **Handler** (`internal/handlers`): decodifica el request, llama al service y
  traduce el resultado a HTTP. Sin lógica de negocio.
- **Service** (`internal/service`): reglas de negocio y validaciones. Devuelve
  *errores de dominio* (`ErrNombreVacio`, `ErrNoEncontrado`, `ErrEmailEnUso`...).
- **Repository** (`internal/storage`): persistencia. La interfaz `Almacen` se
  partió en interfaces por entidad (`ProductoRepository`, `CategoriaRepository`)
  recompuestas por *embedding*; `UserRepository` para usuarios.

Los tres backends previos (`Memoria`, `AlmacenSQLite` con GORM, `AlmacenSQLC`)
siguen cumpliendo `Almacen` sin cambios: el mismo objeto satisface las interfaces
estrechas que consume cada service.

## Autenticación

- `POST /api/v1/auth/register` — crea usuario (contraseña hasheada con bcrypt).
- `POST /api/v1/auth/login` — verifica credenciales y devuelve `{"token": "..."}`.
- El resto de rutas (`/productos`, `/categorias`) exige el header
  `Authorization: Bearer <token>`. El middleware delega la validación al
  `AuthService` (el JWT vive en el service, no en el middleware).

> El secreto del JWT está hardcodeado con un `TODO(S12)`: en producción va en
> una variable de entorno.

## Correr

```bash
go mod tidy            # resuelve golang-jwt y golang.org/x/crypto
go build ./... && go vet ./...
go run ./cmd/cafeteria-api      # backend GORM (por defecto)
STORAGE=sqlc go run ./cmd/cafeteria-api   # backend sqlc para productos/categorias
```

Servidor en `http://localhost:8080`. Los usuarios viven siempre en GORM, incluso
con `STORAGE=sqlc`.

## Dependencias nuevas

- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`
