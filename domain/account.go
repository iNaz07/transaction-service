package domain

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
	CreateAccountRepo(acc *Account) error
	GetAccountByIINRepo(iin string) ([]Account, error)
	GetAccountByNumberRepo(number string) (*Account, error)
	GetAllAccountRepo() ([]Account, error)
	DepositMoneyRepo(deposit *Deposit) error
	TransferMoneyRepo(tr *Transaction) error
	GetAccountByUserIDRepo(userID int64) (*Account, error)
}

type AccountUsecase interface {
	CreateAccount(iin string, userID int64) error
	GetAccountByIIN(iin string) ([]Account, error)
	GetAccountByNumber(number string) (*Account, error)
	GetAllAccount() ([]Account, error)
	DepositMoney(iin, number, balance string) error
	TransferMoney(senderIIN, recipientIIN, amount string) error
	GetAccountByUserID(userID int64) (*Account, error)
}
