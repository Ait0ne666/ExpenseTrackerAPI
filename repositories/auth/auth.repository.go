package auth

import (
	"expense_tracker/models"

	"gorm.io/gorm"
)

type AuthDAO struct {
	db *gorm.DB
}

func NewAuthDAO(db *gorm.DB) *AuthDAO {
	return &AuthDAO{db: db}
}

func (r *AuthDAO) GetUserById(id string) (*models.User, error) {

	user := models.User{}

	if err := r.db.Debug().Table("users").Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *AuthDAO) GetUserByLogin(login string) (*models.User, error) {

	user := models.User{}

	if err := r.db.Debug().Table("users").Where("login = ?", login).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil

}
