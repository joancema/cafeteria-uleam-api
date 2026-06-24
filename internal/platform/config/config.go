// Package config carga la configuracion desde el entorno (con .env opcional).
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Puerto       string
	RutaDB       string
	JWTSecreto   []byte
	JWTDuracion  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Cargar() Config {
	_ = godotenv.Load()
	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		RutaDB:       conTexto("RUTA_DB", "cafeteria.db"),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "cafeteria-uleam-secreto-solo-dev")),
		JWTDuracion:  conDuracion("JWT_DURACION", 24*time.Hour),
		ReadTimeout:  conDuracion("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: conDuracion("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}

func conDuracion(clave string, porDefecto time.Duration) time.Duration {
	if v := os.Getenv(clave); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return porDefecto
}
