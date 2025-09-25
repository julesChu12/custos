package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/your-org/custos/internal/domain/entity"
	"github.com/your-org/custos/internal/domain/repository"
	"github.com/your-org/custos/pkg/constants"
	"github.com/your-org/custos/pkg/errors"
)

type AuthService struct {
	userRepo     repository.UserRepository
	tokenService *TokenService
}

func NewAuthService(userRepo repository.UserRepository, tokenService *TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	if len(username) < constants.UsernameMinLength || len(username) > constants.UsernameMaxLength {
		return nil, errors.NewInvalidPasswordError(
			fmt.Sprintf("Username must be between %d and %d characters",
				constants.UsernameMinLength, constants.UsernameMaxLength))
	}

	if len(password) < constants.PasswordMinLength || len(password) > constants.PasswordMaxLength {
		return nil, errors.NewInvalidPasswordError(
			fmt.Sprintf("Password must be between %d and %d characters",
				constants.PasswordMinLength, constants.PasswordMaxLength))
	}

	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, errors.NewUserAlreadyExistsError(username)
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.NewUserAlreadyExistsError(email)
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entity.NewUser(username, email, hashedPassword)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, *entity.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, errors.NewInvalidCredentialsError()
	}

	if !user.IsActive() {
		return nil, nil, errors.NewInvalidCredentialsError()
	}

	if !s.checkPassword(password, user.Password) {
		return nil, nil, errors.NewInvalidCredentialsError()
	}

	tokenPair, err := s.tokenService.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenPair, user, nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}