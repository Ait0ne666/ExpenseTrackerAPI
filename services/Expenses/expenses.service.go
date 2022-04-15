package expenses

import (
	"expense_tracker/models"
	repository "expense_tracker/repositories"
	common "expense_tracker/services/Common"
	"time"
)

type ExpensesService struct {
	db     repository.ExpensesDAO
	common common.CommonService
}

func NewExpensesService(
	db repository.ExpensesDAO,
	common *common.CommonService,
) *ExpensesService {
	return &ExpensesService{
		db:     db,
		common: *common,
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

func (s *ExpensesService) GetCategoryList(query string) (*[]models.Category, error) {

	list, err := s.db.GetCategoryList(query)

	if err != nil {
		return nil, err
	}

	return list, nil

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

func (s *ExpensesService) GetDayTotalExpenses(date time.Time, currency models.Currency) (*models.DayTotalResulDTO, error) {

	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	dayExpenses, err := s.db.GetDayExpenses(date)

	if err != nil {
		return nil, err
	}

	monthExpense, err := s.db.GetMonthExpenses(&models.MonthDTO{Date: date})

	if err != nil {
		return nil, err
	}

	var dayTotal float64 = 0
	var monthTotal float64 = 0

	for _, exp := range dayExpenses {

		if exp.Currency == currency {
			dayTotal = dayTotal + exp.Amount
		} else {
			dayTotal = dayTotal + exchangeCurrency(exp.Amount, exp.Currency, currency, rates)
		}

	}

	for _, exp := range monthExpense {

		if exp.Currency == currency {
			monthTotal = monthTotal + exp.Amount
		} else {
			monthTotal = monthTotal + exchangeCurrency(exp.Amount, exp.Currency, currency, rates)
		}

	}

	result := models.DayTotalResulDTO{DayTotal: dayTotal, MonthTotal: monthTotal}

	return &result, nil

}

func (s *ExpensesService) GetMonthExpenses(dto *models.MonthDTO) (*models.MonthExpensesDTO, error) {
	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	currency := dto.Currency

	expenses, err := s.db.GetMonthExpenses(dto)

	if err != nil {
		return nil, err
	}

	var monthTotal float64 = 0

	for i, exp := range expenses {

		if exp.Currency == currency {
			monthTotal = monthTotal + exp.Amount
		} else {
			expenseAmount := exchangeCurrency(exp.Amount, exp.Currency, currency, rates)
			monthTotal = monthTotal + expenseAmount

			expenses[i].Amount = expenseAmount
		}

	}

	result := models.MonthExpensesDTO{MonthTotal: monthTotal, Expenses: expenses}

	return &result, nil

}

func (s *ExpensesService) GetMonthExpensesByCategory(date time.Time, currency models.Currency) (*models.MonthExpensesByCategoryDTO, error) {

	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	categories := make(map[string]models.MonthCategory)

	expenses, err := s.db.GetMonthExpenses(&models.MonthDTO{Date: date})

	if err != nil {
		return nil, err
	}

	var monthTotal float64 = 0

	for _, exp := range expenses {

		val, ok := categories[exp.CategoryID]

		if ok {
			curAmount := exchangeCurrency(exp.Amount, exp.Currency, currency, rates)
			val.Amount = val.Amount + curAmount
			categories[exp.CategoryID] = val
			monthTotal = monthTotal + curAmount
		} else {
			month := models.MonthCategory{}

			month.CategoryID = exp.CategoryID
			month.CategoryTitle = exp.Category.Title
			month.Amount = exchangeCurrency(exp.Amount, exp.Currency, currency, rates)
			monthTotal = monthTotal + month.Amount

			categories[exp.CategoryID] = month
		}

	}

	resultCategories := make([]models.MonthCategory, 0)

	for _, category := range categories {
		resultCategories = append(resultCategories, category)
	}

	result := models.MonthExpensesByCategoryDTO{MonthTotal: monthTotal, Categories: resultCategories}

	return &result, nil

}

func (s *ExpensesService) GetCurrencyRate() (*models.CurrencyRate, error) {
	return s.common.GetCurrencyRate()
}

func exchangeCurrency(amount float64, currencyIn, currencyOut models.Currency, rates *models.CurrencyRate) float64 {

	if currencyIn == currencyOut {
		return amount
	}

	if currencyOut == models.RUB {

		rate := getCurrencyRate(currencyIn, rates)
		return amount * rate
	} else if currencyIn == models.RUB {
		rate := getCurrencyRate(currencyOut, rates)
		return amount / rate

	} else {
		rate := getCurrencyRate(currencyIn, rates)
		secondRate := getCurrencyRate(currencyOut, rates)

		return amount * rate / secondRate
	}

}

func getCurrencyRate(currency models.Currency, rates *models.CurrencyRate) float64 {

	switch currency {
	case models.USD:
		return rates.USD
	case models.EUR:
		return rates.EUR
	case models.TBH:
		return rates.TBH
	default:
		return rates.USD
	}

}
