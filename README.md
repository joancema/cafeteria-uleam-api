# Cafetería ULEAM — Variante: Hexagonal (Ports & Adapters)

Esqueleto **ejecutable** con la entidad `producto` cableada de punta a punta para
mostrar el patrón. `categoria` y `auth` se omiten a propósito (ver más abajo cómo
se agregarían).

## La regla del hexagonal

Las dependencias apuntan **hacia el núcleo**. El núcleo no sabe que existen GORM
ni HTTP.

```
   Adaptador ENTRADA            NÚCLEO                 Adaptador SALIDA
   (rest.ProductoHandler) ─▶ producto.Servicio        producto.Repositorio ◀─ (persistencia.ProductoGORM)
                              (puerto entrada)          (puerto salida)
                                    │                        ▲
                                    └── producto.servicio ───┘
                                        (lógica de negocio)
```

## Estructura

```
internal/
  core/producto/             ← NÚCLEO (sin infraestructura)
    producto.go              ←   entidad + errores
    puertos.go               ←   Repositorio (salida) + Servicio (entrada)  [INTERFACES]
    servicio.go              ←   implementación de la lógica (privada)
  adaptadores/
    entrada/rest/            ← adaptador de ENTRADA: depende de producto.Servicio
    salida/persistencia/     ← adaptador de SALIDA: implementa producto.Repositorio (GORM)
  config/
```

## El punto clave

En la versión por capas, la interfaz `Repository` vivía en el paquete `storage`
(lado infraestructura). En hexagonal, **el puerto lo define el núcleo**
(`core/producto/puertos.go`) y la infraestructura lo *implementa*. Esa inversión
es lo que permite cambiar de GORM a otra base —o de HTTP a gRPC— sin tocar el
núcleo. El proyecto ya lo demostraba con sus 3 backends; aquí se formaliza el
principio.

`NuevoServicio` devuelve la **interfaz** `Servicio`, no el struct: el exterior
nunca ve la implementación concreta.

## Cómo se agregarían categoria y auth

- `categoria`: otro paquete `core/categoria` con su entidad + puertos + servicio,
  un `persistencia.CategoriaGORM` y un handler en `rest`.
- `auth`: `core/auth` define un puerto `Servicio` con `ValidarToken`; el
  middleware de autenticación es **otro adaptador de entrada** que depende de ese
  puerto. El JWT/bcrypt viven en el núcleo o en un adaptador de salida de
  criptografía, según qué tan estricto se quiera ser.

## Correr

```bash
cp .env.example .env
go mod tidy
go build ./... && go vet ./...
go run ./cmd/cafeteria-api    # GET http://localhost:8080/api/v1/productos
```
