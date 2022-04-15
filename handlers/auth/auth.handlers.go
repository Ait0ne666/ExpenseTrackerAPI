package auth

import (
	"expense_tracker/models"
	"expense_tracker/pkg/infrastruct"
	"expense_tracker/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	s services.AuthService
}

func NewAuthHandlers(s services.AuthService) *AuthHandlers {
	return &AuthHandlers{
		s: s,
	}
}

func (h *AuthHandlers) Login(ctx *gin.Context) {
	var dto models.LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	jwt, err := h.s.Login(&dto)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}
