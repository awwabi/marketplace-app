package bank_account

import "marketplace-app/repositories"

type BankAccountHandler struct {
	repository *repositories.BankAccountRepository
}

func NewBankAccountHandler(repository *repositories.BankAccountRepository) *BankAccountHandler {
	return &BankAccountHandler{repository}
}
