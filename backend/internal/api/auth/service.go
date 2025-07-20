package auth

import (
	"crocus-calculator/internal/dto"
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
}

type authService struct {
	repo   AuthRepository
	hasher PasswordHasher
	logger *zap.SugaredLogger
}

func NewAuthService(repo AuthRepository, hasher PasswordHasher, logger *zap.SugaredLogger) AuthService {
	return &authService{
		repo:   repo,
		hasher: hasher,
		logger: logger,
	}
}

func (s authService) Register(userDTO dto.RegisterRequest) *shared.HttpError {
	var user = &storage.User{
		Login:     userDTO.Login,
		Password:  userDTO.Password,
		Permit:    userDTO.Permit,
		FullName:  userDTO.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	hashedPassword, err := s.hasher.HashPassword(userDTO.Password)
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
