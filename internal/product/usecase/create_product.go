package usecase

import (
	"context"

	"super-indo-api/internal/product/domain"
	"super-indo-api/internal/product/port"
)

type CreateProduct struct {
	repo port.ProductRepository
}

func NewCreateProduct(repo port.ProductRepository) *CreateProduct {
	return &CreateProduct{repo: repo}
}

func (uc *CreateProduct) Execute(ctx context.Context, product *domain.Product) error {
	return uc.repo.Save(ctx, product)
}
