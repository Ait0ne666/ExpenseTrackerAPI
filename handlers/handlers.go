package handlers

import (
	expenses_handlers "expense_tracker/handlers/expenses"
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
}

type Handlers struct {
	ExpensesHandlers
}

func NewHandlers(s *services.Service) *Handlers {
	return &Handlers{
		ExpensesHandlers: expenses_handlers.NewExpensesHandlers(s),
	}
}
