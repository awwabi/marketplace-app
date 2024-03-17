package bank_account

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"marketplace-app/entities"
	"marketplace-app/utils"
	"net/http"
)

type CreateBankAccountRequest struct {
	BankName          string `json:"bankName" validate:"required,min=5,max=15"`
	BankAccountName   string `json:"bankAccountName" validate:"required,min=5,max=15"`
	BankAccountNumber string `json:"bankAccountNumber" validate:"required,min=5,max=15"`
}

func (bh *BankAccountHandler) CreateBankAccount(c *gin.Context) {
	var request CreateBankAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	bankAccount := &entities.BankAccount{
		UserID:        c.Param("userID"),
		BankName:      request.BankName,
		AccountName:   request.BankAccountName,
		AccountNumber: request.BankAccountNumber,
	}

	if err := bh.repository.Create(bankAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank account created successfully"})
}
