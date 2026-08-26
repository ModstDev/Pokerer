package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ModstDev/Pokerer/internal/httpapi/middleware"
	"github.com/ModstDev/Pokerer/internal/poker"
)

type createTableRequest struct {
	Name       string `json:"name"`
	SmallBlind int64  `json:"small_blind"`
	BigBlind   int64  `json:"big_blind"`
	MinBuyIn   int64  `json:"min_buy_in"`
	MaxBuyIn   int64  `json:"max_buy_in"`
	MaxPlayers int    `json:"max_players"`
}

type joinTableRequest struct {
	BuyIn int64 `json:"buy_in"`
}

func (s *Server) createTable(w http.ResponseWriter, r *http.Request) {
	var req createTableRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	table, err := s.pokerService.CreateTable(r.Context(), poker.CreateTableInput{
		Name:       req.Name,
		SmallBlind: req.SmallBlind,
		BigBlind:   req.BigBlind,
		MinBuyIn:   req.MinBuyIn,
		MaxBuyIn:   req.MaxBuyIn,
		MaxPlayers: req.MaxPlayers,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, table)
}

func (s *Server) listTables(w http.ResponseWriter, r *http.Request) {
	tables, err := s.pokerService.ListTables(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get tables")
		return
	}

	writeJSON(w, http.StatusOK, tables)
}

func (s *Server) joinTable(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	tableID := r.PathValue("id")

	var req joinTableRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := s.pokerService.JoinTable(r.Context(), tableID, userID, req.BuyIn); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to join table")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) leaveTable(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	tableID := r.PathValue("id")

	if err := s.pokerService.LeaveTable(r.Context(), tableID, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "couldn't leave table")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
