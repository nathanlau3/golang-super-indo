package usecase

import (
	"context"

	"super-indo-api/internal/auth/domain"
	"super-indo-api/internal/auth/port"
)

type Register struct {
	repo port.UserRepository
}

func NewRegister(repo port.UserRepository) *Register {
	return &Register{repo: repo}
}

func (uc *Register) Execute(ctx context.Context, user *domain.User) error {
	return uc.repo.Save(ctx, user)
}
