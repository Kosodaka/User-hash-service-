package httpserver

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 30 * time.Second
	defaultWriteTimeout    = 30 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, port string) *Server {
	httpServer := &http.Server{
		Handler:        handler,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
		Addr:           net.JoinHostPort("", port),
		MaxHeaderBytes: 100 << 20,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	return s
}

func (s *Server) Serve() {
	log.Println("Starting the server...")

	err := s.server.ListenAndServe()
	if err != nil {
		log.Printf("Server stopped with error: %v", err)
	}
	s.notify <- err
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
