package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/your-org/custos/pkg/constants"
	"github.com/your-org/custos/pkg/errors"
	"github.com/your-org/custos/pkg/types"
)

type TokenService struct {
	secretKey string
	issuer    string
	ttl       time.Duration
}

type TokenClaims struct {
	UserID   uint           `json:"user_id"`
	Username string         `json:"username"`
	Role     types.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewTokenService(secretKey string, ttl time.Duration) *TokenService {
	return &TokenService{
		secretKey: secretKey,
		issuer:    constants.JWTIssuer,
		ttl:       ttl,
	}
}

func (s *TokenService) GenerateAccessToken(userID uint, username string, role types.UserRole) (*TokenPair, error) {
	now := time.Now()
	claims := &TokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &TokenPair{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.ttl.Seconds()),
	}, nil
}

func (s *TokenService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		if jwt.IsValidationError(err, jwt.ValidationErrorExpired) {
			return nil, errors.NewTokenExpiredError()
		}
		return nil, errors.NewTokenInvalidError()
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.NewTokenInvalidError()
	}

	return claims, nil
}