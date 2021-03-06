package usecase

import (
	"fmt"
	"net/http"
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
		return -1, &domain.LogError{"invalid token", err, http.StatusBadRequest}
	}
	id, ok := claims["id"].(float64)
	if !ok {
		return -1, &domain.LogError{"invalid token", fmt.Errorf("id not found from token"), http.StatusBadRequest}
	}
	return int64(id), nil
}

func (j *jwtUsecase) ParseTokenAndGetRole(token string) (string, error) {
	claims, err := j.ParseToken(token)
	if err != nil {
		return "", &domain.LogError{"invalid token", err, http.StatusBadRequest}
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", &domain.LogError{"invalid token", fmt.Errorf("role not found from token"), http.StatusBadRequest}
	}
	return role, nil
}

func (j *jwtUsecase) ParseTokenAndGetIIN(token string) (string, error) {
	claims, err := j.ParseToken(token)
	if err != nil {
		return "", &domain.LogError{"invalid token", err, http.StatusBadRequest}
	}
	iin, ok := claims["iin"].(string)
	if !ok {
		return "", &domain.LogError{"invalid token", fmt.Errorf("iin not found from token"), http.StatusBadRequest}
	}
	return iin, nil
}

func (j *jwtUsecase) GetAccessTTL() time.Duration {
	return j.token.AccessTtl
}

func (j *jwtUsecase) ParseToken(token string) (jwt.MapClaims, error) {
	JWTToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &domain.LogError{"invalid token", fmt.Errorf("failed to extract token metadata, unexpected signing method: %v", token.Header["alg"]), http.StatusBadRequest}
		}
		return []byte(j.token.AccessSecret), nil
	})

	if err != nil {
		return nil, &domain.LogError{"invalid token", err, http.StatusBadRequest}
	}

	claims, ok := JWTToken.Claims.(jwt.MapClaims)
	if ok && JWTToken.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, &domain.LogError{"invalid token", fmt.Errorf("field exp not found from token"), http.StatusBadRequest}
		}
		expiredTime := time.Unix(int64(exp), 0)
		if time.Now().After(expiredTime) {
			return nil, &domain.LogError{"token expired", fmt.Errorf("token expired"), http.StatusBadRequest}
		}
	}
	return claims, nil
}
