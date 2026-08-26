package poker

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	generated "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

type Service struct {
	tables *repository.PokerTableRepository
}

func NewService(tables *repository.PokerTableRepository) *Service {
	return &Service{
		tables: tables,
	}
}

type CreateTableInput struct {
	Name       string
	SmallBlind int64
	BigBlind   int64
	MinBuyIn   int64
	MaxBuyIn   int64
	MaxPlayers int
}

func (s *Service) CreateTable(ctx context.Context, input CreateTableInput) (generated.PokerTable, error) {
	if input.Name == "" {
		return generated.PokerTable{}, fmt.Errorf("table name is required")
	}
	if input.SmallBlind <= 0 {
		return generated.PokerTable{}, fmt.Errorf("small blind must be positive")
	}

	if input.BigBlind <= input.SmallBlind {
		return generated.PokerTable{}, fmt.Errorf("big blind must be greater than small blind")
	}

	if input.MinBuyIn <= 0 {
		return generated.PokerTable{}, fmt.Errorf("minimum buy-in must be positive")
	}

	if input.MaxBuyIn < input.MinBuyIn {
		return generated.PokerTable{}, fmt.Errorf("maximum buy-in must be greater than or equal to minimum buy-in")
	}

	if input.MaxPlayers < 2 || input.MaxPlayers > 9 {
		return generated.PokerTable{}, fmt.Errorf("maximum players must be between 2 and 9")
	}

	id := uuid.New().String()

	err := s.tables.Create(ctx, generated.CreatePokerTableParams{
		ID:         id,
		Name:       input.Name,
		SmallBlind: input.SmallBlind,
		BigBlind:   input.BigBlind,
		MinBuyIn:   input.MinBuyIn,
		MaxBuyIn:   input.MaxBuyIn,
		MaxPlayers: int32(input.MaxPlayers),
		Status:     "waiting",
	})
	if err != nil {
		return generated.PokerTable{}, fmt.Errorf("creating poker table: %w", err)
	}

	return s.tables.GetByID(ctx, id)
}

func (s *Service) ListTables(ctx context.Context) ([]generated.ListPokerTablesRow, error) {
	return s.tables.List(ctx)
}
