package postgres

import (
	"context"
	"time"
	"transaction-service/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AccountRepo struct {
	Conn *pgxpool.Pool
}

func NewAccountRepo(conn *pgxpool.Pool) domain.AccountRepo {
	return &AccountRepo{Conn: conn}
}

func (ar *AccountRepo) CreateAccountRepo(acc *domain.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	if _, err := ar.Conn.Exec(ctx, "INSERT INTO accounts(iin, balance, number, registerDate) VALUES ($1, $2, $3, $4)",
		acc.IIN, acc.Balance, acc.AccountNumber, acc.RegisterDate); err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepo) DeleteAccountRepo(iin string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	_, err := ar.Conn.Exec(ctx, "DELETE FROM accounts WHERE iin=$1", iin)
	return err
}

func (ar *AccountRepo) DepositMoneyRepo(iin string, amount int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	_, err := ar.Conn.Exec(ctx, "UPDATE accounts SET balance = balance+$1 WHERE iin=$2", amount, iin)
	return err

}

func (ar *AccountRepo) TransferMoneyRepo(senderIIN, recipientIIN string, amount int64) error {
	
	return nil
}

func (ar *AccountRepo) GetAccountByIINRepo(iin string) (*domain.Account, error) {
	acc := &domain.Account{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	if err := ar.Conn.QueryRow(ctx, "SELECT id, iin, balance, number, registerDate FROM accounts WHERE iin=$1", iin).
		Scan(&acc.ID, &acc.IIN, &acc.Balance, &acc.AccountNumber, &acc.RegisterDate); err != nil {
		return nil, err
	}
	return acc, nil
}

func (ar *AccountRepo) GetAllAccountRepo() ([]*domain.Account, error) {
	account := &domain.Account{}
	allAcoount := []*domain.Account{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	rows, err := ar.Conn.Query(ctx, "SELECT id, iin, balance, number, registerDate FROM accounts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&account.ID, &account.IIN, &account.Balance, &account.AccountNumber, &account.RegisterDate); err != nil {
			return nil, err
		}
		allAcoount = append(allAcoount, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allAcoount, nil
}
