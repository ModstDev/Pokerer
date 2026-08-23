package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
	issuer string
}

func NewJWT(secret, issuer string) *JWT {
	return &JWT{
		secret: []byte(secret),
		issuer: issuer,
	}
}

func (j *JWT) Generate(userID string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signedToken, nil
}

func (j *JWT) Validate(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return j.secret, nil
		})
	if err != nil {
		return "", fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims.Subject == "" {
		return "", fmt.Errorf("token subject is missing")
	}

	return claims.Subject, nil
}
