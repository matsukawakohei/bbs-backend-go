package services

import (
	"bbs/models"
	"bbs/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(name string, email string, password string) error
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Signup(name string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repository.CreateUser(user)
}
