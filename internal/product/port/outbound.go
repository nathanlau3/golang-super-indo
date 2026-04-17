package port

import (
	"context"

	"super-indo-api/internal/product/domain"
	"super-indo-api/pkg/common"
)

type ProductRepository interface {
	Save(ctx context.Context, product *domain.Product) error
	FindAll(ctx context.Context, filter common.Filter) ([]domain.Product, int64, error)
	FindByID(ctx context.Context, id uint) (*domain.Product, error)
}
