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
func (au *AccountUsecase) CreateAccount(iin string, userid int64) error {
	acc := &domain.Account{
		IIN:           iin,
		UserID:        userid,
		Balance:       0,
		RegisterDate:  time.Now().Format("2006-01-02 15:04:05"),
		AccountNumber: utils.GenerateNumber(),
	}

	if err := au.AccRepo.CreateAccountRepo(acc); err != nil {
		return fmt.Errorf("create account error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) DepositMoney(iin, number, balance string) error {
	amount, err := strconv.Atoi(balance)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}
	//need to check acc number?
	deposit := &domain.Deposit{
		IIN:    iin,
		Number: number,
		Amount: int64(amount),
		Date:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := au.AccRepo.DepositMoneyRepo(deposit); err != nil {
		return fmt.Errorf("deposit money error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) TransferMoney(senderAccNum, recipientACCNum string, amount string) error {
	money, err := strconv.Atoi(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	acc, err := au.AccRepo.GetAccountByNumberRepo(senderAccNum) //check if account exists

	if err != nil {
		return fmt.Errorf("sender account doesn't exist: %w", err)
	}
	if acc.Balance <= int64(money) {
		return fmt.Errorf("not enough balance to transfer")
	}
	recAcc, err := au.AccRepo.GetAccountByNumberRepo(recipientACCNum)
	if err != nil {
		return fmt.Errorf("recipient account doesn't exist: %w", err)
	}

	transaction := &domain.Transaction{
		SenderIIN:           acc.IIN,
		SenderAccountNumber: senderAccNum,
		RecipientAccNumber:  recipientACCNum,
		RecipientIIN:        recAcc.IIN,
		Amount:              int64(money),
		Date:                time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := au.AccRepo.TransferMoneyRepo(transaction); err != nil {
		return fmt.Errorf("transaction money error: %w", err)
	}

	return nil
}

func (au *AccountUsecase) GetAccountByIIN(iin string) ([]domain.Account, error) {
	account, err := au.AccRepo.GetAccountByIINRepo(iin)
	if err != nil {
		return nil, fmt.Errorf("accounts not found: %w", err)
	}
	return account, nil
}

func (au *AccountUsecase) GetAccountByNumber(number string) (*domain.Account, error) {
	account, err := au.AccRepo.GetAccountByNumberRepo(number)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return account, nil
}

func (au *AccountUsecase) GetAllAccount() ([]domain.Account, error) {
	all, err := au.AccRepo.GetAllAccountRepo()
	if err != nil {
		return nil, fmt.Errorf("get all account err: %w", err)
	}
	return all, nil
}

func (au *AccountUsecase) GetAccountByUserID(userID int64) (*domain.Account, error) {
	acc, err := au.AccRepo.GetAccountByUserIDRepo(userID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return acc, nil
}
