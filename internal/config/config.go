// Package config carga la configuracion de la aplicacion desde variables de
// entorno (con soporte para un archivo .env opcional) y expone valores por
// defecto razonables para desarrollo.
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa toda la configuracion del servidor en un solo lugar.
//
// Antes estos valores estaban dispersos y hardcodeados: el secreto JWT vivia en
// una var global de service/auth.go, el puerto y la ruta de la DB eran literales
// en main.go, y el backend se leia con os.Getenv suelto. Ahora hay UNA sola
// fuente de verdad.
type Config struct {
	Puerto       string        // puerto HTTP, ej ":8080"
	RutaDB       string        // archivo SQLite, ej "cafeteria.db"
	Backend      string        // backend de productos/categorias: "gorm" (default) o "sqlc"
	JWTSecreto   []byte        // clave para firmar/verificar JWT
	JWTDuracion  time.Duration // validez del token
	ReadTimeout  time.Duration // timeout de lectura del servidor HTTP
	WriteTimeout time.Duration // timeout de escritura del servidor HTTP
}

// Cargar lee la configuracion. Primero intenta cargar un archivo .env (si no
// existe, no es un error: en produccion las variables vienen del entorno real).
// Luego lee cada variable con un valor por defecto seguro para desarrollo.
func Cargar() Config {
	// godotenv.Load NO sobreescribe variables ya presentes en el entorno; si no
	// hay archivo .env, devuelve un error que ignoramos a proposito.
	_ = godotenv.Load()

	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		RutaDB:       conTexto("RUTA_DB", "cafeteria.db"),
		Backend:      conTexto("STORAGE", "gorm"),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "cafeteria-uleam-secreto-solo-dev")),
		JWTDuracion:  conDuracion("JWT_DURACION", 24*time.Hour),
		ReadTimeout:  conDuracion("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: conDuracion("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

// conTexto devuelve la variable de entorno o el valor por defecto si esta vacia.
func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}

// conDuracion parsea una duracion (ej "24h", "30m") o usa el valor por defecto.
func conDuracion(clave string, porDefecto time.Duration) time.Duration {
	v := os.Getenv(clave)
	if v == "" {
		return porDefecto
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return porDefecto
	}
	return d
}
