package product

import "marketplace-app/repositories"

type ProductHandler struct {
	repository *repositories.ProductRepository
}

func NewProductHandler(repository *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repository}
}
