package expenses

import (
	"expense_tracker/models"
	repository "expense_tracker/repositories"
	"time"
)

type ExpensesService struct {
	db repository.ExpensesDAO
}

func NewExpensesService(
	db repository.ExpensesDAO,
) *ExpensesService {
	return &ExpensesService{
		db: db,
	}
}

func (s *ExpensesService) CreateCategory(title string) (*models.Category, error) {

	category, err := s.db.CreateCategory(title)

	if err != nil {
		return nil, err
	}

	return category, nil

}

func (s *ExpensesService) DeleteCategory(id string) (*string, error) {

	err := s.db.DeleteCategory(id)

	if err != nil {
		return nil, err
	}

	return &id, nil

}

func (s *ExpensesService) DeleteExpense(id string) (*string, error) {

	err := s.db.DeleteExpense(id)

	if err != nil {
		return nil, err
	}

	return &id, nil

}

func (s *ExpensesService) CreateExpense(dto models.ExpenseDTO) (*models.Expense, error) {

	expense := models.Expense{}

	if dto.CategoryID == nil {

		category, err := s.db.CreateCategory(dto.CategoryName)

		if err != nil {
			return nil, err
		}

		expense.CategoryID = category.ID
	} else {
		expense.CategoryID = *dto.CategoryID
	}

	expense.Date = dto.Date
	if dto.ID != nil {
		expense.ID = *dto.ID
	}
	expense.Title = dto.Title
	expense.Amount = dto.Amount

	err := s.db.UpsertExpense(&expense)

	if err != nil {
		return nil, err
	}

	return &expense, nil

}

func (s *ExpensesService) GetDayTotalExpenses(date time.Time) (*models.DayTotalResulDTO, error) {

	total, err := s.db.GetDayTotalExpenses(date)

	if err != nil {
		return nil, err
	}

	monthTotal, err := s.db.GetMonthTotalExpenses(&models.MonthDTO{Date: date})

	if err != nil {
		return nil, err
	}

	result := models.DayTotalResulDTO{DayTotal: *total, MonthTotal: *monthTotal}

	return &result, nil

}

func (s *ExpensesService) GetMonthExpenses(dto *models.MonthDTO) (*models.MonthExpensesDTO, error) {

	monthTotal, err := s.db.GetMonthTotalExpenses(dto)

	if err != nil {
		return nil, err
	}

	expenses, err := s.db.GetMonthExpenses(dto)

	if err != nil {
		return nil, err
	}
	result := models.MonthExpensesDTO{MonthTotal: *monthTotal, Expenses: expenses}

	return &result, nil

}

func (s *ExpensesService) GetMonthExpensesByCategory(date time.Time) (*models.MonthExpensesByCategoryDTO, error) {

	monthTotal, err := s.db.GetMonthTotalExpenses(&models.MonthDTO{Date: date})

	if err != nil {
		return nil, err
	}

	categories, err := s.db.GetMonthExpensesByCategory(date)

	if err != nil {
		return nil, err
	}
	result := models.MonthExpensesByCategoryDTO{MonthTotal: *monthTotal, Categories: categories}

	return &result, nil

}
