package auth

import (
	"fmt"
	"os"
	"time"

	"clean-architecture-api/internal/domain/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthService interface {
	GenerateTokenPair(userID uuid.UUID, email, role string) (*TokenPair, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshTokenPair(refreshToken string) (*TokenPair, error)
}

type authService struct {
	secretKey []byte
}

func NewAuthService() (AuthService, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is required")
	}
	return &authService{
		secretKey: []byte(secretKey),
	}, nil
}

func (s *authService) GenerateTokenPair(userID uuid.UUID, email, role string) (*TokenPair, error) {
	accessTokenExp := time.Now().Add(15 * time.Minute)
	accessTokenClaims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "clean-architecture-api",
			Subject:   userID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, errors.ErrFailedToGenerateAccessToken
	}

	refreshTokenExp := time.Now().Add(7 * 24 * time.Hour)
	refreshTokenClaims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "clean-architecture-api",
			Subject:   userID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, errors.ErrFailedToGenerateRefreshToken
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessTokenExp.Unix(),
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrUnexpectedSigningMethod
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, errors.ErrFailedToParseToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidToken
}

func (s *authService) RefreshTokenPair(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
}
