package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tiagobnarita/go_learn/internal/repository"
	"github.com/tiagobnarita/go_learn/internal/service"
)

const defaultPort = 8080

type Server struct {
	port            int
	bookmarkService service.BookmarkService
}

// NewServer creates a new Server and returns it alongside the configured *http.Server.
func NewServer(repo repository.Repository) (*Server, *http.Server) {
	s := &Server{
		port:            defaultPort,
		bookmarkService: service.NewBookmarkService(repo),
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	return s, httpServer
}
