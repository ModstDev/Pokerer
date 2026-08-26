package httpapi

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ModstDev/Pokerer/internal/httpapi/middleware"
)

type walletResponse struct {
	ID        string    `json:"id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type walletTransactionResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Amount       int64     `json:"amount"`
	BalanceAfter int64     `json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}

type depositRequest struct {
	Amount int64 `json:"amount"`
}

func (s *Server) getWallet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	wallet, err := s.walletService.GetByUserID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "wallet not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, walletResponse{
		ID:        wallet.ID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	})
}

func (s *Server) deposit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req depositRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	if err := s.walletService.Deposit(r.Context(), userID, req.Amount); err != nil {
		writeError(w, http.StatusInternalServerError, "deposit failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getWalletTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	transactions, err := s.walletService.GetTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "transactions not found")
		return
	}

	response := make([]walletTransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		response = append(response, walletTransactionResponse{
			ID:           transaction.ID,
			Type:         transaction.Type,
			Amount:       transaction.Amount,
			BalanceAfter: transaction.BalanceAfter,
			CreatedAt:    transaction.CreatedAt,
		})
	}

	writeJSON(w, http.StatusOK, response)
}
