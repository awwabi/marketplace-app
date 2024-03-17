package entities

import (
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type ProductCondition string

const (
	PRODUCT_NEW    ProductCondition = "new"
	PRODUCT_SECOND ProductCondition = "second"
)

type Product struct {
	Model
	Name          string
	Price         decimal.Decimal
	ImageURL      string
	Stock         int
	Condition     ProductCondition
	Tags          pq.StringArray
	IsPurchasable bool
	UserID        string
}

func (pc ProductCondition) String() string {
	return string(pc)
}

type ProductFilter struct {
	UserOnly       *bool          `form:"userOnly"`
	Limit          *int           `form:"limit"`
	Offset         *int           `form:"offset"`
	Tags           pq.StringArray `form:"tags"`
	Condition      *string        `form:"condition" binding:"oneof=new second"`
	ShowEmptyStock *bool          `form:"showEmptyStock"`
	MaxPrice       *float64       `form:"maxPrice"`
	MinPrice       *float64       `form:"minPrice"`
	SortBy         *string        `form:"sortBy" binding:"oneof=price date"`
	OrderBy        *string        `form:"orderBy" binding:"oneof=asc desc"`
	Search         *string        `form:"search"`
	UserID         string         `form:"userId"`
}
