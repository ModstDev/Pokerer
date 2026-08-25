package httpapi

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ModstDev/Pokerer/internal/auth"
	"github.com/ModstDev/Pokerer/internal/httpapi/middleware"
)

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
	}

	user, err := s.authService.Register(r.Context(), auth.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
	}

	response := userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	writeJSON(w, http.StatusCreated, response)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := s.authService.Login(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	},
		s.tokenGenerator,
	)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := loginResponse{
		AccessToken: result.AccessToken,
		User: userResponse{
			ID:        result.User.ID,
			Username:  result.User.Username,
			Email:     result.User.Email,
			CreatedAt: result.User.CreatedAt.Format(time.RFC3339),
		},
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	user, err := s.userRepository.GetByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusOK, userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
