package repository

import (
	"expense_tracker/models"
	"expense_tracker/repositories/expenses"
	"time"

	"gorm.io/gorm"
)

type ExpensesDAO interface {
	CreateCategory(title string) (*models.Category, error)
	DeleteCategory(id string) error
	UpsertExpense(expense *models.Expense) error
	DeleteExpense(id string) error
	GetDayTotalExpenses(date time.Time) (*float64, error)
	GetMonthTotalExpenses(dto *models.MonthDTO) (*float64, error)
	GetMonthExpenses(dto *models.MonthDTO) ([]models.Expense, error)
	GetMonthExpensesByCategory(date time.Time) ([]models.MonthCategory, error)
}

type Repository struct {
	ExpensesDAO
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		ExpensesDAO: expenses.NewExpensesDAO(db),
	}
}
