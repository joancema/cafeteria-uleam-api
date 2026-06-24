// Command cafeteria-api arranca la API en su variante HEXAGONAL (ports & adapters).
//
// El cableado deja ver las tres piezas: adaptador de salida -> núcleo ->
// adaptador de entrada. Las dependencias apuntan hacia el núcleo.
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

	"cafeteria-hex/internal/adaptadores/entrada/rest"
	"cafeteria-hex/internal/adaptadores/salida/persistencia"
	"cafeteria-hex/internal/config"
	"cafeteria-hex/internal/core/producto"
)

func main() {
	if err := run(config.Cargar()); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	gdb, cerrar, err := persistencia.Abrir(cfg.RutaDB)
	if err != nil {
		return err
	}
	defer func() { _ = cerrar() }()
	persistencia.Sembrar(gdb)

	repo := persistencia.NuevoProductoGORM(gdb) // adaptador de SALIDA (implementa el puerto)
	svc := producto.NuevoServicio(repo)         // NÚCLEO (recibe el puerto de salida)
	h := rest.NuevoProductoHandler(svc)         // adaptador de ENTRADA (recibe el puerto de entrada)

	r := chi.NewRouter()
	r.Route("/api/v1/productos", h.Rutas)

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
		log.Printf("Servidor (hexagonal) escuchando en http://localhost%s", cfg.Puerto)
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
