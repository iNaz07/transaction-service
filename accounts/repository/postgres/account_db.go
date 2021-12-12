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
	if _, err := ar.Conn.Exec(ctx, "INSERT INTO accounts(iin, userid, number, registerDate, balance) VALUES ($1, $2, $3, $4, $5)",
		acc.IIN, acc.UserID, acc.AccountNumber, acc.RegisterDate, acc.Balance); err != nil {
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

func (ar *AccountRepo) DepositMoneyRepo(deposit *domain.Deposit) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	tx, err := ar.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance+$1 and lasttransaction = $2 WHERE number=$3",
		deposit.Amount, deposit.Date, deposit.Number); err != nil {
		tx.Rollback(ctx)
		return err
	}
	if _, err := tx.Exec(ctx, "INSERT INTO deposits(iin, number, amount, date) VALUES ($1, $2, $3, $4)",
		deposit.IIN, deposit.Number, deposit.Amount, deposit.Date); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return err

}

func (ar *AccountRepo) TransferMoneyRepo(tr *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	tx, err := ar.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance+$1 AND lasttransaction = $2 WHERE number=$3", tr.Amount, tr.Date, tr.RecipientAccNumber); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance-$1 AND lasttransaction = $2 WHERE number = $3", tr.Amount, tr.Date, tr.SenderAccountNumber); err != nil {
		tx.Rollback(ctx)
		return err
	}
	if _, err := tx.Exec(ctx, "INSERT INTO transactions(sender, sender_number, recipient_number, recipient, amount, date) VALUES ($1, $2, $3, $4, $5, $6)",
		tr.SenderIIN, tr.SenderAccountNumber, tr.RecipientAccNumber, tr.RecipientIIN, tr.Amount, tr.Date); err != nil {
		tx.Rollback(ctx)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (ar *AccountRepo) GetAccountByIINRepo(iin string) ([]domain.Account, error) {
	acc := domain.Account{}
	userAccount := []domain.Account{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	rows, err := ar.Conn.Query(ctx, "SELECT id, iin, balance, number, registerDate FROM accounts WHERE iin=$1", iin)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&acc.ID, &acc.IIN, &acc.Balance, &acc.AccountNumber, &acc.RegisterDate); err != nil {
			return nil, err
		}
		userAccount = append(userAccount, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userAccount, nil
}

func (ar *AccountRepo) GetAccountByNumberRepo(number string) (*domain.Account, error) {
	acc := &domain.Account{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	if err := ar.Conn.QueryRow(ctx, "SELECT id, iin, balance, number, registerDate FROM accounts WHERE number = $1", number).
		Scan(&acc.ID, &acc.IIN, &acc.Balance, &acc.AccountNumber, &acc.RegisterDate); err != nil {
		return nil, err
	}
	return acc, nil
}

func (ar *AccountRepo) GetAllAccountRepo() ([]domain.Account, error) {
	account := domain.Account{}
	allAcoount := []domain.Account{}

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

func (ar *AccountRepo) GetAccountByUserIDRepo(userID int64) (*domain.Account, error) {
	acc := &domain.Account{}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := ar.Conn.QueryRow(ctx, "SELECT id, userid, iin FROM accounts WHERE userid = $1", userID).
		Scan(&acc.ID, &acc.UserID, &acc.IIN); err != nil {
		return nil, err
	}
	return acc, nil
}
