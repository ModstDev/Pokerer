package httpapi

import (
	"net/http"

	"github.com/ModstDev/Pokerer/internal/auth"
	"github.com/ModstDev/Pokerer/internal/auth/token"
	"github.com/ModstDev/Pokerer/internal/httpapi/middleware"
	"github.com/ModstDev/Pokerer/internal/poker"
	"github.com/ModstDev/Pokerer/internal/repository"
	"github.com/ModstDev/Pokerer/internal/wallet"
)

type Server struct {
	authService    *auth.Service
	tokenGenerator *token.JWT
	userRepository *repository.UserRepository
	walletService  *wallet.Service
	pokerService   *poker.Service
}

func NewServer(authService *auth.Service,
	tokenGenerator *token.JWT,
	userRepository *repository.UserRepository,
	walletService *wallet.Service,
	pokerService *poker.Service,
) *Server {
	return &Server{
		authService:    authService,
		tokenGenerator: tokenGenerator,
		userRepository: userRepository,
		walletService:  walletService,
		pokerService:   pokerService,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	protected := middleware.Auth(s.tokenGenerator)

	mux.Handle("GET /api/v1/me", protected(http.HandlerFunc(s.me)))
	mux.Handle("GET /api/v1/wallet", protected(http.HandlerFunc(s.getWallet)))
	mux.Handle("GET /api/v1/wallet/transactions", protected(http.HandlerFunc(s.getWalletTransactions)))
	mux.HandleFunc("GET /api/v1/tables", s.listTables)

	mux.HandleFunc("POST /api/v1/auth/register", s.register)
	mux.HandleFunc("POST /api/v1/auth/login", s.login)
	mux.Handle("POST /api/v1/wallet/deposit", protected(http.HandlerFunc(s.deposit)))
	mux.Handle("POST /api/v1/tables", protected(http.HandlerFunc(s.createTable)))
	mux.Handle("POST /api/v1/tables/{id}/join", protected(http.HandlerFunc(s.joinTable)))

	return mux
}
