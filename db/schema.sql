CREATE TABLE categorias (
    id          INTEGER PRIMARY KEY,
    nombre      TEXT    NOT NULL,
    descripcion TEXT    NOT NULL
);

CREATE TABLE productos (
    id           INTEGER PRIMARY KEY,
    nombre       TEXT    NOT NULL,
    precio       REAL    NOT NULL,
    stock        INTEGER NOT NULL,
    categoria_id INTEGER NOT NULL
);
