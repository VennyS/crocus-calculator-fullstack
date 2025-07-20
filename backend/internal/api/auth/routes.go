package auth

import (
	"crocus-calculator/internal/api/middleware"
	"crocus-calculator/internal/dto"
	"crocus-calculator/internal/shared"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(userDTO dto.RegisterRequest) *shared.HttpError
}

type AuthController struct {
	authService AuthService
	logger      *zap.SugaredLogger
}

func NewAuthController(authService AuthService, logger *zap.SugaredLogger) *AuthController {
	return &AuthController{authService: authService, logger: logger}
}

func (c AuthController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.With(middleware.JsonBodyMiddleware[dto.RegisterRequest](c.logger)).Post("/register", c.RegisterHandler())
	})
}
