package auth

import (
	"crocus-calculator/internal/api/middleware"
	"crocus-calculator/internal/dto"
	tokenservice "crocus-calculator/internal/lib/tokenService"
	"crocus-calculator/internal/shared"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(userDTO dto.RegisterRequest) *shared.HttpError
	Login(LoginRequest dto.LoginRequest) (*tokenservice.TokenPair, *shared.HttpError)
	Refresh(userID uint) (*tokenservice.TokenPair, *shared.HttpError)
}

type AuthController struct {
	authService  AuthService
	logger       *zap.SugaredLogger
	tokenService TokenService
}

func NewAuthController(authService AuthService, logger *zap.SugaredLogger, tokenService TokenService) *AuthController {
	return &AuthController{authService: authService, logger: logger, tokenService: tokenService}
}

func (c AuthController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.With(middleware.JsonBodyMiddleware[dto.RegisterRequest](c.logger)).Post("/register", c.RegisterHandler())
		r.With(middleware.JsonBodyMiddleware[dto.LoginRequest](c.logger)).Post("/login", c.LoginHandler())
		r.With(middleware.ValidateToken(c.tokenService, c.logger)).Post("/refresh", c.RefreshHandler())
	})
}
