package auth

import (
	"context"

	"github.com/your-org/custos/internal/application/dto"
	"github.com/your-org/custos/internal/domain/entity"
	"github.com/your-org/custos/internal/domain/service"
	"github.com/your-org/custos/pkg/types"
)

type RegisterUseCase struct {
	authService *service.AuthService
}

func NewRegisterUseCase(authService *service.AuthService) *RegisterUseCase {
	return &RegisterUseCase{
		authService: authService,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req *dto.RegisterRequest) (*dto.UserInfo, error) {
	user, err := uc.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     string(user.Role),
		Status:   string(user.Status),
	}, nil
}

type LoginUseCase struct {
	authService *service.AuthService
}

func NewLoginUseCase(authService *service.AuthService) *LoginUseCase {
	return &LoginUseCase{
		authService: authService,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	tokenPair, user, err := uc.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: tokenPair.AccessToken,
		TokenType:   tokenPair.TokenType,
		ExpiresIn:   tokenPair.ExpiresIn,
		User: &dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     string(user.Role),
			Status:   string(user.Status),
		},
	}, nil
}

func entityToUserInfo(user *entity.User) *dto.UserInfo {
	return &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     string(user.Role),
		Status:   string(user.Status),
	}
}