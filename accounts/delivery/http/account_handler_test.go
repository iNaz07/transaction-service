package http_test

import (
	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	// "strconv"
	"strings"
	"testing"
	accHTTP "transaction-service/accounts/delivery/http"
	"transaction-service/domain"
	"transaction-service/domain/mocks"
	// utils "transaction-service/utils"
)

func TestHomePage(t *testing.T) {

	var mockNewUser *domain.User
	err := faker.FakeData(&mockNewUser)
	assert.NoError(t, err)

	mockUCase := new(mocks.AccountUsecase)

	e := echo.New()

	req, err := http.NewRequest(echo.GET, "/account/", strings.NewReader(""))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.Set("user", mockNewUser)

	handler := accHTTP.AccountHandler{
		AccUsecase: mockUCase,
	}
	err = handler.HomePage(c)
	require.NoError(t, err)

	mockUCase.AssertExpectations(t)
}

func TestOpenAcc(t *testing.T) {

	var mockNewUser *domain.User
	err := faker.FakeData(&mockNewUser)
	assert.NoError(t, err)

	mockUCase := new(mocks.AccountUsecase)
	mockUCase.On("CreateAccount", mockNewUser.IIN, mockNewUser.ID).Return(nil)

	e := echo.New()

	req, err := http.NewRequest(echo.POST, "/account/open?iin="+mockNewUser.IIN, strings.NewReader(""))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.Set("user", mockNewUser)

	handler := accHTTP.AccountHandler{
		AccUsecase: mockUCase,
	}
	err = handler.OpenAcc(c)
	require.NoError(t, err)

	mockUCase.AssertExpectations(t)
}

func TestTransferMoney(t *testing.T) {

	var mockAcc domain.Account
	err := faker.FakeData(&mockAcc)
	assert.NoError(t, err)
	mockUser := &domain.User{
		ID:   mockAcc.UserID,
		IIN:  mockAcc.IIN,
		Role: "user",
	}

	mockUCase := new(mocks.AccountUsecase)
	mockUCase.On("GetAccountByNumber", mockAcc.AccountNumber).Return(&mockAcc, nil)
	mockUCase.On("TransferMoney", mockAcc.AccountNumber, "KZ05777326581634", "500").Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/account/transfer/:number?recipient=KZ05777326581634&amount=500", strings.NewReader(""))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.Set("user", mockUser)
	c.SetPath("/account/transfer/:number")
	c.SetParamNames("number")
	c.SetParamValues(mockAcc.AccountNumber)

	handler := accHTTP.AccountHandler{
		AccUsecase: mockUCase,
	}
	err = handler.TransferMoney(c)
	require.NoError(t, err)

	mockUCase.AssertExpectations(t)
}

func TestDepositAcc(t *testing.T) {

	var mockAcc domain.Account
	err := faker.FakeData(&mockAcc)
	assert.NoError(t, err)
	mockUser := &domain.User{
		ID:   mockAcc.UserID,
		IIN:  mockAcc.IIN,
		Role: "admin",
	}

	mockUCase := new(mocks.AccountUsecase)
	mockUCase.On("GetAccountByNumber", mockAcc.AccountNumber).Return(&mockAcc, nil)
	mockUCase.On("DepositMoney", mockAcc.IIN, mockAcc.AccountNumber, "500").Return(nil)

	e := echo.New()

	req, err := http.NewRequest(echo.POST, "/account/deposit/:number?amount=500", strings.NewReader(""))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.Set("user", mockUser)
	c.SetPath("/account/deposit/:number")
	c.SetParamNames("number")
	c.SetParamValues(mockAcc.AccountNumber)

	handler := accHTTP.AccountHandler{
		AccUsecase: mockUCase,
	}
	err = handler.DepositAcc(c)
	require.NoError(t, err)

	mockUCase.AssertExpectations(t)
}

func TestGetAccountInfo(t *testing.T) {
	var mockAcc domain.Account
	err := faker.FakeData(&mockAcc)
	assert.NoError(t, err)

	mockUser := &domain.User{
		ID:   mockAcc.UserID,
		IIN:  mockAcc.IIN,
		Role: "user",
	}
	mockAllAcc := []domain.Account{}
	mockAllAcc = append(mockAllAcc, mockAcc)
	mockUCase := new(mocks.AccountUsecase)
	mockUCase.On("GetAccountByUserID", mockUser.ID).Return(&mockAcc, nil)
	mockUCase.On("GetAccountByIIN", mockUser.IIN).Return(mockAllAcc, nil)

	req, err := http.NewRequest(echo.GET, "/account/info/"+mockAcc.IIN, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	c.Set("user", mockUser)
	c.SetPath("/account/info/:iin")
	c.SetParamNames("iin")
	c.SetParamValues(mockAcc.IIN)

	handler := accHTTP.AccountHandler{
		AccUsecase: mockUCase,
	}
	err = handler.GetAccountInfo(c)
	require.NoError(t, err)

	mockUCase.AssertExpectations(t)
}
