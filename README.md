# Cafetería ULEAM — Variante: Monolito Modular

Esqueleto **ejecutable** con dos módulos (`producto` y `auth`) para mostrar el
patrón: un solo binario, pero por dentro dividido en módulos con **fronteras
públicas explícitas**. Es hacia dónde escala el proyecto: límites tipo
microservicio, sin la complejidad operativa de microservicios.

## La idea central: la frontera la fuerza el lenguaje

Cada módulo expone una **fachada pública** (`producto.go`, `auth.go`) y esconde
sus tripas en un subdirectorio `internal/`. **Go garantiza** que el código bajo
`modules/producto/internal/` solo lo puede importar `modules/producto/...`. Es
decir: el módulo `auth` **no puede** —ni por accidente— importar el modelo o el
service interno de `producto`. La frontera no es una convención: es una regla del
compilador.

```
internal/
  modules/
    producto/
      producto.go          ← FACHADA pública: Nuevo(db), Rutas(r)
      internal/            ← tripas privadas (modelo, repo, service, handler)
    auth/
      auth.go              ← FACHADA pública: Nuevo(...), RegistrarRutas(r), ValidarToken(...)
      internal/            ← tripas privadas (Usuario, JWT, bcrypt)
  platform/                ← infraestructura compartida (sin lógica de negocio)
    config/
    db/                    ← abre la conexión; NO conoce los modelos
    middleware/            ← CORS + Auth (vía interfaz Validador)
```

## Dos consecuencias que se ven en el código

1. **Cada módulo es dueño de su esquema.** `platform/db` solo abre la conexión;
   no importa ningún modelo. Es cada módulo el que migra su tabla en su `Nuevo`
   (ver `interno.Migrar`). Si quitas un módulo, su tabla se va con él.

2. **Los módulos se comunican por contratos públicos, nunca por tripas.** El
   middleware de la plataforma necesita validar tokens. No importa el `Service`
   interno de auth (no podría); usa el método público `ValidarToken` de la fachada
   `auth.Modulo`, que satisface la interfaz `middleware.Validador`. Si mañana
   `producto` necesitara algo de `auth`, dependería igual de esa API pública.

## Cómo se agregaría otro módulo (p. ej. `pedido`)

1. `internal/modules/pedido/internal/` con sus tripas (modelo, repo, service, handler).
2. `internal/modules/pedido/pedido.go` con la fachada `Nuevo(db)` + `Rutas(r)`.
3. Si `pedido` necesita validar productos, **no** importa las tripas de producto:
   se le da a `producto` un método público en su fachada (p. ej.
   `ProductoExiste(id) bool`) y `pedido` depende de ese contrato.
4. Una línea en `main`: `pedidoMod := pedido.Nuevo(gdb)` + montar sus rutas.

## Correr

```bash
cp .env.example .env
go mod tidy
go build ./... && go vet ./...
go run ./cmd/cafeteria-api
# POST /api/v1/auth/register -> POST /api/v1/auth/login -> usar token en /api/v1/productos
```
