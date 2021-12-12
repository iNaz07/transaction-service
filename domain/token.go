package domain

import (
	"time"
)

type User struct {
	ID   int64
	IIN  string `json:"iin"`
	Role string
}

type JwtToken struct {
	AccessSecret string
	AccessTtl    time.Duration
}

type JwtTokenUsecase interface {
	ParseTokenAndGetID(token string) (int64, error)
	ParseTokenAndGetRole(token string) (string, error)
	ParseTokenAndGetIIN(token string) (string, error)
	GetAccessTTL() time.Duration
}
