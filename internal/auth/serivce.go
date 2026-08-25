package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ModstDev/Pokerer/internal/auth/token"
	database "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

type Service struct {
	users *repository.UserRepository
}

func NewService(users *repository.UserRepository) *Service {
	return &Service{
		users: users,
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
	User        database.User
	AccessToken string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (database.User, error) {
	username := strings.TrimSpace(input.Username)
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if username == "" {
		return database.User{}, fmt.Errorf("username is required")
	}

	if email == "" {
		return database.User{}, fmt.Errorf("email is required")
	}

	if len(input.Password) < 8 {
		return database.User{}, fmt.Errorf("password must contain at least 8 characters")
	}

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return database.User{}, fmt.Errorf("email already exists")
	}

	if err != sql.ErrNoRows {
		return database.User{}, fmt.Errorf("checking email: %w", err)
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return database.User{}, fmt.Errorf("hashing password: %w", err)
	}

	userID := uuid.New().String()

	err = s.users.Create(ctx, database.CreateUserParams{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return database.User{}, fmt.Errorf("creating user: %w", err)
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return database.User{}, fmt.Errorf("getting created user: %w", err)
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
