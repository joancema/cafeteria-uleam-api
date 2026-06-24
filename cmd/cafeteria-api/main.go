// Command cafeteria-api arranca el servidor HTTP de la Cafeteria Universitaria.
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

	"cafeteria-uleam-api/internal/config"
	"cafeteria-uleam-api/internal/handlers"
	"cafeteria-uleam-api/internal/httpserver"
	"cafeteria-uleam-api/internal/middleware"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

// main queda DELGADO: carga la configuracion, delega en run y traduce el error
// a un exit code. Toda la logica de arranque vive en run, que devuelve error en
// lugar de llamar a log.Fatal en cada paso (mas testeable y mas limpio).
func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

// run construye las dependencias, levanta el servidor y bloquea hasta recibir
// una senal de apagado (Ctrl+C / SIGTERM); en ese momento hace un cierre
// ordenado (graceful shutdown): deja de aceptar conexiones, termina las que
// estan en curso y cierra la base de datos.
func run(cfg config.Config) error {
	// 1. Recursos de almacenamiento (Factory): abre DB, migra, siembra y elige backend.
	recursos, err := storage.Inicializar(cfg.RutaDB, cfg.Backend)
	if err != nil {
		return err
	}
	defer func() { _ = recursos.Cerrar() }()
	log.Printf("Backend de productos/categorias: %s", recursos.BackendUsado)

	// 2. Capa de servicio. AuthService recibe secreto y duracion por Options,
	//    tomados de la configuracion (antes eran globales hardcodeadas).
	productoSvc := service.NuevoProductoService(recursos.Almacen)
	categoriaSvc := service.NuevoCategoriaService(recursos.Almacen)
	authSvc := service.NuevoAuthService(
		recursos.Usuarios,
		service.WithSecreto(cfg.JWTSecreto),
		service.WithDuracionToken(cfg.JWTDuracion),
	)

	// 3. Server con sus dependencias agrupadas en un struct (escala sin crecer
	//    la firma del constructor).
	servidor := handlers.NewServer(handlers.Deps{
		Productos:  productoSvc,
		Categorias: categoriaSvc,
		Auth:       authSvc,
	})

	// 4. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 5. Rutas versionadas /api/v1/.
	r.Route("/api/v1", func(r chi.Router) {
		// Publicas: registro y login.
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Protegidas: exigen JWT valido en Authorization: Bearer <token>.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			r.Get("/productos", servidor.ListarProductos)
			r.Post("/productos", servidor.CrearProducto)
			r.Get("/productos/{id}", servidor.ObtenerProducto)
			r.Put("/productos/{id}", servidor.ActualizarProducto)
			r.Delete("/productos/{id}", servidor.BorrarProducto)

			r.Get("/categorias", servidor.ListarCategorias)
			r.Post("/categorias", servidor.CrearCategoria)
			r.Get("/categorias/{id}", servidor.ObtenerCategoria)
			r.Put("/categorias/{id}", servidor.ActualizarCategoria)
			r.Delete("/categorias/{id}", servidor.BorrarCategoria)
		})
	})

	// 6. Servidor HTTP configurado por Options (puerto + timeouts desde config).
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	// 7. Contexto que se cancela al recibir Ctrl+C o SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 8. Arrancar el servidor en una goroutine para no bloquear la espera de la senal.
	errServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor escuchando en http://localhost%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	// 9. Esperar: o el servidor falla, o llega la senal de apagado.
	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Senal de apagado recibida, cerrando ordenadamente...")
	}

	// 10. Graceful shutdown: hasta 10s para terminar las requests en curso.
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}
	log.Println("Servidor detenido limpiamente.")
	return nil
}
