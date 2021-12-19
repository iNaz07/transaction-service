package domain

import (
	"context"
)

type Account struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"userid"`
	IIN             string `json:"iin"`
	AccountNumber   string `json:"number"`
	Balance         int64  `json:"balance"`
	RegisterDate    string `json:"registerDate"`
	LastTransaction string `json:"lasttransaction"`
}

type AccountRepo interface {
	CreateAccountRepo(ctx context.Context, acc *Account) error
	GetAccountByIINRepo(ctx context.Context, iin string) ([]Account, error)
	GetAccountByNumberRepo(ctx context.Context, number string) (*Account, error)
	GetAllAccountRepo(ctx context.Context) ([]Account, error)
	DepositMoneyRepo(ctx context.Context, deposit *Deposit) error
	TransferMoneyRepo(ctx context.Context, tr *Transaction) error
	GetAccountByUserIDRepo(ctx context.Context, userID int64) (*Account, error)
}

type AccountUsecase interface {
	CreateAccount(ctx context.Context, iin string, userID int64) error
	GetAccountByIIN(ctx context.Context, iin string) ([]Account, error)
	GetAccountByNumber(ctx context.Context, number string) (*Account, error)
	GetAllAccount(ctx context.Context) ([]Account, error)
	DepositMoney(ctx context.Context, iin, number, balance string) error
	TransferMoney(ctx context.Context, senderAccNum, recipientACCNum, amount string) error
	GetAccountByUserID(ctx context.Context, userID int64) (*Account, error)
}
