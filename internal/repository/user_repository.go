package repository

import (
	"context"
	"database/sql"

	database "github.com/ModstDev/Pokerer/internal/database/generated"
)

type UserRepository struct {
	queries *database.Queries
}

func NewUserRepository(queries *database.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) Create(ctx context.Context, params database.CreateUserParams) error {
	return r.queries.CreateUser(ctx, params)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (database.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *UserRepository) WithTx(tx *sql.Tx) *UserRepository {
	return &UserRepository{
		queries: r.queries.WithTx(tx),
	}
}
