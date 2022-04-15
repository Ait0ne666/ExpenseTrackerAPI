package auth

import (
	"expense_tracker/models"
	"expense_tracker/pkg/infrastruct"
	"expense_tracker/pkg/jwt_auth"
	repository "expense_tracker/repositories"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	db  repository.AuthDAO
	jwt *jwt_auth.JwtAuth
}

func NewAuthService(
	db repository.AuthDAO,
	jwt *jwt_auth.JwtAuth,
) *AuthService {
	return &AuthService{
		db:  db,
		jwt: jwt,
	}
}

func (s *AuthService) Login(dto *models.LoginDTO) (*string, gin.H) {

	user, err := s.db.GetUserByLogin(dto.Login)

	if err != nil {
		return nil, infrastruct.ErrorForbidden
	}

	if user.Password != dto.Password {
		return nil, infrastruct.ErrorForbidden
	}

	jwt, err := s.jwt.GenerateJWT(user.ID)

	if user.Password != dto.Password {
		return nil, infrastruct.ErrorInternalServerError
	}

	return &jwt, nil

}

func (s *AuthService) GetUserById(id string) (*models.User, error) {

	user, err := s.db.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return user, nil

}
