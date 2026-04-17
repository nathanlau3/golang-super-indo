package port

import (
	"context"

	"super-indo-api/internal/product/domain"
)

// Outbound port — kontrak yang dibutuhkan application dari infrastructure (adapter)

type ProductRepository interface {
	Save(ctx context.Context, product *domain.Product) error
	FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error)
	FindByID(ctx context.Context, id uint) (*domain.Product, error)
}
