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

---

## Testing (Semana 11)

La capa de tests es **aditiva**: no cambia ni una línea del código de S10, solo
agrega archivos `_test.go`. Cubre **tres estilos** que conviene conocer:

| Archivo | Qué prueba | Estilo |
| --- | --- | --- |
| `internal/service/producto_test.go` | La regla de negocio (`validarProducto`, `Crear`) aislada de la BD | **Mock** con `testify/mock` + **table-driven** |
| `internal/service/auth_test.go` | Registro, login y round-trip de JWT (bcrypt + token) | Repo *fake* en memoria |
| `internal/handlers/producto_test.go` | Los endpoints a través del **router y `middleware.Auth` reales** (201/404/400/**401 sin token**) | `httptest` |
| `internal/storage/memoria_test.go` | El almacén en memoria (CRUD + comma-ok) | Librería **estándar** (`testing`, sin testify) |

### Comandos

```bash
go test ./...                 # corre todo
go test ./... -v              # con el nombre de cada test
go test ./... -cover          # porcentaje por paquete
go test ./internal/handlers/ -coverpkg=./internal/... -cover   # incluye el middleware
```

### Filosofía de cobertura

No perseguimos el 100%. La **lógica de negocio** (`validarProducto`, `Crear`)
está al **100%** porque es donde un error duele; los *getters* triviales y el
wiring valen menos un test. **Cubrir lo que importa, no inflar el número.**

> Dependencia nueva de esta semana: `github.com/stretchr/testify`.
