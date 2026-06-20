-- ===================== PRODUCTOS =====================

-- name: ListarProductos :many
SELECT id, nombre, precio, stock, categoria_id FROM productos;

-- name: BuscarProductoPorID :one
SELECT id, nombre, precio, stock, categoria_id FROM productos
WHERE id = ?;

-- name: CrearProducto :one
INSERT INTO productos (nombre, precio, stock, categoria_id)
VALUES (?, ?, ?, ?)
RETURNING id, nombre, precio, stock, categoria_id;

-- name: ActualizarProducto :one
UPDATE productos
SET nombre = ?, precio = ?, stock = ?, categoria_id = ?
WHERE id = ?
RETURNING id, nombre, precio, stock, categoria_id;

-- name: BorrarProducto :execrows
DELETE FROM productos WHERE id = ?;

-- ===================== CATEGORIAS =====================

-- name: ListarCategorias :many
SELECT id, nombre, descripcion FROM categorias;

-- name: BuscarCategoriaPorID :one
SELECT id, nombre, descripcion FROM categorias
WHERE id = ?;

-- name: CrearCategoria :one
INSERT INTO categorias (nombre, descripcion)
VALUES (?, ?)
RETURNING id, nombre, descripcion;

-- name: ActualizarCategoria :one
UPDATE categorias
SET nombre = ?, descripcion = ?
WHERE id = ?
RETURNING id, nombre, descripcion;

-- name: BorrarCategoria :execrows
DELETE FROM categorias WHERE id = ?;
