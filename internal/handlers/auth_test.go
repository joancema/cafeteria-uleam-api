package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	h, _, _ := construirEntorno()

	t.Run("valido -> 201", func(t *testing.T) {
		body := `{"email":"nuevo@uleam.edu.ec","password":"secreta123"}`
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", body, ""))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
	t.Run("email duplicado -> 409", func(t *testing.T) {
		body := `{"email":"dup@uleam.edu.ec","password":"secreta123"}`
		ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", body, ""))        // primero
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", body, "")) // repetido
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
	t.Run("JSON malformado -> 400", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", `{roto`, ""))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestLogin(t *testing.T) {
	h, _, _ := construirEntorno()
	cred := `{"email":"ana@uleam.edu.ec","password":"secreta123"}`
	ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", cred, ""))

	t.Run("correcto -> 200 + token", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/login", cred, ""))
		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Token string `json:"token"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.NotEmpty(t, resp.Token)
	})
	t.Run("contrasena incorrecta -> 401", func(t *testing.T) {
		malo := `{"email":"ana@uleam.edu.ec","password":"incorrecta"}`
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/login", malo, ""))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
