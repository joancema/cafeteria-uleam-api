// Command cafeteria-api arranca la API en su variante VERTICAL SLICE.
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

	"cafeteria-uleam-api/internal/auth"
	"cafeteria-uleam-api/internal/categoria"
	"cafeteria-uleam-api/internal/platform/config"
	"cafeteria-uleam-api/internal/platform/httpserver"
	"cafeteria-uleam-api/internal/platform/middleware"
	"cafeteria-uleam-api/internal/platform/storage"
	"cafeteria-uleam-api/internal/producto"
)

func main() {
	if err := run(config.Cargar()); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	gdb, cerrar, err := storage.Abrir(cfg.RutaDB)
	if err != nil {
		return err
	}
	defer func() { _ = cerrar() }()
	storage.Sembrar(gdb)

	// Cada slice arma su propia vertical: repo GORM -> service -> handler.
	prodH := producto.NuevoHandler(producto.NuevoService(producto.NuevoRepoGORM(gdb)))
	catH := categoria.NuevoHandler(categoria.NuevoService(categoria.NuevoRepoGORM(gdb)))
	authSvc := auth.NuevoService(
		auth.NuevoRepoGORM(gdb),
		auth.WithSecreto(cfg.JWTSecreto),
		auth.WithDuracion(cfg.JWTDuracion),
	)
	authH := auth.NuevoHandler(authSvc)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		// Publicas.
		r.Post("/auth/register", authH.Registrar)
		r.Post("/auth/login", authH.Login)
		// Protegidas: cada slice monta SUS rutas con r.Route(prefijo, h.Rutas).
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Route("/productos", prodH.Rutas)
			r.Route("/categorias", catH.Rutas)
		})
	})

	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor (vertical slice) escuchando en http://localhost%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Senal de apagado recibida, cerrando ordenadamente...")
	}

	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	return srv.Shutdown(ctxApagado)
}
