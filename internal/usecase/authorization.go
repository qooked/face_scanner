package usecase

import (
	"context"
	"faceScanner/pkg/password"
	"fmt"
)

type AuthUsecase struct {
	repository AuthRepository
}

type AuthRepository interface {
	GetUserCredentials(ctx context.Context, email string) (password string, err error)
	SaveUserCredentials(ctx context.Context, email string, hashPassword string) (err error)
}

func NewAuthUsecase(repository AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
	}
}

func (uc *AuthUsecase) GetUserCredentials(ctx context.Context, email string) (hashedPassword string, err error) {

	hashedPassword, err = uc.repository.GetUserCredentials(ctx, email)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetUserCredentials(...): %w", err)
		return "", err
	}

	return hashedPassword, nil
}

func (uc *AuthUsecase) SaveUserCredentials(ctx context.Context, email string, unhashedPassword string) (err error) {
	hashedPassword, err := password.HashPassword(unhashedPassword)
	if err != nil {
		err = fmt.Errorf("password.HashPassword(...): %w", err)
		return err
	}

	err = uc.repository.SaveUserCredentials(ctx, email, hashedPassword)
	if err != nil {
		err = fmt.Errorf("uc.repository.SaveUserCredentials(...): %w", err)
		return err
	}

	return nil
}
