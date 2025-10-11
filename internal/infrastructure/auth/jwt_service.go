package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var OneMonth = 30 * 24 * time.Hour

type JWTService interface {
	GenerateToken(userId uuid.UUID, tokenLifetime time.Duration) (string, error)
	ValidateToken(tokenStr string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{secretKey}
}

type JWTClaims struct {
	UserId uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateToken(userId uuid.UUID, tokenLifetime time.Duration) (string, error) {
	claims := &JWTClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
