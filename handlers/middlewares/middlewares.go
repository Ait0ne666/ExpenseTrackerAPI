package middleware

import (
	"expense_tracker/pkg/infrastruct"
	"expense_tracker/pkg/jwt_auth"
	"expense_tracker/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	jwt *jwt_auth.JwtAuth
	s   services.AuthService
}

func NewMiddlewares(jwt *jwt_auth.JwtAuth, auth services.AuthService) *Middlewares {
	return &Middlewares{jwt: jwt, s: auth}
}

func (h *Middlewares) AuthGuard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := h.jwt.GetClaimsByRequest(ctx)

		if err != nil {
			ctx.JSON(http.StatusForbidden, infrastruct.ErrorForbidden)
			ctx.Abort()
			return
		}

		user, err := h.s.GetUserById(claims.UserID)

		if err != nil {
			ctx.JSON(http.StatusForbidden, infrastruct.ErrorForbidden)
			ctx.Abort()
			return
		}

		ctx.Set("user", user.ID)

	}
}
