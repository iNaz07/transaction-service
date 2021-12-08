package usecase

import (
	"fmt"
	"time"
	"transaction-service/domain"

	"github.com/dgrijalva/jwt-go"
)

type jwtUsecase struct {
	token *domain.JwtToken
}

func NewJWTUseCase(token *domain.JwtToken) domain.JwtTokenUsecase {
	return &jwtUsecase{token: token}
}

func (j *jwtUsecase) ParseTokenAndGetID(token string) (int64, error) {
	claims, err := j.ParseToken(token)
	if err != nil {
		return -1, fmt.Errorf("invalid token: %w", err)
	}
	id, ok := claims["id"].(float64)
	if !ok {
		return -1, fmt.Errorf("id not found from token")
	}
	return int64(id), nil
}

func (j *jwtUsecase) ParseTokenAndGetRole(token string) (string, error) {
	claims, err := j.ParseToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found from token")
	}
	return role, nil
}

func (j *jwtUsecase) GetAccessTTL() time.Duration {
	return j.token.AccessTtl
}

func (j *jwtUsecase) ParseToken(token string) (jwt.MapClaims, error) {
	JWTToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to extract token metadata, unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.token.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := JWTToken.Claims.(jwt.MapClaims)
	if ok && JWTToken.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("field exp not found from token")
		}
		expiredTime := time.Unix(int64(exp), 0)
		if time.Now().After(expiredTime) {
			return nil, fmt.Errorf("token expired")
		}
	}
	return claims, nil
}
