package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marketplace-app/entities"
	"math/big"
	"net/http"
)

type GetProductsRequest struct {
	UserOnly       *bool    `form:"userOnly"`
	Limit          *int     `form:"limit"`
	Offset         *int     `form:"offset"`
	Tags           []string `form:"tags"`
	Condition      *string  `form:"condition" binding:"omitempty,oneof=new second"`
	ShowEmptyStock *bool    `form:"showEmptyStock"`
	MaxPrice       *float64 `form:"maxPrice"`
	MinPrice       *float64 `form:"minPrice"`
	SortBy         *string  `form:"sortBy" binding:"omitempty,oneof=price date"`
	OrderBy        *string  `form:"orderBy" binding:"omitempty,oneof=asc desc"`
	Search         *string  `form:"search"`
}

type GetProductsResponse struct {
	ProductID     string   `json:"productId"`
	Name          string   `json:"name"`
	Price         big.Int  `json:"price"`
	ImageURL      string   `json:"imageUrl"`
	Stock         int      `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
	PurchaseCount int      `json:"purchaseCount"`
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var request GetProductsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid query parameters: %s", err.Error())})
		return
	}

	filter := entities.ProductFilter{
		UserOnly:       request.UserOnly,
		Limit:          request.Limit,
		Offset:         request.Offset,
		Tags:           request.Tags,
		Condition:      request.Condition,
		ShowEmptyStock: request.ShowEmptyStock,
		MaxPrice:       request.MaxPrice,
		MinPrice:       request.MinPrice,
		SortBy:         request.SortBy,
		OrderBy:        request.OrderBy,
		Search:         request.Search,
		UserID:         c.GetString("userId"),
	}

	products, meta, err := h.repository.Search(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch products: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    mapProductsToResponse(products),
		"meta":    meta,
	})
}

func mapProductsToResponse(products []entities.Product) []GetProductsResponse {
	var response []GetProductsResponse
	for _, product := range products {
		response = append(response, GetProductsResponse{
			ProductID:     product.ID.String(),
			Name:          product.Name,
			Price:         *product.Price.BigInt(),
			ImageURL:      product.ImageURL,
			Stock:         product.Stock,
			Condition:     product.Condition.String(),
			Tags:          product.Tags,
			IsPurchasable: product.IsPurchasable,
			PurchaseCount: 2,
		})
	}

	return response
}
