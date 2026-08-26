package repository

import (
	"context"

	generated "github.com/ModstDev/Pokerer/internal/database/generated"
)

type PokerTableRepository struct {
	queries *generated.Queries
}

func NewPokerTableRepository(queries *generated.Queries) *PokerTableRepository {
	return &PokerTableRepository{
		queries: queries,
	}
}

func (r *PokerTableRepository) Create(ctx context.Context, params generated.CreatePokerTableParams) error {
	return r.queries.CreatePokerTable(ctx, params)
}

func (r *PokerTableRepository) GetByID(ctx context.Context, id string) (generated.PokerTable, error) {
	return r.queries.GetPokerTableByID(ctx, id)
}

func (r *PokerTableRepository) List(ctx context.Context) ([]generated.ListPokerTablesRow, error) {
	return r.queries.ListPokerTables(ctx)
}

func (r *PokerTableRepository) AddPlayerr(ctx context.Context, params generated.AddTablePlayerParams) error {
	return r.queries.AddTablePlayer(ctx, params)
}

func (r *PokerTableRepository) RemovePlayer(ctx context.Context, tableID, userID string) error {
	return r.queries.RemoveTablePlayer(ctx, generated.RemoveTablePlayerParams{
		TableID: tableID,
		UserID:  userID,
	})
}

func (r *PokerTableRepository) GetPlayer(ctx context.Context, tableID string, userID string) (generated.TablePlayer, error) {
	return r.queries.GetTablePlayer(ctx, generated.GetTablePlayerParams{
		TableID: tableID,
		UserID:  userID,
	})
}

func (r *PokerTableRepository) ListPlayers(ctx context.Context, tableID string) ([]generated.TablePlayer, error) {
	return r.queries.ListTablePlayers(ctx, tableID)
}
