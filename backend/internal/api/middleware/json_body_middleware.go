package middleware

import (
	"context"
	ht "crocus-calculator/internal/lib/http"
	"crocus-calculator/internal/lib/validator"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

var DataKey key = "data"

func JsonBodyMiddleware[T any](logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				logger.Warnw("Failed to decode JSON body", "error", err.Error())
				ht.SendMessage(w, r, http.StatusBadRequest, "invalid request body")
				return
			}

			if err := validator.ValidateStruct(body); err != nil {
				logger.Warnw("Validation failed", "error", err.Error())
				ht.SendMessage(w, r, http.StatusBadRequest, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), DataKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
