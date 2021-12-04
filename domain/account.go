package domain

type Account struct {
	ID              int64  `json:"id"`
	IIN             string `json:"iin"`
	AccountNumber   string `json:"number"`
	Balance         int64  `json:"balance"`
	RegisterDate    string `json:"registerDate"`
	LastTransaction string `json:"lasttransaction"`
}

type AccountRepo interface {
	CreateAccountRepo(acc *Account) error
	DeleteAccountRepo(iin string) error
	GetAccountByIINRepo(iin string) ([]*Account, error)
	GetAccountByNumberRepo(number string) (*Account, error)
	GetAllAccountRepo() ([]*Account, error)
	DepositMoneyRepo(deposit *Deposit) error
	TransferMoneyRepo(tr *Transaction) error
}

type AccountUsecase interface {
	CreateAccount(iin string) error
	DeleteAccount(iin string) error
	GetAccountByIIN(iin string) ([]*Account, error)
	GetAccountByNumber(number string) (*Account, error)
	GetAllAccount() ([]*Account, error)
	DepositMoney(iin, number, balance string) error
	TransferMoney(senderIIN, recipientIIN string, amount int64) error
}
