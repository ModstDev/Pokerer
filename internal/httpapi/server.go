package httpapi

import (
	"net/http"

	"github.com/ModstDev/Pokerer/internal/auth"
	"github.com/ModstDev/Pokerer/internal/auth/token"
	"github.com/ModstDev/Pokerer/internal/httpapi/middleware"
	"github.com/ModstDev/Pokerer/internal/repository"
)

type Server struct {
	authService    *auth.Service
	tokenGenerator *token.JWT
	userRepository *repository.UserRepository
}

func NewServer(authService *auth.Service, tokenGenerator *token.JWT, userRepository *repository.UserRepository) *Server {
	return &Server{
		authService:    authService,
		tokenGenerator: tokenGenerator,
		userRepository: userRepository,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	protected := middleware.Auth(s.tokenGenerator)

	mux.Handle("GET /api/v1/me", protected(http.HandlerFunc(s.me)))

	mux.HandleFunc("POST /api/v1/auth/register", s.register)
	mux.HandleFunc("POST /api/v1/auth/login", s.login)

	return mux
}
