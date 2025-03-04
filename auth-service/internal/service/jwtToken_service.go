package service

import (
	"auth-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type jwtTokenService struct {
	secretKey     string
	redisClient   *redis.Client
	tokenDuration time.Duration
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJwtTokenService(secretKey string, redisClient *redis.Client) domain.TokenService {
	return &jwtTokenService{
		secretKey:     secretKey,
		redisClient:   redisClient,
		tokenDuration: 24 * time.Hour, // token 24 jam
	}
}

func (s *jwtTokenService) GenerateToken(userID uint64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *jwtTokenService) ValidateToken(tokenString string) (uint64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil

}

func (s *jwtTokenService) BlacklistToken(token string) error {
	// if Redis is not available, assume success
	if s.redisClient == nil {
		log.Println("Warning: Redis not available, token blacklisting skipped")
		return nil
	}

	ctx := context.Background()
	// store token in redis with expiration time
	err := s.redisClient.Set(ctx,
		"blacklist:"+token,
		true,
		s.tokenDuration,
	).Err()

	if err != nil {
		log.Printf("Error blacklisting token: %v", err)
		return fmt.Errorf("failed to blacklist token: %v", err)
	}

	return nil
}

func (s *jwtTokenService) IsTokenBlacklisted(token string) bool {
	if s.redisClient == nil {
		// if Redis is not available, assume token is valid
		return false
	}

	ctx := context.Background()

	exists, _ := s.redisClient.Exists(ctx, "blacklist:"+token).Result()
	return exists > 0
}
