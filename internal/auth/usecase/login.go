package usecase

import (
	"context"

	"super-indo-api/internal/auth/domain"
	"super-indo-api/internal/auth/port"
)

type Login struct {
	repo     port.UserRepository
	tokenSvc port.TokenService
}

func NewLogin(repo port.UserRepository, tokenSvc port.TokenService) *Login {
	return &Login{repo: repo, tokenSvc: tokenSvc}
}

func (uc *Login) Execute(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	if !user.CheckPassword(password) {
		return nil, "", domain.ErrInvalidCredentials
	}

	token, err := uc.tokenSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
