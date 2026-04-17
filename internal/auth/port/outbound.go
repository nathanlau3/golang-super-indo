package port

import (
	"context"

	"super-indo-api/internal/auth/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TokenService interface {
	GenerateToken(userID uint, email string) (string, error)
}
