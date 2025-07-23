package auth

import (
	"crocus-calculator/internal/dto"
	tokenservice "crocus-calculator/internal/lib/tokenService"
	"crocus-calculator/internal/shared"
	"crocus-calculator/internal/storage"
	"crocus-calculator/internal/storage/repository"
	"errors"
	"time"

	"go.uber.org/zap"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPwd string, plainPwd string) error
}

type AuthRepository interface {
	Create(user *storage.User) error
	Login(credentials *dto.LoginRequest) (int64, string, error)
}

type TokenService interface {
	GenerateTokenPair(userID uint) (*tokenservice.TokenPair, error)
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) (*tokenservice.Claims, error)
}

type authService struct {
	repo         AuthRepository
	hasher       PasswordHasher
	logger       *zap.SugaredLogger
	tokenService TokenService
}

func NewAuthService(repo AuthRepository, hasher PasswordHasher, logger *zap.SugaredLogger, tokenService TokenService) AuthService {
	return &authService{
		repo:         repo,
		hasher:       hasher,
		logger:       logger,
		tokenService: tokenService,
	}
}

func (s authService) Register(registerDTO dto.RegisterRequest) *shared.HttpError {
	var user = &storage.User{
		Login:     registerDTO.Login,
		Password:  registerDTO.Password,
		Permit:    registerDTO.Permit,
		FullName:  registerDTO.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	hashedPassword, err := s.hasher.HashPassword(registerDTO.Password)
	if err != nil {
		s.logger.Errorf("failed to hash password: %v", err)
		return shared.InternalError
	}
	user.Password = hashedPassword

	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, repository.ErrLoginAlreadyExists) {
			s.logger.Warnf("user with login %s already exists", user.Login)
			return shared.InvalidCredentials
		}

		s.logger.Errorf("failed to create user: %v", err)
		return shared.InternalError
	}

	s.logger.Infof("user registered successfully: %v", user)
	return nil
}

func (s authService) Login(loginDTO dto.LoginRequest) (*tokenservice.TokenPair, *shared.HttpError) {
	s.logger.Infof("attempting login for user: %s", loginDTO.Login)

	userID, hashedPassword, err := s.repo.Login(&loginDTO)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			s.logger.Warnf("invalid credentials for user: %s", loginDTO.Login)
			return nil, shared.InvalidCredentials
		}

		s.logger.Errorf("failed to login user %s: %v", loginDTO.Login, err)
		return nil, shared.InternalError
	}

	if err := s.hasher.ComparePasswords(hashedPassword, loginDTO.Password); err != nil {
		s.logger.Warnf("password mismatch for user: %s", loginDTO.Login)
		return nil, shared.InvalidCredentials
	}

	tokens, err := s.tokenService.GenerateTokenPair(uint(userID))
	if err != nil {
		s.logger.Errorf("failed to generate tokens for user %d: %v", userID, err)
		return nil, shared.InternalError
	}

	s.logger.Infof("user logged in successfully: %d", userID)
	return tokens, nil
}
