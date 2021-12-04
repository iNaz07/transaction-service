package domain

type Account struct {
	ID            int64  `json:"id"`
	IIN           string `json:"iin"`
	AccountNumber string `json:"number"`
	Balance       int64  `json:"balance"`
	RegisterDate  string `json:"registerDate"`
}

type AccountRepo interface {
	CreateAccountRepo(acc *Account) error
	DeleteAccountRepo(iin string) error
	GetAccountByIINRepo(iin string) ([]*Account, error)
	GetAccountByNumberRepo(number string) (*Account, error)
	GetAllAccountRepo() ([]*Account, error)
	DepositMoneyRepo(iin string, amount int64) error
	TransferMoneyRepo(senderIIN, recipientIIN string, amount int64) error
}

type AccountUsecase interface {
	CreateAccount(iin string) error
	DeleteAccount(iin string) error
	GetAccountByIIN(iin string) ([]*Account, error)
	GetAccountByNumber(number string) (*Account, error)
	GetAllAccount() ([]*Account, error)
	DepositMoney(iin, balance string) error
	TransferMoney(senderIIN, recipientIIN string, amount int64) error
}
