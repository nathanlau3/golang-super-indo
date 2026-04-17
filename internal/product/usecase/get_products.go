package usecase

import (
	"context"

	"super-indo-api/internal/product/domain"
	"super-indo-api/internal/product/port"
)

type GetProducts struct {
	repo port.ProductRepository
}

func NewGetProducts(repo port.ProductRepository) *GetProducts {
	return &GetProducts{repo: repo}
}

func (uc *GetProducts) Execute(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	return uc.repo.FindAll(ctx, filter)
}
