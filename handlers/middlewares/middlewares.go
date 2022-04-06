package middleware

import (
	"expense_tracker/pkg/infrastruct"
	"expense_tracker/pkg/jwt_auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	jwt *jwt_auth.JwtAuth
}

func NewMiddlewares(jwt *jwt_auth.JwtAuth) *Middlewares {
	return &Middlewares{jwt: jwt}
}

func (h *Middlewares) AuthGuard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := h.jwt.GetClaimsByRequest(ctx)
		if err != nil {
			ctx.JSON(http.StatusForbidden, infrastruct.ErrorForbidden)
			ctx.Abort()
			return
		}

		ctx.Set("user", claims.UserID)

	}
}
