package domain

import (
	"time"
)

type JwtToken struct {
	AccessSecret string
	AccessTtl    time.Duration
}

type JwtTokenUsecase interface {
	ParseTokenAndGetID(token string) (int64, error)
	ParseTokenAndGetRole(token string) (string, error)
	// JWTErrorChecker(err error, c echo.Context) error
	GetAccessTTL() time.Duration
}
