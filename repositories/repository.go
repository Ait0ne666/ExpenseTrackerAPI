package repository

import (
	"expense_tracker/models"
	"expense_tracker/repositories/auth"
	"expense_tracker/repositories/common"
	"expense_tracker/repositories/expenses"
	"time"

	"gorm.io/gorm"
)

type ExpensesDAO interface {
	CreateCategory(title string, userID string) (*models.Category, error)
	DeleteCategory(id string, userID string) error
	UpsertExpense(expense *models.Expense, userID string) error
	DeleteExpense(id string, userID string) error
	GetDayTotalExpenses(date time.Time, userID string) (*float64, error)
	GetMonthTotalExpenses(dto *models.MonthDTO, userID string) (*float64, error)
	GetMonthExpenses(dto *models.MonthDTO, userID string) ([]models.Expense, error)
	GetMonthExpensesByCategory(date time.Time, userID string) ([]models.MonthCategory, error)
	GetCategoryList(query string, userID string) (*[]models.Category, error)
	GetDayExpenses(date time.Time, userID string) ([]models.Expense, error)
}

type CommonDAO interface {
	InsertCurrencyRate(currency *models.CurrencyRate) error
	GetCurrencyRate(date time.Time) (*models.CurrencyRate, error)
}

type AuthDAO interface {
	GetUserByLogin(login string) (*models.User, error)
	GetUserById(id string) (*models.User, error)
}

type Repository struct {
	ExpensesDAO
	CommonDAO
	AuthDAO
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		ExpensesDAO: expenses.NewExpensesDAO(db),
		CommonDAO:   common.NewCommonDAO(db),
		AuthDAO:     auth.NewAuthDAO(db),
	}
}
