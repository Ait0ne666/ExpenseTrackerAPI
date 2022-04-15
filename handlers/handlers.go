package handlers

import (
	expenses_handlers "expense_tracker/handlers/expenses"

	middleware "expense_tracker/handlers/middlewares"
	"expense_tracker/pkg/jwt_auth"
	"expense_tracker/services"

	"github.com/gin-gonic/gin"
)

type ExpensesHandlers interface {
	CreateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
	DeleteExpense(ctx *gin.Context)
	CreateExpense(ctx *gin.Context)
	GetDayTotalExpenses(ctx *gin.Context)
	GetMonthExpenses(ctx *gin.Context)
	GetMonthExpensesByCategory(ctx *gin.Context)
	GetCategoryList(ctx *gin.Context)
	GetCurrencyRate(ctx *gin.Context)
}

type Middlewares interface {
	AuthGuard() gin.HandlerFunc
}

type Handlers struct {
	ExpensesHandlers
	Middlewares
}

func NewHandlers(s *services.Service, jwt *jwt_auth.JwtAuth) *Handlers {
	return &Handlers{
		ExpensesHandlers: expenses_handlers.NewExpensesHandlers(s.ExpensesService),
		Middlewares:      middleware.NewMiddlewares(jwt),
	}
}
