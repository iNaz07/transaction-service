package http

import (
	"fmt"
	"net/http"
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
	iin := c.FormValue("login") //temporary

	//TODO: check authorization from service-1
	//TODO: get IIN from token
	err := aH.AccUsecase.CreateAccount(iin)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("create account error: %v", err))
	}

	return c.String(http.StatusOK, "Account created. Your balance: 0")

}

func (aH *AccountHandler) OpenAccPage(c echo.Context) error {
	return c.HTML(http.StatusOK, `
	<html>
<head>
</head>
<body>
<form action="/create account" method="post">
	<label for="iin">Please enter your IIN:</label> <br>
	<input type="text" id="iin" name="iin"> <br>
	<input type="submit" value="Create an account">
</form>
</body>
</html>
	`)
}
