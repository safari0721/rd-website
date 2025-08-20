package services

import (
	"github.com/parthvinchhi/rd-website/pkg/models"
	"github.com/parthvinchhi/rd-website/pkg/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repo.UserRepo
}

func (s *AuthService) Signup(name, agentID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{Name: name, AgentID: agentID, Password: string(hashedPassword)}
	return s.Repo.Create(user)
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}
