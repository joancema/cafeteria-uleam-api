// Package httpserver construye un *http.Server configurable mediante el patron
// funcional Options. Es el SEGUNDO ejemplo del patron en el proyecto (el primero
// es service.AuthService): aqui los valores opcionales con default son la
// direccion de escucha y los timeouts.
package httpserver

import (
	"net/http"
	"time"
)

// Valores por defecto del servidor si no se pasa ninguna Option.
const (
	puertoPorDefecto       = ":8080"
	readTimeoutPorDefecto  = 10 * time.Second
	writeTimeoutPorDefecto = 10 * time.Second
	idleTimeoutPorDefecto  = 60 * time.Second
)

// Opcion configura el *http.Server en su construccion.
type Opcion func(*http.Server)

// ConPuerto fija la direccion de escucha, ej ":8080".
func ConPuerto(puerto string) Opcion {
	return func(s *http.Server) {
		if puerto != "" {
			s.Addr = puerto
		}
	}
}

// ConReadTimeout fija el timeout de lectura.
func ConReadTimeout(d time.Duration) Opcion {
	return func(s *http.Server) {
		if d > 0 {
			s.ReadTimeout = d
		}
	}
}

// ConWriteTimeout fija el timeout de escritura.
func ConWriteTimeout(d time.Duration) Opcion {
	return func(s *http.Server) {
		if d > 0 {
			s.WriteTimeout = d
		}
	}
}

// Nuevo construye un *http.Server con el handler dado y aplica las Options.
//
// Los timeouts no son un lujo: un servidor sin ellos deja conexiones lentas o
// maliciosas ocupando recursos indefinidamente. Por eso aqui SIEMPRE hay un
// default, y las Options solo permiten ajustarlo.
func Nuevo(handler http.Handler, opts ...Opcion) *http.Server {
	s := &http.Server{
		Addr:         puertoPorDefecto,
		Handler:      handler,
		ReadTimeout:  readTimeoutPorDefecto,
		WriteTimeout: writeTimeoutPorDefecto,
		IdleTimeout:  idleTimeoutPorDefecto,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
