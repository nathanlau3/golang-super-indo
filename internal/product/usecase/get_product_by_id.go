package usecase

import (
	"context"

	"super-indo-api/internal/product/domain"
	"super-indo-api/internal/product/port"
)

type GetProductByID struct {
	repo port.ProductRepository
}

func NewGetProductByID(repo port.ProductRepository) *GetProductByID {
	return &GetProductByID{repo: repo}
}

func (uc *GetProductByID) Execute(ctx context.Context, id uint) (*domain.Product, error) {
	return uc.repo.FindByID(ctx, id)
}
