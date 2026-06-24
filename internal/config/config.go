package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Puerto string
	RutaDB string
}

func Cargar() Config {
	_ = godotenv.Load()
	return Config{
		Puerto: conTexto("PUERTO", ":8080"),
		RutaDB: conTexto("RUTA_DB", "cafeteria.db"),
	}
}

func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}
