package port

import (
	"context"

	"super-indo-api/internal/product/domain"
)

type CreateProductUseCase interface {
	Execute(ctx context.Context, product *domain.Product) error
}

type GetProductsUseCase interface {
	Execute(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error)
}

type GetProductByIDUseCase interface {
	Execute(ctx context.Context, id uint) (*domain.Product, error)
}
