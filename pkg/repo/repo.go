package repo

import (
	"github.com/parthvinchhi/rd-website/pkg/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}
