package auth

import (
	"crocus-calculator/internal/api/middleware"
	"crocus-calculator/internal/dto"
	ht "crocus-calculator/internal/lib/http"
	"net/http"
)

func (c AuthController) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("Handling user registration")

		req := r.Context().Value(middleware.DataKey).(dto.RegisterRequest)

		err := c.authService.Register(req)
		if err != nil {
			ht.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		ht.SendMessage(w, r, http.StatusOK, "User registered successfully")
		c.logger.Info("User registered successfully", "login", req.Login)
	}
}
