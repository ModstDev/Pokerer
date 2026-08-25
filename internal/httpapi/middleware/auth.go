package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ModstDev/Pokerer/internal/auth/token"
)

type contextKey string

const userIDKey contextKey = "user_id"

func Auth(jwt *token.JWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				writeError(w, http.StatusUnauthorized, "authorization required")
				return
			}

			parts := strings.SplitN(header, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			userID, err := jwt.Validate(parts[1])
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}
