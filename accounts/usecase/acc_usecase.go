package usecase

import (
	"fmt"
	"strconv"
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

//TODO: check generate number
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

func (au *AccountUsecase) DepositMoney(iin, balance string) error {
	amount, err := strconv.Atoi(balance)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}
	if err := au.AccRepo.DepositMoneyRepo(iin, int64(amount)); err != nil {
		return fmt.Errorf("deposit money error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) TransferMoney(senderAccNum, recipientACCNum string, amount int64) error {
	acc, err := au.AccRepo.GetAccountByNumberRepo(senderAccNum) //check if account exists
	if err != nil {
		return fmt.Errorf("sender account doesn't exist: %w", err)
	}
	if acc.Balance <= amount {
		return fmt.Errorf("not enough balance to transfer")
	}
	if _, err := au.AccRepo.GetAccountByNumberRepo(recipientACCNum); err != nil {
		return fmt.Errorf("recipient account doesn't exist: %w", err)
	}

	if err := au.AccRepo.TransferMoneyRepo(senderAccNum, recipientACCNum, amount); err != nil {
		return fmt.Errorf("transfer money error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) DeleteAccount(iin string) error {
	if err := au.AccRepo.DeleteAccountRepo(iin); err != nil {
		return fmt.Errorf("delete account error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) GetAccountByIIN(iin string) ([]*domain.Account, error) {
	account, err := au.AccRepo.GetAccountByIINRepo(iin)
	if err != nil {
		return nil, fmt.Errorf("accounts not found: %w", err)
	}
	return account, nil
}

//TODO: unnessecary method, must be deleted
func (au *AccountUsecase) GetAccountByNumber(number string) (*domain.Account, error) {
	account, err := au.AccRepo.GetAccountByNumberRepo(number)
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
