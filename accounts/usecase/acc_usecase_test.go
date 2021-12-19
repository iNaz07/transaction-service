package usecase_test

import (
	"github.com/bxcodec/faker"

	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ucase "transaction-service/accounts/usecase"
	"transaction-service/domain"
	"transaction-service/domain/mocks"
)

func TestCreateAccount(t *testing.T) {
	mockAccRepo := new(mocks.AccountRepo)
	iin := "940217450216"
	userid := int64(24)

	t.Run("success", func(t *testing.T) {
		mockAccRepo.On("CreateAccountRepo", mock.AnythingOfType("*domain.Account")).Return(nil).Once()

		u := ucase.NewAccountUsecase(mockAccRepo)
		err := u.CreateAccount(iin, userid)

		assert.NoError(t, err)

		mockAccRepo.AssertExpectations(t)
	})

}

func TestDepositMoney(t *testing.T) {
	var mockDeposit *domain.Deposit
	err := faker.FakeData(&mockDeposit)
	assert.NoError(t, err)

	mockRepo := new(mocks.AccountRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DepositMoneyRepo", mock.AnythingOfType("*domain.Deposit")).Return(nil)
		u := ucase.NewAccountUsecase(mockRepo)
		err := u.DepositMoney(mockDeposit.IIN, mockDeposit.Number, "1000")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransferMoney(t *testing.T) {
	var mockSender *domain.Account
	err := faker.FakeData(&mockSender)
	assert.NoError(t, err)

	var mockRecipient *domain.Account
	err = faker.FakeData(&mockRecipient)
	assert.NoError(t, err)

	mockRepo := new(mocks.AccountRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAccountByNumberRepo", mock.AnythingOfType("string")).Return(mockSender, nil).Once()
		mockRepo.On("GetAccountByNumberRepo", mock.AnythingOfType("string")).Return(mockRecipient, nil).Once()
		mockRepo.On("TransferMoneyRepo", mock.AnythingOfType("*domain.Transaction")).Return(nil).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		err := u.TransferMoney(mockSender.AccountNumber, mockRecipient.AccountNumber, strconv.Itoa(int(mockSender.Balance)-1))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})

	t.Run("error-failed", func(t *testing.T) {
		mockRepo.On("GetAccountByNumberRepo", mock.AnythingOfType("string")).Return(mockSender, nil).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		err := u.TransferMoney(mockSender.AccountNumber, mockRecipient.AccountNumber, strconv.Itoa(int(mockSender.Balance)+1))
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAccountByIIN(t *testing.T) {
	mockRepo := new(mocks.AccountRepo)
	mockAcc := &domain.Account{}
	err := faker.FakeData(mockAcc)
	assert.NoError(t, err)
	allMock := []domain.Account{}
	allMock = append(allMock, *mockAcc)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAccountByIINRepo", mock.AnythingOfType("string")).Return(allMock, nil).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAccountByIIN(mockAcc.IIN)
		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockRepo.On("GetAccountByIINRepo", mock.AnythingOfType("string")).Return([]domain.Account{}, errors.New("no rows in result set")).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAccountByIIN(mockAcc.IIN)
		assert.Error(t, err)
		assert.NotSame(t, []domain.Account{}, a)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAccountByNumber(t *testing.T) {
	mockRepo := new(mocks.AccountRepo)
	mockAcc := &domain.Account{}
	err := faker.FakeData(mockAcc)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAccountByNumberRepo", mock.AnythingOfType("string")).Return(mockAcc, nil).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAccountByNumber(mockAcc.AccountNumber)
		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockRepo.On("GetAccountByNumberRepo", mock.AnythingOfType("string")).Return(&domain.Account{}, errors.New("no rows in result set")).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAccountByNumber(mockAcc.AccountNumber)
		assert.Error(t, err)
		assert.NotSame(t, &domain.Account{}, a)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllAccount(t *testing.T) {
	mockRepo := new(mocks.AccountRepo)
	var mockAllAcc []domain.Account
	err := faker.FakeData(&mockAllAcc)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAllAccountRepo").Return(mockAllAcc, nil).Once()
		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAllAccount()
		assert.NoError(t, err)
		assert.NotNil(t, &domain.Account{}, a)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAccountByUserID(t *testing.T) {
	mockRepo := new(mocks.AccountRepo)
	mockAcc := &domain.Account{
		ID:              int64(25),
		UserID:          int64(10),
		IIN:             "940217450216",
		AccountNumber:   "KZ65777351365634",
		Balance:         0,
		RegisterDate:    time.Now().Format("2006-01-02 15:04:05"),
		LastTransaction: "",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAccountByUserIDRepo", mock.AnythingOfType("int64")).Return(mockAcc, nil).Once()

		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAccountByUserID(mockAcc.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockRepo.On("GetAccountByUserIDRepo", mock.AnythingOfType("int64")).Return(&domain.Account{}, errors.New("Unexpected")).Once()

		u := ucase.NewAccountUsecase(mockRepo)
		a, err := u.GetAccountByUserID(mockAcc.UserID)

		assert.Error(t, err)
		assert.NotSame(t, &domain.Account{}, a)

		mockRepo.AssertExpectations(t)
	})
}
