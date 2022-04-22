package services

import (
	"expense_tracker/models"
	"expense_tracker/pkg/jwt_auth"
	repository "expense_tracker/repositories"
	auth "expense_tracker/services/Auth"
	common "expense_tracker/services/Common"
	expenses "expense_tracker/services/Expenses"
	"time"

	"github.com/gin-gonic/gin"
)

type ExpensesService interface {
	CreateCategory(title string, userID string) (*models.Category, error)
	DeleteCategory(id string, userID string) (*string, error)
	DeleteExpense(id string, userID string) (*string, error)
	CreateExpense(dto models.ExpenseDTO, userID string) (*models.Expense, error)
	GetDayTotalExpenses(date time.Time, currency models.Currency, userID string) (*models.DayTotalResulDTO, error)
	GetMonthExpenses(dto *models.MonthDTO, userID string) (*models.MonthExpensesDTO, error)
	GetMonthExpensesByCategory(date time.Time, currency models.Currency, userID string) (*models.MonthExpensesByCategoryDTO, error)
	GetCategoryList(query string, userID string) (*[]models.Category, error)
	GetCurrencyRate() (*models.CurrencyRate, error)
	SyncDatabase(syncData models.SyncDTO, userID string) (*models.SyncResultDTO, error)
}

type CommonService interface {
	GetCurrencyRate() (*models.CurrencyRate, error)
}

type AuthService interface {
	Login(dto *models.LoginDTO) (*string, gin.H)
	GetUserById(id string) (*models.User, error)
}

type Service struct {
	ExpensesService
	CommonService
	AuthService
}

func NewService(
	db *repository.Repository,
	jwt *jwt_auth.JwtAuth,
) *Service {
	common := common.NewCommonService(db.CommonDAO)
	return &Service{
		ExpensesService: expenses.NewExpensesService(db.ExpensesDAO, common),
		AuthService:     auth.NewAuthService(db.AuthDAO, jwt),
		CommonService:   common,
	}
}
