package usecase

import (
	"errors"
	"user_service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo      domain.Authorization
	repoToken domain.TokenManager
}

func NewUseCase(repo domain.Authorization, token domain.TokenManager) *UseCase {
	return &UseCase{repo: repo, repoToken: token}
}

func (u *UseCase) RegisterUser(name, password string) (*domain.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.Register(domain.AuthorizationInput{
		Name:     name,
		Password: string(hashPassword),
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UseCase) LoginUser(name, password string) (*domain.Token, error) {
	user, err := u.repo.Login(domain.AuthorizationInput{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := u.repoToken.CreateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	return token, nil
}
