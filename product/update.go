package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"marketplace-app/entities"
	"marketplace-app/utils"
	"net/http"
	"strings"
)

type UpdateProductRequest struct {
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         float64  `json:"price" validate:"required,min=0"`
	ImageURL      string   `json:"imageUrl" validate:"required,url"`
	Condition     string   `json:"condition" validate:"required,oneof=new second"`
	Tags          []string `json:"tags" validate:"required"`
	IsPurchasable *bool    `json:"isPurchasable" validate:"required"`
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request data: %s", err.Error())})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	productID := c.Param("id")
	if err := uuid.Validate(productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.repository.FindByID(productID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to find product: %s", err.Error())})
		return
	}

	// check user permission
	userID := c.GetString("userID")
	if product.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this product"})
		return
	}

	product.Name = req.Name
	product.Price = decimal.NewFromFloat(req.Price)
	product.ImageURL = req.ImageURL
	product.Condition = entities.ProductCondition(req.Condition)
	product.Tags = req.Tags

	var purchasable bool
	if req.IsPurchasable != nil {
		purchasable = *req.IsPurchasable
	}
	product.IsPurchasable = purchasable

	if err := h.repository.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update product: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
