package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"marketplace-app/entities"
	"marketplace-app/utils"
	"net/http"
)

type CreateProductRequest struct {
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         float64  `json:"price" validate:"required,min=0"`
	ImageURL      string   `json:"imageUrl" validate:"required,url"`
	Stock         int      `json:"stock" validate:"required,min=0"`
	Condition     string   `json:"condition" validate:"required,oneof=new second"`
	Tags          []string `json:"tags" validate:"required"`
	IsPurchasable *bool    `json:"isPurchasable" validate:"required"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request data: %s", err.Error())})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	userID := c.GetString("userID")

	newProduct := &entities.Product{
		Model:         entities.Model{ID: uuid.New()},
		Name:          req.Name,
		Price:         decimal.NewFromFloat(req.Price),
		ImageURL:      req.ImageURL,
		Stock:         req.Stock,
		Condition:     entities.ProductCondition(req.Condition),
		Tags:          pq.StringArray(req.Tags),
		IsPurchasable: *req.IsPurchasable,
		UserID:        userID,
	}

	if err := h.repository.Create(newProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create product: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
}
