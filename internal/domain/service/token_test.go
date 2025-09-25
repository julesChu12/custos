package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/custos/pkg/types"
)

func TestTokenService_GenerateAccessToken(t *testing.T) {
	tokenService := NewTokenService("test-secret-key", 15*60) // 15 minutes

	tests := []struct {
		name     string
		userID   uint
		username string
		role     types.UserRole
		wantErr  bool
	}{
		{
			name:     "valid user token",
			userID:   1,
			username: "testuser",
			role:     types.UserRoleUser,
			wantErr:  false,
		},
		{
			name:     "valid admin token",
			userID:   2,
			username: "admin",
			role:     types.UserRoleAdmin,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPair, err := tokenService.GenerateAccessToken(tt.userID, tt.username, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tokenPair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.Equal(t, "Bearer", tokenPair.TokenType)
				assert.Greater(t, tokenPair.ExpiresIn, int64(0))
			}
		})
	}
}

func TestTokenService_ValidateToken(t *testing.T) {
	tokenService := NewTokenService("test-secret-key", 15*60)

	// Generate a token first
	tokenPair, err := tokenService.GenerateAccessToken(1, "testuser", types.UserRoleUser)
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   tokenPair.AccessToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := tokenService.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, uint(1), claims.UserID)
				assert.Equal(t, "testuser", claims.Username)
				assert.Equal(t, types.UserRoleUser, claims.Role)
			}
		})
	}
}