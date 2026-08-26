package poker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/ModstDev/Pokerer/internal/database"
	generated "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

type Service struct {
	db      *sql.DB
	tables  *repository.PokerTableRepository
	wallets *repository.WalletRepository
}

func NewService(db *sql.DB, tables *repository.PokerTableRepository, wallets *repository.WalletRepository) *Service {
	return &Service{
		db:      db,
		tables:  tables,
		wallets: wallets,
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

func (s *Service) JoinTable(ctx context.Context, tableID string, userID string, buyIn int64) error {
	if buyIn <= 0 {
		return errors.New("buy-ing must be positive")
	}

	return database.WithTransaction(ctx, s.db, func(tx *sql.Tx) error {
		tables := s.tables.WithTx(tx)
		wallets := s.wallets.WithTx(tx)

		table, err := tables.GetByID(ctx, tableID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("table not found")
			}

			return fmt.Errorf("getting table: %w", err)
		}

		if table.Status != "waiting" {
			return errors.New("table is not accepting plaayers")
		}

		if buyIn < table.MinBuyIn || buyIn > table.MaxBuyIn {
			return errors.New("buy-in is outside the table limits")
		}

		players, err := tables.ListPlayers(ctx, tableID)
		if err != nil {
			return fmt.Errorf("getting table players: %w", err)
		}

		if len(players) >= int(table.MaxPlayers) {
			return errors.New("table is full")
		}

		_, err = tables.GetPlayer(ctx, tableID, userID)

		if err == nil {
			return errors.New("user is already at this table")
		}

		if err != sql.ErrNoRows {
			return fmt.Errorf("checking table membership: %w", err)
		}

		wallet, err := wallets.GetByUserID(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("wallet not found")
			}

			return fmt.Errorf("getting wallet: %w", err)
		}

		if wallet.Balance < buyIn {
			return errors.New("insufficient balance")
		}

		seat, err := findFreeSeat(players, int(table.MaxPlayers))
		if err != nil {
			return errors.New("no available seat")
		}

		newBalance := wallet.Balance - buyIn

		if err := wallets.UpdateBalance(ctx, wallet.ID, newBalance); err != nil {
			return fmt.Errorf("updating wallet: %w", err)
		}

		if err := wallets.CreateTransaction(ctx, generated.CreateWalletTransactionParams{
			ID:           uuid.New().String(),
			WalletID:     wallet.ID,
			Type:         "table_buy_in",
			Amount:       -buyIn,
			BalanceAfter: newBalance,
		},
		); err != nil {
			return fmt.Errorf("creating wallet transaction: %w", err)
		}

		if err := tables.AddPlayer(ctx, generated.AddTablePlayerParams{
			ID:         uuid.New().String(),
			TableID:    tableID,
			UserID:     userID,
			SeatNumber: int32(seat),
			Chips:      buyIn,
		},
		); err != nil {
			return fmt.Errorf("adding player to table: %w", err)
		}

		return nil
	})
}

func findFreeSeat(players []generated.TablePlayer, maxPlayers int) (int, error) {
	occupied := make(map[int]struct{}, len(players))

	for _, player := range players {
		occupied[int(player.SeatNumber)] = struct{}{}
	}

	for seat := 0; seat < maxPlayers; seat++ {
		if _, exists := occupied[seat]; !exists {
			return seat, nil
		}
	}

	return 0, fmt.Errorf("no free seat")
}

// TODO
func (s *Service) LeaveTable(ctx context.Context, tableID string, userID string) error {
	return nil
}
