# Cafetería ULEAM — Variante: Vertical Slice (package-by-feature)

Misma API, **reorganizada por funcionalidad** en lugar de por capa técnica. Es la
variante recomendada cuando el equipo trabaja en paralelo: **un slice = una
carpeta = un dueño**.

## Estructura

```
internal/
  platform/            ← lo TRANSVERSAL (no es de ninguna entidad)
    config/            ← carga de .env
    web/               ← RespondJSON, RespondError, IDDeURL (sin dominio)
    httpserver/        ← *http.Server con Options
    middleware/        ← CORS + Auth (vía interfaz Validador)
    storage/           ← abre DB, migra y siembra (importa los modelos de los slices)
  producto/            ← VERTICAL COMPLETA de productos:
    producto.go        ←   modelo
    errores.go         ←   errores de dominio del slice
    repository.go      ←   contrato (puerto)
    gorm.go            ←   adaptador GORM (implementación)
    service.go         ←   lógica de negocio
    handler.go         ←   HTTP + mapeo de errores a status
    routes.go          ←   monta sus propias rutas
    service_test.go    ←   test con mock del Repository del slice
  categoria/           ← misma estructura que producto
  auth/                ← misma estructura (+ JWT/bcrypt y Options en el service)
```

## La idea

Comparado con la versión por capas (donde `producto` vivía repartido entre
`handlers/producto.go`, `service/producto.go` y `storage/...`), aquí **todo lo de
productos está en una sola carpeta**. Beneficio directo para los grupos del
proyecto: dos estudiantes editando módulos distintos **no se pisan** —editan
carpetas distintas— y desaparecen los conflictos de merge en `handlers/` y
`service/`.

### Cómo se evitan los ciclos de importación

- `platform/web` y `platform/middleware` **no conocen el dominio** → los slices
  pueden depender de ellos sin ciclo.
- `platform/storage` **importa** los slices (para migrar sus modelos), pero
  **ningún slice importa `platform/storage`** → dependencia en una sola dirección.
- El middleware de auth recibe una interfaz `Validador` (un método
  `ValidarToken`), no el paquete `auth` → no se acoplan.
- Cada slice **mapea sus propios errores** a HTTP en su `handler.go`.

## Nota de alcance

Para mantener la variante enfocada, **solo se incluye el backend GORM**. Los
backends `Memoria` y `sqlc` de la versión por capas colapsan exactamente igual:
serían un archivo más dentro de cada slice (`memoria.go`, `sqlc.go`) implementando
el mismo `Repository`. La idea del patrón ya queda demostrada con GORM.

## Correr

```bash
cp .env.example .env          # opcional
go mod tidy
go build ./... && go vet ./...
go test ./...
go run ./cmd/cafeteria-api    # http://localhost:8080 ; Ctrl+C = apagado limpio
```
