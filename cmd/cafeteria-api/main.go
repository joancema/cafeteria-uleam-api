// Command cafeteria-api arranca la API en su variante MONOLITO MODULAR.
//
// main solo conoce las FACHADAS públicas de cada módulo (auth, producto) y las
// cablea. Ningún módulo puede tocar las tripas de otro: Go lo impide vía internal/.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"cafeteria-mod/internal/modules/auth"
	"cafeteria-mod/internal/modules/producto"
	"cafeteria-mod/internal/platform/config"
	"cafeteria-mod/internal/platform/db"
	"cafeteria-mod/internal/platform/middleware"
)

func main() {
	if err := run(config.Cargar()); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	gdb, cerrar, err := db.Abrir(cfg.RutaDB)
	if err != nil {
		return err
	}
	defer func() { _ = cerrar() }()

	// Cada módulo se ensambla solo (migra su esquema y arma sus capas).
	authMod := auth.Nuevo(gdb, cfg.JWTSecreto, cfg.JWTDuracion)
	prodMod := producto.Nuevo(gdb)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		authMod.RegistrarRutas(r) // públicas: /auth/register, /auth/login
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authMod)) // la fachada de auth ES el Validador
			r.Route("/productos", prodMod.Rutas)
		})
	})

	srv := &http.Server{
		Addr:         cfg.Puerto,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor (monolito modular) escuchando en http://localhost%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()
	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Apagado ordenado...")
	}
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	return srv.Shutdown(ctxApagado)
}
