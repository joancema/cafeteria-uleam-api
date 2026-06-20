// Command cafeteria-api arranca el servidor HTTP de la Cafeteria Universitaria.
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite" // driver database/sql "sqlite" (pure-Go) para el backend sqlc
	"github.com/glebarez/sqlite"      // driver GORM (pure-Go)
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"cafeteria-uleam-api/internal/handlers"
	"cafeteria-uleam-api/internal/middleware"
	"cafeteria-uleam-api/internal/models"
	"cafeteria-uleam-api/internal/service"
	"cafeteria-uleam-api/internal/storage"
)

func main() {
	// 1. GORM es el DUENO DEL ESQUEMA: abre la DB, migra y siembra.
	//    Ahora tambien migra la tabla de usuarios.
	gdb, err := gorm.Open(sqlite.Open("cafeteria.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}
	if err := gdb.AutoMigrate(&models.Producto{}, &models.Categoria{}, &models.Usuario{}); err != nil {
		log.Fatal("fallo AutoMigrate: ", err)
	}
	almacenGorm := storage.NuevoAlmacenSQLite(gdb)
	almacenGorm.SembrarSiVacio()

	// 2. Elegir el backend de productos/categorias segun STORAGE (igual que antes).
	var almacen storage.Almacen
	switch os.Getenv("STORAGE") {
	case "sqlc":
		sdb, err := sql.Open("sqlite", "cafeteria.db")
		if err != nil {
			log.Fatal("no se pudo abrir sql.DB para sqlc: ", err)
		}
		almacen = storage.NuevoAlmacenSQLC(sdb)
		log.Println("Backend de productos/categorias: sqlc (database/sql)")
	default:
		almacen = almacenGorm
		log.Println("Backend de productos/categorias: GORM")
	}

	// 3. Los usuarios viven SIEMPRE en GORM (decision de la semana). Por eso NO
	//    cerramos gdb aunque el backend de productos sea sqlc: GORM mantiene su
	//    conexion para la tabla de usuarios.
	usuarioRepo := storage.NuevoUsuarioGORM(gdb)

	// 4. Capa de servicio. Cada servicio recibe SOLO la interfaz estrecha que
	//    necesita; almacen (Almacen) cumple ProductoRepository y CategoriaRepository
	//    por embedding, asi que es asignable a ambos parametros.
	productoSvc := service.NuevoProductoService(almacen)
	categoriaSvc := service.NuevoCategoriaService(almacen)
	authSvc := service.NuevoAuthService(usuarioRepo)

	// 5. Server con los servicios inyectados.
	servidor := handlers.NewServer(productoSvc, categoriaSvc, authSvc)

	// 6. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 7. Rutas versionadas /api/v1/.
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

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
