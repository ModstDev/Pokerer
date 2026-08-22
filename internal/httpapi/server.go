package httpapi

import (
	"net/http"

	"github.com/ModstDev/Pokerer/internal/auth"
)

type Server struct {
	authService *auth.Service
}

func NewServer(authService *auth.Service) *Server {
	return &Server{
		authService: authService,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", s.register)

	return mux
}
