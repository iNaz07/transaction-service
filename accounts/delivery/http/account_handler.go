package http

import (
	"fmt"
	"net/http"
	"strconv"
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
	e.POST("/account/deposit", handler.DepositAcc)
	e.GET("/account/info/:iin", handler.GetAccountInfo)
	e.POST("/account/transfer", handler.TransferMoney)
}

//TODO: check auth
func (aH *AccountHandler) TransferMoney(c echo.Context) error {
	senderIIN := c.FormValue("iin")                    //must get from cookie
	recipientIIN := c.FormValue("recipiin")            //get from front
	amount, err := strconv.Atoi(c.FormValue("amount")) //get from front
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := aH.AccUsecase.TransferMoney(senderIIN, recipientIIN, int64(amount)); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("transfer money error: %v", err))
	}
	return c.String(http.StatusOK, "Money successfully transfered")
}

func (aH *AccountHandler) GetAccountInfo(c echo.Context) error {
	iin := c.FormValue("iin") // must get iin from cookie of context
	acc, err := aH.AccUsecase.GetAccountByIIN(iin)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("account not found: %v", err))
	}
	return c.JSON(http.StatusOK, acc)
}

func (aH *AccountHandler) DepositAcc(c echo.Context) error {
	balance := c.FormValue("amount") //temporary
	iin := c.FormValue("iin")        //get from cookie
	if err := aH.AccUsecase.DepositMoney(iin, balance); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("deposit account error: %v", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("%v deposited into your account", balance))

}

func (aH *AccountHandler) OpenAcc(c echo.Context) error {
	iin := c.FormValue("iin") //temporary

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
