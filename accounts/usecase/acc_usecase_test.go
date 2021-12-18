package usecase_test

import (
	// "errors"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// "transaction-service/domain"
	ucase "transaction-service/accounts/usecase"
	"transaction-service/domain/mocks"
	// utils "transaction-service/utils"
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
