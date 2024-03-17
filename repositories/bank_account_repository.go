package repositories

import (
	"database/sql"
	"marketplace-app/entities"
)

type BankAccountRepository struct {
	db *sql.DB
}

func NewBankAccountRepository(db *sql.DB) *BankAccountRepository {
	return &BankAccountRepository{db}
}

func (r *BankAccountRepository) Create(bankAccount *entities.BankAccount) error {
	query := `
		INSERT INTO bank_accounts (id, user_id, bank_name, account_name, account_number)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, bankAccount.ID, bankAccount.UserID, bankAccount.BankName, bankAccount.AccountName, bankAccount.AccountNumber)
	return err
}

func (r *BankAccountRepository) FindByID(id string) (*entities.BankAccount, error) {
	query := "SELECT * FROM bank_accounts WHERE id = $1 AND deleted_at IS NULL"
	row := r.db.QueryRow(query, id)

	var bankAccount entities.BankAccount
	if err := row.Scan(&bankAccount.ID, &bankAccount.UserID, &bankAccount.BankName, &bankAccount.AccountName, &bankAccount.AccountNumber, &bankAccount.CreatedAt, &bankAccount.UpdatedAt, &bankAccount.DeletedAt); err != nil {
		return nil, err
	}
	return &bankAccount, nil
}

func (r *BankAccountRepository) FindByUserID(userID string) ([]entities.BankAccount, error) {
	query := "SELECT * FROM bank_accounts WHERE user_id = $1 AND deleted_at IS NULL"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bankAccounts []entities.BankAccount
	for rows.Next() {
		var bankAccount entities.BankAccount
		if err := rows.Scan(&bankAccount.ID, &bankAccount.UserID, &bankAccount.BankName, &bankAccount.AccountName, &bankAccount.AccountNumber, &bankAccount.CreatedAt, &bankAccount.UpdatedAt, &bankAccount.DeletedAt); err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, bankAccount)
	}
	return bankAccounts, nil
}

func (r *BankAccountRepository) Update(bankAccount *entities.BankAccount) error {
	query := `
		UPDATE bank_accounts
		SET bank_name = $1, account_name = $2, account_number = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, bankAccount.BankName, bankAccount.AccountName, bankAccount.AccountNumber, bankAccount.ID)
	return err
}

func (r *BankAccountRepository) Delete(id string) error {
	query := "UPDATE bank_accounts SET deleted_at = now() AND updated_at = now() WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
