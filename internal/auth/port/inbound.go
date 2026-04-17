package port

import (
	"context"

	"super-indo-api/internal/auth/domain"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, user *domain.User) error
}

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (*domain.User, string, error)
}
