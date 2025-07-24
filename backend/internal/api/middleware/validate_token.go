package middleware

import (
	"context"
	ht "crocus-calculator/internal/lib/http"
	tokenservice "crocus-calculator/internal/lib/tokenService"
	"net/http"

	"go.uber.org/zap"
)

var UserIDKey key = "userIDKey"

type TokenService interface {
	ValidateToken(tokenString string) (*tokenservice.Claims, error)
}

func ValidateToken(tokenservice TokenService, logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				logger.Warn("No authorization token provided")
				ht.SendMessage(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims, err := tokenservice.ValidateToken(tokenString)
			if err != nil {
				logger.Warn("Invalid authorization token", zap.Error(err))
				ht.SendMessage(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Attach the claims to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
