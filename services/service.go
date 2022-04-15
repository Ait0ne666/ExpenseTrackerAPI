package services

import (
	"expense_tracker/models"
	repository "expense_tracker/repositories"
	common "expense_tracker/services/Common"
	expenses "expense_tracker/services/Expenses"
	"time"
)

type ExpensesService interface {
	CreateCategory(title string) (*models.Category, error)
	DeleteCategory(id string) (*string, error)
	DeleteExpense(id string) (*string, error)
	CreateExpense(dto models.ExpenseDTO) (*models.Expense, error)
	GetDayTotalExpenses(date time.Time, currency models.Currency) (*models.DayTotalResulDTO, error)
	GetMonthExpenses(dto *models.MonthDTO) (*models.MonthExpensesDTO, error)
	GetMonthExpensesByCategory(date time.Time, currency models.Currency) (*models.MonthExpensesByCategoryDTO, error)
	GetCategoryList(query string) (*[]models.Category, error)
	GetCurrencyRate() (*models.CurrencyRate, error)
}

type CommonService interface {
	GetCurrencyRate() (*models.CurrencyRate, error)
}

type Service struct {
	ExpensesService
	CommonService
}

func NewService(
	db *repository.Repository,
) *Service {
	common := common.NewCommonService(db.CommonDAO)
	return &Service{
		ExpensesService: expenses.NewExpensesService(db.ExpensesDAO, common),
		CommonService:   common,
	}
}
