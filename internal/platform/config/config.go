package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Puerto      string
	RutaDB      string
	JWTSecreto  []byte
	JWTDuracion time.Duration
}

func Cargar() Config {
	_ = godotenv.Load()
	dur := 24 * time.Hour
	if v := os.Getenv("JWT_DURACION"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			dur = d
		}
	}
	return Config{
		Puerto:      conTexto("PUERTO", ":8080"),
		RutaDB:      conTexto("RUTA_DB", "cafeteria.db"),
		JWTSecreto:  []byte(conTexto("JWT_SECRETO", "cafeteria-uleam-secreto-solo-dev")),
		JWTDuracion: dur,
	}
}

func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}
