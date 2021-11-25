package domain

type Account struct {
	ID            int64  `json:"id"`
	IIN           string `json:"iin"`
	AccountNumber int64  `json:"number"`
	Balance       int64  `json:"balance"`
	RegisterDate  string `json:"registerDate"`
}

type AccountRepo interface {
	CreateAccountRepo(acc *Account) error
	DeleteAccountRepo(iin string) error
	GetAccountByIINRepo(iin string) (*Account, error)
	GetAllAccountRepo() ([]*Account, error)
}

type AccountUsecase interface {
	CreateAccount(iin string) error
	DeleteAccount(iin string) error
	GetAccountByIIN(iin string) (*Account, error)
	GetAllAccount() ([]*Account, error)
}
