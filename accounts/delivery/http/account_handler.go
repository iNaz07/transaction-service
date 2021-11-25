package http

import (
	"transaction-service/domain"

	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	AccUsecase domain.AccountUsecase
}

func NewAccountHandler(e *echo.Echo, acc domain.AccountUsecase) {
	handler := &AccountHandler{AccUsecase: acc}

	e.GET("/account/open", handler.OpenAccPage)
	e.POST("/account/open", handler.OpenAcc)
}

func (aH *AccountHandler) OpenAcc(c echo.Context) error {
	return nil
}

func (aH *AccountHandler) OpenAccPage(c echo.Context) error {
	return nil
}
