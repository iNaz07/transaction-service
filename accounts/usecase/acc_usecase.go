package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"transaction-service/domain"
	utils "transaction-service/utils"
)

type AccountUsecase struct {
	AccRepo        domain.AccountRepo
	ContextTimeout time.Duration
}

func NewAccountUsecase(repo domain.AccountRepo, time time.Duration) domain.AccountUsecase {
	return &AccountUsecase{AccRepo: repo, ContextTimeout: time}
}

//TODO: check generate number
func (au *AccountUsecase) CreateAccount(ctx context.Context, iin string, userid int64) error {
	acc := &domain.Account{
		IIN:           iin,
		UserID:        userid,
		Balance:       0,
		RegisterDate:  time.Now().Format("2006-01-02 15:04:05"),
		AccountNumber: utils.GenerateNumber(),
	}
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	if err := au.AccRepo.CreateAccountRepo(context, acc); err != nil {
		return fmt.Errorf("create account error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) DepositMoney(ctx context.Context, iin, number, balance string) error {
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
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	if err := au.AccRepo.DepositMoneyRepo(context, deposit); err != nil {
		return fmt.Errorf("deposit money error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) TransferMoney(ctx context.Context, senderAccNum, recipientACCNum string, amount string) error {
	money, err := strconv.Atoi(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	acc, err := au.AccRepo.GetAccountByNumberRepo(context, senderAccNum) //check if account exists

	if err != nil {
		return fmt.Errorf("sender account doesn't exist: %w", err)
	}
	if acc.Balance <= int64(money) {
		return fmt.Errorf("not enough balance to transfer")
	}
	recAcc, err := au.AccRepo.GetAccountByNumberRepo(context, recipientACCNum)
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

	if err := au.AccRepo.TransferMoneyRepo(context, transaction); err != nil {
		return fmt.Errorf("transaction money error: %w", err)
	}
	return nil
}

func (au *AccountUsecase) GetAccountByIIN(ctx context.Context, iin string) ([]domain.Account, error) {

	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	account, err := au.AccRepo.GetAccountByIINRepo(context, iin)
	if err != nil {
		return nil, fmt.Errorf("accounts not found: %w", err)
	}
	return account, nil
}

func (au *AccountUsecase) GetAccountByNumber(ctx context.Context, number string) (*domain.Account, error) {
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	account, err := au.AccRepo.GetAccountByNumberRepo(context, number)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return account, nil
}

func (au *AccountUsecase) GetAllAccount(ctx context.Context) ([]domain.Account, error) {
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	all, err := au.AccRepo.GetAllAccountRepo(context)
	if err != nil {
		return nil, fmt.Errorf("get all account err: %w", err)
	}
	return all, nil
}

func (au *AccountUsecase) GetAccountByUserID(ctx context.Context, userID int64) (*domain.Account, error) {
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	acc, err := au.AccRepo.GetAccountByUserIDRepo(context, userID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return acc, nil
}
