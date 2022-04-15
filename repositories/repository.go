package repository

import (
	"expense_tracker/models"
	"expense_tracker/repositories/common"
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
	GetCategoryList(query string) (*[]models.Category, error)
	GetDayExpenses(date time.Time) ([]models.Expense, error)
}

type CommonDAO interface {
	InsertCurrencyRate(currency *models.CurrencyRate) error
	GetCurrencyRate(date time.Time) (*models.CurrencyRate, error)
}

type Repository struct {
	ExpensesDAO
	CommonDAO
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		ExpensesDAO: expenses.NewExpensesDAO(db),
		CommonDAO:   common.NewCommonDAO(db),
	}
}
