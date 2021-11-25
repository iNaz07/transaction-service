package usecase

import (
	"fmt"
	"time"
	"transaction-service/domain"
	utils "transaction-service/utils"
)

type AccountUsecase struct {
	AccRepo domain.AccountRepo
}

func NewAccountUsecase(repo domain.AccountRepo) domain.AccountUsecase {
	return &AccountUsecase{AccRepo: repo}
}

//TODO: generate number
func (au *AccountUsecase) CreateAccount(iin string) error {
	acc := &domain.Account{
		IIN:           iin,
		RegisterDate:  time.Now().Format("2006-01-02 15:04:05"),
		Balance:       0,
		AccountNumber: utils.GenerateNumber(),
	}

	if err := au.AccRepo.CreateAccountRepo(acc); err != nil {
		return fmt.Errorf("create account error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) DeleteAccount(iin string) error {
	if err := au.AccRepo.DeleteAccountRepo(iin); err != nil {
		return fmt.Errorf("delete account error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) GetAccountByIIN(iin string) (*domain.Account, error) {
	account, err := au.AccRepo.GetAccountByIINRepo(iin)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return account, nil
}
func (au *AccountUsecase) GetAllAccount() ([]*domain.Account, error) {
	all, err := au.AccRepo.GetAllAccountRepo()
	if err != nil {
		return nil, fmt.Errorf("get all account err: %w", err)
	}
	return all, nil
}
