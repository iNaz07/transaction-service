package http

import (
	"fmt"
	"log"
	"net/http"
	"transaction-service/domain"

	_mid "transaction-service/accounts/delivery/http/middleware"

	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	AccUsecase   domain.AccountUsecase
	TokenUsecase domain.JwtTokenUsecase
}

func NewAccountHandler(e *echo.Echo, acc domain.AccountUsecase, token domain.JwtTokenUsecase) {
	handler := &AccountHandler{AccUsecase: acc, TokenUsecase: token}

	accGroup := e.Group("/account")
	midd := _mid.InitAuth(token)
	accGroup.Use(midd.GetCookie, midd.SetHeader)

	accGroup.GET("/open", handler.OpenAccPage)
	accGroup.POST("/open", handler.OpenAcc)
	accGroup.POST("/deposit", handler.DepositAcc)
	accGroup.POST("/transfer", handler.TransferMoney)

	accGroup.GET("/info/:iin", handler.GetAccountInfo)
	accGroup.GET("/info", handler.GetAllAccountInfo) //not need?
}

func (aH *AccountHandler) TransferMoney(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		return c.String(http.StatusBadRequest, "cannot get meta info")
	}

	tr := &domain.Transaction{}
	if err := c.Bind(tr); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bind error: %v", err))
	}
	senderNumber := c.FormValue("sender number") //get from front
	acc, err := aH.AccUsecase.GetAccountByNumber(senderNumber)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("account not found: %v", err))
	}
	if acc.UserID != meta.ID {
		return c.String(http.StatusForbidden, "У вас недостаточно прав для данной операции")
	}

	if err := aH.AccUsecase.TransferMoney(senderNumber, tr.RecipientAccNumber, tr.Amount); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("transfer money error: %v", err))
	}

	return c.JSON(http.StatusOK, "Successful transfer")
}

func (aH *AccountHandler) DepositAcc(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		return c.String(http.StatusBadRequest, "cannot get meta info")
	}

	dep := &domain.Deposit{}
	if err := c.Bind(dep); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bind error: %v", err))
	}
	number := c.FormValue("sender number") //get from front
	acc, err := aH.AccUsecase.GetAccountByNumber(number)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("account not found: %v", err))
	}
	if acc.UserID != meta.ID {
		return c.String(http.StatusForbidden, "У вас недостаточно прав для данной операции")
	}

	balance := c.FormValue("amount") //temporary, should be dep.Amount

	if err := aH.AccUsecase.DepositMoney(acc.IIN, number, balance); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("deposit account error: %v", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("%v deposited into your account", balance))

}
func (aH *AccountHandler) GetAccountInfo(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		return c.String(http.StatusBadRequest, "cannot get meta info")
	}
	userAcc, err := aH.AccUsecase.GetAccountByUserID(meta.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("account not found: %v)", err))
	}

	iin := c.Param("iin") // must get iin from cookie of context
	log.Println("what is iin from auth service", iin)
	if meta.Role != "admin" {
		if userAcc.IIN != iin {
			return c.String(http.StatusForbidden, "У вас недостаточно прав")
		}
	}
	acc, err := aH.AccUsecase.GetAccountByIIN(iin)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("account not found: %v", err))
	}
	return c.JSON(http.StatusOK, acc)
}

func (aH *AccountHandler) GetAllAccountInfo(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		return c.String(http.StatusBadRequest, "cannot get meta info")
	}
	if meta.Role != "admin" {
		return c.String(http.StatusForbidden, "У вас недостаточно прав")
	}
	accounts, err := aH.AccUsecase.GetAllAccount()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("get all account error: %v", err))
	}
	return c.JSON(http.StatusOK, accounts)
}

func (aH *AccountHandler) OpenAcc(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		return c.String(http.StatusBadRequest, "cannot get meta info")
	}

	user := &domain.User{}
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bind user error: %v", err))
	}
	fmt.Println("are you here", meta, user)
	if meta.IIN != user.IIN {
		return c.String(http.StatusBadRequest, "invalid iin")
	}
	if err := aH.AccUsecase.CreateAccount(meta.IIN, meta.ID); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("create account error: %v", err))
	}
	return c.String(http.StatusOK, "Account created. Your balance: 0")
}

func (aH *AccountHandler) OpenAccPage(c echo.Context) error {
	fmt.Println("meta info", c.Get("user"))
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
