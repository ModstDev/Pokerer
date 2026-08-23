package httpapi

import (
	"net/http"

	"github.com/ModstDev/Pokerer/internal/auth"
	"github.com/ModstDev/Pokerer/internal/auth/token"
)

type Server struct {
	authService    *auth.Service
	tokenGenerator *token.JWT
}

func NewServer(authService *auth.Service, tokenGenerator *token.JWT) *Server {
	return &Server{
		authService:    authService,
		tokenGenerator: tokenGenerator,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", s.register)
	mux.HandleFunc("POST /api/v1/auth/login", s.login)

	return mux
}
