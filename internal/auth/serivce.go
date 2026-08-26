package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ModstDev/Pokerer/internal/auth/token"
	"github.com/ModstDev/Pokerer/internal/database"
	generated "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

type Service struct {
	db      *sql.DB
	users   *repository.UserRepository
	wallets *repository.WalletRepository
}

func NewService(db *sql.DB, users *repository.UserRepository, wallets *repository.WalletRepository) *Service {
	return &Service{
		db:      db,
		users:   users,
		wallets: wallets,
	}
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	User        generated.User
	AccessToken string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (generated.User, error) {
	username := strings.TrimSpace(input.Username)
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if username == "" {
		return generated.User{}, fmt.Errorf("username is required")
	}

	if email == "" {
		return generated.User{}, fmt.Errorf("email is required")
	}

	if len(input.Password) < 8 {
		return generated.User{}, fmt.Errorf("password must contain at least 8 characters")
	}

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return generated.User{}, fmt.Errorf("email already exists")
	}

	if err != sql.ErrNoRows {
		return generated.User{}, fmt.Errorf("checking email: %w", err)
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return generated.User{}, fmt.Errorf("hashing password: %w", err)
	}

	userID := uuid.New().String()
	walletID := uuid.New().String()

	err = database.WithTransaction(ctx, s.db, func(tx *sql.Tx) error {
		users := s.users.WithTx(tx)
		wallets := s.wallets.WithTx(tx)

		if err := users.Create(ctx, generated.CreateUserParams{
			ID:           userID,
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}); err != nil {
			return fmt.Errorf("creating user: %w", err)
		}

		if err := wallets.Create(ctx, generated.CreateWalletParams{
			ID:     walletID,
			UserID: userID,
		}); err != nil {
			return fmt.Errorf("creating wallet: %w", err)
		}

		return nil
	})
	if err != nil {
		return generated.User{}, err
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return generated.User{}, fmt.Errorf("getting created user: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput, tokenGenerator *token.JWT) (LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if email == "" {
		return LoginResult{}, fmt.Errorf("email is required")
	}

	if input.Password == "" {
		return LoginResult{}, fmt.Errorf("password is required")
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return LoginResult{}, fmt.Errorf("invalid email or password")
		}

		return LoginResult{}, fmt.Errorf("getting user: %w", err)
	}

	if !CheckPassword(input.Password, user.PasswordHash) {
		return LoginResult{}, fmt.Errorf("invalid email or password")
	}

	accessToken, err := tokenGenerator.Generate(user.ID, 15*time.Minute)
	if err != nil {
		return LoginResult{}, fmt.Errorf("generating access token: %w", err)
	}

	return LoginResult{
		User:        user,
		AccessToken: accessToken,
	}, nil
}
