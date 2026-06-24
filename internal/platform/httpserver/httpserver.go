// Package httpserver construye un *http.Server configurable por Options.
package httpserver

import (
	"net/http"
	"time"
)

const (
	puertoPorDefecto       = ":8080"
	readTimeoutPorDefecto  = 10 * time.Second
	writeTimeoutPorDefecto = 10 * time.Second
	idleTimeoutPorDefecto  = 60 * time.Second
)

type Opcion func(*http.Server)

func ConPuerto(p string) Opcion {
	return func(s *http.Server) {
		if p != "" {
			s.Addr = p
		}
	}
}

func ConReadTimeout(d time.Duration) Opcion {
	return func(s *http.Server) {
		if d > 0 {
			s.ReadTimeout = d
		}
	}
}

func ConWriteTimeout(d time.Duration) Opcion {
	return func(s *http.Server) {
		if d > 0 {
			s.WriteTimeout = d
		}
	}
}

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
