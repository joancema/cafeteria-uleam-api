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

> El secreto del JWT ya **no** está hardcodeado: se carga desde el `.env`
> (`JWT_SECRETO`) vía `internal/config`. Con un default seguro solo para dev.

## Correr

```bash
cp .env.example .env   # opcional: ajusta puerto, secreto JWT, backend...
go mod tidy            # resuelve godotenv, golang-jwt y golang.org/x/crypto
go build ./... && go vet ./...
go run ./cmd/cafeteria-api      # backend GORM (por defecto)
STORAGE=sqlc go run ./cmd/cafeteria-api   # backend sqlc para productos/categorias
```

Para detener el servidor, `Ctrl+C`: hace un **cierre ordenado** (graceful
shutdown) que termina las peticiones en curso y cierra la base de datos.

Servidor en `http://localhost:8080`. Los usuarios viven siempre en GORM, incluso
con `STORAGE=sqlc`.

## Dependencias nuevas

- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`

---

## Testing (Semana 11)

La capa de tests es **aditiva**: no cambia una sola línea del código de S10, solo
agrega archivos `_test.go`. La estrategia es la de un proyecto real, por capa:

| Archivo | Qué prueba | Cómo |
| --- | --- | --- |
| `internal/service/producto_test.go` | Regla de negocio de productos (Crear/Obtener/Actualizar/Borrar/Listar) aislada de la base | **Mock** de la interfaz estrecha + table-driven |
| `internal/service/categoria_test.go` | Regla de negocio de categorías | **Mock** |
| `internal/service/auth_test.go` | Registro, login y validación de JWT (bcrypt + token) | Repo *fake* en memoria |
| `internal/handlers/producto_test.go` | Endpoints de productos por HTTP, con el **middleware real** (200/201/204/400/404 y el **401**) | `httptest` |
| `internal/handlers/auth_test.go` | Endpoints de auth (register 201/409/400, login 200/401) | `httptest` |
| `internal/handlers/helpers_test.go` | Dobles en memoria + router idéntico a `main` | (infra de tests) |
| `internal/storage/sqlite_test.go` | El repositorio **GORM real** sobre SQLite `:memory:` + el índice único de email | Librería **estándar** (`testing`) |

**La regla clave del diseño:** los tests de *service* y *handler* **no tocan la
base de datos** —usan dobles para ser rápidos y aislados—; la base de datos real
(GORM) se prueba **aparte**, en `sqlite_test.go`, contra una SQLite desechable.

### Comandos

```bash
go test ./...                 # toda la suite
go test ./... -v              # con el nombre de cada caso
go test ./... -cover          # cobertura por paquete
go test ./internal/handlers/ -coverpkg=./internal/... -cover   # incluye el middleware
```

> Dependencia nueva de esta semana: `github.com/stretchr/testify`.

---

## Refactorizaciones (Semana 12) — Principios de arquitectura y patrones

Esta versión **no agrega features**: mejora *mantenibilidad* y *ciclo de vida*,
sin tocar la lógica de negocio ni los tests de S11 (salvo un call-site).

| # | Refactor | Patrón | Archivo(s) | Qué mejora |
| --- | --- | --- | --- | --- |
| 1 | Configuración por `.env` | Config centralizada + godotenv | `internal/config/config.go`, `.env.example` | Una sola fuente de verdad; nada hardcodeado |
| 2 | Secreto/duración del JWT inyectables | **Options** (funcional) | `internal/service/auth.go` | Quita estado global; configurable y testeable |
| 3 | Servidor con timeouts + apagado limpio | **Options** + graceful shutdown | `internal/httpserver/httpserver.go`, `cmd/cafeteria-api/main.go` | Seguridad en prod; cierre ordenado |
| 4 | `NewServer` con struct de dependencias | Parameter object (Deps) | `internal/handlers/server.go` | El constructor escala sin crecer la firma |
| 5 | Selección de backend encapsulada | **Factory** | `internal/storage/factory.go` | `main` deja de hacer plumbing de DB |
| 6 | Extracción de `id` de la URL | DRY (helper) | `internal/handlers/params.go` | Elimina 6 repeticiones idénticas |

> `main()` quedó delgado: delega en `run(cfg) error`. La construcción de
> dependencias y el ciclo de vida son ahora explícitos y testeables.

### Patrones arquitectónicos de referencia (no aplicados aquí)

El código sigue siendo **por capas**. Como referencia para el proyecto, hay tres
variantes en ramas del repositorio (ver `SETUP-GIT.md`):

- `arch/vertical-slice` — un paquete por entidad (todas sus capas juntas).
- `arch/hexagonal` — ports & adapters (el dominio define los puertos).
- `arch/modular-monolith` — módulos con frontera pública explícita.

### Dependencia nueva de esta semana

- `github.com/joho/godotenv` — carga el archivo `.env` en desarrollo.
