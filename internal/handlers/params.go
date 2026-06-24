package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// idDeURL extrae y valida el parametro de ruta {id} como entero.
//
// Antes este bloque (strconv.Atoi sobre chi.URLParam) estaba repetido en CADA
// handler que recibe un id: 6 veces entre productos y categorias. Centralizarlo
// elimina la duplicacion y deja un unico lugar donde cambiar la validacion del id.
func idDeURL(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}
