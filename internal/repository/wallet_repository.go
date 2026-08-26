package repository

import (
	"context"
	"database/sql"

	database "github.com/ModstDev/Pokerer/internal/database/generated"
)

type WalletRepository struct {
	queries *database.Queries
}

func NewWalletRepository(queries *database.Queries) *WalletRepository {
	return &WalletRepository{
		queries: queries,
	}
}

func (r *WalletRepository) Create(ctx context.Context, params database.CreateWalletParams) error {
	return r.queries.CreateWallet(ctx, params)
}

func (r *WalletRepository) GetByUserID(ctx context.Context, userID string) (database.Wallet, error) {
	return r.queries.GetWalletByUserID(ctx, userID)
}

func (r *WalletRepository) CreateTransaction(ctx context.Context, params database.CreateWalletTransactionParams) error {
	return r.queries.CreateWalletTransaction(ctx, params)
}

func (r *WalletRepository) GetTransactionByUserID(ctx context.Context, userID string) ([]database.WalletTransaction, error) {
	return r.queries.GetWalletTransactionsByUserID(ctx, userID)
}

func (r *WalletRepository) WithTx(tx *sql.Tx) *WalletRepository {
	return &WalletRepository{
		queries: r.queries.WithTx(tx),
	}
}
