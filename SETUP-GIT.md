# Configuración del repositorio: refactor + 3 variantes arquitectónicas

Objetivo: **un solo repo** con el proyecto refactorizado en `main` y las tres
variantes arquitectónicas como **ramas** (para explorar y comparar) + **tags**
(para congelar cada estado de referencia).

Se asume que tienes las 4 carpetas de proyecto:

```
cafeteria-build-mejorado/      <- va a main
cafeteria-vertical-slice/      <- va a la rama arch/vertical-slice
cafeteria-hexagonal/           <- va a la rama arch/hexagonal
cafeteria-modular-monolith/    <- va a la rama arch/modular-monolith
```

## 1. Inicializar `main` con el proyecto refactorizado

```bash
cd cafeteria-build-mejorado
git init
git add .
git commit -m "S12: refactor (config/.env, Options, graceful shutdown, Deps, Factory, DRY)"
git branch -M main
# Si vas a usar GitHub:
# git remote add origin git@github.com:<usuario>/cafeteria-uleam-api.git
# git push -u origin main
git tag -a v2.0-refactor -m "Refactor S12 sobre arquitectura en capas"
```

## 2. Crear cada rama de variante

El patrón es el mismo para las tres. Se parte de `main`, se **reemplaza** el
contenido por el de la variante y se commitea en su rama. Ejemplo con vertical slice:

```bash
# Desde la raíz del repo (main), crea la rama
git checkout main
git checkout -b arch/vertical-slice

# Reemplaza el código por el de la variante (conserva la carpeta .git)
# (ejecuta desde DENTRO del repo; ajusta la ruta de origen a tu disco)
rm -rf cmd internal
cp -r ../cafeteria-vertical-slice/cmd ../cafeteria-vertical-slice/internal .
cp ../cafeteria-vertical-slice/go.mod ../cafeteria-vertical-slice/README.md .

git add -A
git commit -m "Variante arquitectónica: Vertical Slice (package-by-feature)"
git tag -a arch-vertical-slice-v1 -m "Referencia: Vertical Slice"

# Volver a main para la siguiente
git checkout main
```

Repite cambiando los nombres:

| Rama | Carpeta de origen | Tag |
| --- | --- | --- |
| `arch/vertical-slice` | `cafeteria-vertical-slice/` | `arch-vertical-slice-v1` |
| `arch/hexagonal` | `cafeteria-hexagonal/` | `arch-hexagonal-v1` |
| `arch/modular-monolith` | `cafeteria-modular-monolith/` | `arch-modular-monolith-v1` |

## 3. Subir todo a GitHub (ramas + tags)

```bash
git push origin --all     # sube main y las 3 ramas arch/*
git push origin --tags    # sube los 4 tags
```

## 4. El truco de enseñanza: ver el diff de cada variante

En GitHub, la vista **Compare** muestra *qué cambia* al reorganizar a un patrón:

```
https://github.com/<usuario>/cafeteria-uleam-api/compare/main...arch/vertical-slice
```

Localmente es lo mismo:

```bash
git diff --stat main arch/vertical-slice     # qué archivos se movieron/cambiaron
git diff main arch/hexagonal                 # el diff completo
```

Esto deja ver, por ejemplo, cómo `internal/handlers/producto.go` +
`internal/service/producto.go` + `internal/storage/...` se funden en un solo
paquete `internal/producto/` en la variante vertical slice.

## Resumen

- **Ramas** `arch/*` → explorables, comparables con `compare`/`git diff`.
- **Tags** `*-v1` / `v2.0-refactor` → marcadores inmutables de cada referencia.
- **Un repo**, no cuatro.
