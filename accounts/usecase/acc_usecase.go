package usecase

import (
	"context"
	"net/http"

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

func (au *AccountUsecase) CreateAccount(ctx context.Context, iin string, userid int64) error {
	acc := &domain.Account{
		IIN:             iin,
		UserID:          userid,
		Balance:         0,
		RegisterDate:    time.Now().Format("2006-01-02 15:04:05"),
		AccountNumber:   utils.GenerateNumber(),
		LastTransaction: "no data",
	}
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	if err := au.AccRepo.CreateAccountRepo(context, acc); err != nil {
		return &domain.LogError{"create account error", err, http.StatusInternalServerError}
	}
	return nil
}

func (au *AccountUsecase) DepositMoney(ctx context.Context, iin, number, balance string) error {
	amount, err := strconv.Atoi(balance)
	if err != nil {
		return &domain.LogError{"invalid amount.", err, http.StatusBadRequest}
	}
	if _, err := au.AccRepo.GetAccountByNumberRepo(ctx, number); err != nil {
		return &domain.LogError{fmt.Sprintf("account: %s not found", number), err, http.StatusBadRequest}
	}
	deposit := &domain.Deposit{
		IIN:    iin,
		Number: number,
		Amount: int64(amount),
		Date:   time.Now().Format("2006-01-02 15:04:05"),
	}
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()
	if err := au.AccRepo.DepositMoneyRepo(context, deposit); err != nil {
		return &domain.LogError{"top up account error.", err, http.StatusInternalServerError}
	}
	return nil
}

func (au *AccountUsecase) TransferMoney(ctx context.Context, senderAccNum, recipientACCNum string, amount string) error {
	money, err := strconv.Atoi(amount)
	if err != nil {
		return &domain.LogError{"invalid amount", err, http.StatusBadRequest}
	}

	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()

	acc, err := au.AccRepo.GetAccountByNumberRepo(context, senderAccNum) //check if account exists
	if err != nil {
		return &domain.LogError{"sender account not found", err, http.StatusNotFound}
	}

	if acc.Balance <= int64(money) {
		return &domain.LogError{"not enough balance to transfer", fmt.Errorf("insufficient funds"), http.StatusBadRequest}
	}

	recAcc, err := au.AccRepo.GetAccountByNumberRepo(context, recipientACCNum)
	if err != nil {
		return &domain.LogError{"recipient account not found", err, http.StatusNotFound}
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
		return &domain.LogError{"error occured while transaction", err, http.StatusInternalServerError}
	}
	return nil
}

func (au *AccountUsecase) GetAccountByIIN(ctx context.Context, iin string) ([]domain.Account, error) {

	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()

	account, err := au.AccRepo.GetAccountByIINRepo(context, iin)
	if err != nil {
		return nil, &domain.LogError{"account not found", err, http.StatusNotFound}
	}
	return account, nil
}

func (au *AccountUsecase) GetAccountByNumber(ctx context.Context, number string) (*domain.Account, error) {
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()

	account, err := au.AccRepo.GetAccountByNumberRepo(context, number)
	if err != nil {
		return nil, &domain.LogError{"account not found", err, http.StatusNotFound}
	}
	return account, nil
}

func (au *AccountUsecase) GetAllAccount(ctx context.Context) ([]domain.Account, error) {
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()

	all, err := au.AccRepo.GetAllAccountRepo(context)
	if err != nil {
		return nil, &domain.LogError{"get all account error", err, http.StatusInternalServerError}
	}
	return all, nil
}

func (au *AccountUsecase) GetAccountByUserID(ctx context.Context, userID int64) (*domain.Account, error) {
	context, cancel := context.WithTimeout(ctx, au.ContextTimeout)
	defer cancel()

	acc, err := au.AccRepo.GetAccountByUserIDRepo(context, userID)
	if err != nil {
		return nil, &domain.LogError{"account not found", err, http.StatusNotFound}
	}
	return acc, nil
}
