package services

import (
	"expense_tracker/models"
	repository "expense_tracker/repositories"
	expenses "expense_tracker/services/Expenses"
	"time"
)

type ExpensesService interface {
	CreateCategory(title string) (*models.Category, error)
	DeleteCategory(id string) (*string, error)
	DeleteExpense(id string) (*string, error)
	CreateExpense(dto models.ExpenseDTO) (*models.Expense, error)
	GetDayTotalExpenses(date time.Time) (*models.DayTotalResulDTO, error)
	GetMonthExpenses(dto *models.MonthDTO) (*models.MonthExpensesDTO, error)
	GetMonthExpensesByCategory(date time.Time) (*models.MonthExpensesByCategoryDTO, error)
}

type Service struct {
	ExpensesService
}

func NewService(
	db *repository.Repository,
) *Service {
	return &Service{
		ExpensesService: expenses.NewExpensesService(db.ExpensesDAO),
	}
}
