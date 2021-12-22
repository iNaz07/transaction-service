package http

import (
	"fmt"
	"html/template"
	"io"

	"net/http"
	"transaction-service/domain"

	"github.com/rs/zerolog/log"

	_mid "transaction-service/accounts/delivery/http/middleware"

	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	AccUsecase   domain.AccountUsecase
	TokenUsecase domain.JwtTokenUsecase
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewAccountHandler(e *echo.Echo, acc domain.AccountUsecase, token domain.JwtTokenUsecase) {
	t := &Template{
		templates: template.Must(template.ParseGlob("../templates/*.html")),
	}
	e.Renderer = t
	handler := &AccountHandler{AccUsecase: acc, TokenUsecase: token}

	accGroup := e.Group("/account")
	midd := _mid.InitAuth(token)
	accGroup.Use(midd.GetCookie, midd.SetHeader)
	accGroup.GET("/", handler.HomePage)
	accGroup.GET("/open", handler.OpenAccPage)
	accGroup.GET("/deposit/:number", handler.DepositPage)
	accGroup.GET("/transfer/:number", handler.TransferPage)

	accGroup.POST("/open", handler.OpenAcc)
	accGroup.POST("/deposit/:number", handler.DepositAcc)
	accGroup.POST("/transfer/:number", handler.TransferMoney)

	accGroup.GET("/info/:iin/:service", handler.GetAccountInfo)
	accGroup.GET("/info", handler.GetAllAccountInfo) //no need?
}

func (aH *AccountHandler) HomePage(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		log.Err(domain.ErrorMetaNotFound).Msg("unauthorized")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied. Please authorize")
	}
	// return c.String(http.StatusOK, meta.IIN)
	return c.Render(http.StatusOK, "home.html", meta)
}

func (aH *AccountHandler) TransferPage(c echo.Context) error {
	number := c.Param("number")
	return c.Render(http.StatusOK, "transfer.html", number)
}

func (aH *AccountHandler) DepositPage(c echo.Context) error {
	number := c.Param("number")
	return c.Render(http.StatusOK, "deposit.html", number)
}

func (aH *AccountHandler) TransferMoney(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		log.Err(domain.ErrorMetaNotFound).Msg("unauthorized")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied. Please authorize")
	}

	senderNumber := c.Param("number")
	recipientNumber := c.FormValue("recipient")
	amount := c.FormValue("amount")

	ctx := c.Request().Context()
	acc, err := aH.AccUsecase.GetAccountByNumber(ctx, senderNumber)
	if err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", "unavailable account")
	}
	if acc.UserID != meta.ID {
		log.Log().Msg("proceeding acc doesn't belong to user")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied to proceed")
	}

	if err := aH.AccUsecase.TransferMoney(ctx, senderNumber, recipientNumber, amount); err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", logerr.Message)
	}
	// return c.String(http.StatusOK, fmt.Sprintf("%v KZT successfully transfered to account %v", amount, recipientNumber))
	return c.Render(http.StatusOK, "notify.html", fmt.Sprintf("%v KZT successfully transfered to account %v", amount, recipientNumber))
}

func (aH *AccountHandler) DepositAcc(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		log.Err(domain.ErrorMetaNotFound).Msg("unauthorized")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied. Please authorize")
	}

	number := c.Param("number")
	if number == "" {
		log.Log().Msg("sender account number is empty")
		return c.Render(http.StatusBadRequest, "notify.html", "Please choose account")
	}
	ctx := c.Request().Context()
	acc, err := aH.AccUsecase.GetAccountByNumber(ctx, number)
	if err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", logerr.Message)
	}
	if acc.UserID != meta.ID {
		log.Log().Msg("proceeding acc doesn't belong to user")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied to proceed")
	}

	balance := c.FormValue("amount")
	fmt.Println("balance", balance)
	if err := aH.AccUsecase.DepositMoney(ctx, acc.IIN, number, balance); err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", logerr.Message+" Please try again")
	}
	// return c.String(http.StatusOK, fmt.Sprintf("Account %v topped up amount: %v", number, balance))
	return c.Render(http.StatusOK, "notify.html", fmt.Sprintf("Account %v topped up amount: %v", number, balance))
}

func (aH *AccountHandler) GetAccountInfo(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		log.Err(domain.ErrorMetaNotFound).Msg("unauthorized")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied. Please authorize")
	}

	ctx := c.Request().Context()
	userAcc, err := aH.AccUsecase.GetAccountByUserID(ctx, meta.ID)
	if err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", "No available accounts to proceed")
	}

	iin := c.Param("iin")
	if meta.Role != "admin" {
		if userAcc.IIN != iin {
			log.Log().Msg("access denied")
			return c.Render(http.StatusForbidden, "notify.html", "Access denied")
		}
	}
	acc, err := aH.AccUsecase.GetAccountByIIN(ctx, iin)
	if err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", "No available accounts to proceed")
	}
	service := c.Param("service")
	if service == "auth" {
		return c.JSON(http.StatusOK, acc)
	}
	return c.Render(http.StatusOK, "info.html", acc)
}

func (aH *AccountHandler) GetAllAccountInfo(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		log.Err(domain.ErrorMetaNotFound).Msg("unauthorized")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied. Please authorize")
	}
	if meta.Role != "admin" {
		log.Log().Msg("user role not admin")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied")
	}

	ctx := c.Request().Context()
	accounts, err := aH.AccUsecase.GetAllAccount(ctx)
	if err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", logerr.Message)
	}
	// return c.JSON(http.StatusOK, accounts)
	return c.Render(http.StatusOK, "info.html", accounts)

}

func (aH *AccountHandler) OpenAcc(c echo.Context) error {
	meta, ok := c.Get("user").(*domain.User)
	if !ok {
		log.Err(domain.ErrorMetaNotFound).Msg("unauthorized")
		return c.Render(http.StatusForbidden, "notify.html", "Access denied. Please authorize")
	}

	user := &domain.User{
		IIN: c.FormValue("iin"),
	}

	if meta.IIN != user.IIN {
		log.Log().Msg("invalid IIN")
		return c.Render(http.StatusBadRequest, "notify.html", "Invalid IIN")
	}

	ctx := c.Request().Context()
	if err := aH.AccUsecase.CreateAccount(ctx, meta.IIN, meta.ID); err != nil {
		logerr := err.(*domain.LogError)
		log.Err(logerr.Err).Msg(logerr.Message)
		return c.Render(logerr.Code, "notify.html", "Unexpected error occured. Please try again")
	}
	// return c.String(http.StatusOK, "Account created. Your balance: 0")
	return c.Render(http.StatusOK, "notify.html", "Account created. Your balance: 0")
}

func (aH *AccountHandler) OpenAccPage(c echo.Context) error {
	return c.Render(http.StatusOK, "open.html", nil)
}
