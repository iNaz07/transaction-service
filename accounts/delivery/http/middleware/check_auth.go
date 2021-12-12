package middleware

import (
	"fmt"
	"net/http"
	"transaction-service/domain"

	"github.com/labstack/echo/v4"
)

type Auth struct {
	JwtUsecase domain.JwtTokenUsecase
}

func InitAuth(token domain.JwtTokenUsecase) *Auth {
	return &Auth{JwtUsecase: token}
}

func (a *Auth) GetCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access-token")
		// fmt.Println("get cookie error", err)
		if err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		// fmt.Println("cookie is", cookie, cookie.Value, err)
		if err := a.CheckToken(c, cookie.Value); err != nil {
			return c.String(http.StatusUnauthorized, fmt.Sprintf("%v", err))
		}
		next(c)
		return nil
	}
}

func (a *Auth) CheckToken(c echo.Context, auth string) error {

	id, err := a.JwtUsecase.ParseTokenAndGetID(auth)
	if err != nil {
		return fmt.Errorf("get id from token error: %w", err)
	}

	role, err := a.JwtUsecase.ParseTokenAndGetRole(auth)
	if err != nil {
		return fmt.Errorf("get role from token error: %w", err)
	}

	iin, err := a.JwtUsecase.ParseTokenAndGetIIN(auth)
	if err != nil {
		return fmt.Errorf("get iin from token error: %w", err)
	}

	info := &domain.User{
		ID:   int64(id),
		IIN:  iin,
		Role: role,
	}

	c.Set("user", info)

	return nil
}

func (a *Auth) SetHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Content-Type", "application/json")
		next(c)
		return nil
	}
}
