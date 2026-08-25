package wallet

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ModstDev/Pokerer/internal/database"
	generated "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	db      *sql.DB
	wallets *repository.WalletRepository
}

func NewService(db *sql.DB, wallets *repository.WalletRepository) *Service {
	return &Service{
		db:      db,
		wallets: wallets,
	}
}

func (s *Service) Deposit(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}

	return database.WithTransaction(ctx, s.db, func(tx *sql.Tx) error {
		queries := generated.New(tx)

		wallet, err := queries.GetWalletByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("getting wallet: %w", err)
		}

		newBalance := wallet.Balance + amount

		if err := queries.UpdateWalletBalance(ctx, generated.UpdateWalletBalanceParams{
			Balance: newBalance,
			ID:      wallet.ID,
		},
		); err != nil {
			return fmt.Errorf("updating wallet balance: %w", err)
		}

		if err := queries.CreateWalletTransaction(ctx, generated.CreateWalletTransactionParams{
			ID:           uuid.New().String(),
			WalletID:     wallet.ID,
			Type:         "deposit",
			Amount:       amount,
			BalanceAfter: newBalance,
		},
		); err != nil {
			return fmt.Errorf("creating wallet transaction: %w", err)
		}

		return nil

	})
}
