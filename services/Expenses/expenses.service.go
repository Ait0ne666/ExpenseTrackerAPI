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

func (s *ExpensesService) CreateCategory(title string, userID string) (*models.Category, error) {

	category, err := s.db.CreateCategory(title, userID)

	if err != nil {
		return nil, err
	}

	return category, nil

}

func (s *ExpensesService) DeleteCategory(id string, userID string) (*string, error) {

	err := s.db.DeleteCategory(id, userID)

	if err != nil {
		return nil, err
	}

	return &id, nil

}

func (s *ExpensesService) GetCategoryList(query string, userID string) (*[]models.Category, error) {

	list, err := s.db.GetCategoryList(query, userID)

	if err != nil {
		return nil, err
	}

	return list, nil

}

func (s *ExpensesService) DeleteExpense(id string, userID string) (*string, error) {

	err := s.db.DeleteExpense(id, userID)

	if err != nil {
		return nil, err
	}

	return &id, nil

}

func (s *ExpensesService) CreateExpense(dto models.ExpenseDTO, userID string) (*models.Expense, error) {

	expense := models.Expense{}

	if dto.CategoryID == nil {

		category, err := s.db.CreateCategory(dto.CategoryName, userID)

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

	err := s.db.UpsertExpense(&expense, userID)

	if err != nil {
		return nil, err
	}

	return &expense, nil

}

func (s *ExpensesService) GetDayTotalExpenses(date time.Time, currency models.Currency, userID string) (*models.DayTotalResulDTO, error) {

	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	dayExpenses, err := s.db.GetDayExpenses(date, userID)

	if err != nil {
		return nil, err
	}

	monthExpense, err := s.db.GetMonthExpenses(&models.MonthDTO{Date: date}, userID)

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

func (s *ExpensesService) GetMonthExpenses(dto *models.MonthDTO, userID string) (*models.MonthExpensesDTO, error) {
	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	currency := dto.Currency

	expenses, err := s.db.GetMonthExpenses(dto, userID)

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
			expenses[i].Currency = currency
		}

	}

	result := models.MonthExpensesDTO{MonthTotal: monthTotal, Expenses: expenses}

	return &result, nil

}

func (s *ExpensesService) GetMonthExpensesByCategory(date time.Time, currency models.Currency, userID string) (*models.MonthExpensesByCategoryDTO, error) {

	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	categories := make(map[string]models.MonthCategory)

	expenses, err := s.db.GetMonthExpenses(&models.MonthDTO{Date: date}, userID)

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

func (s *ExpensesService) SyncDatabase(syncData models.SyncDTO, userID string) (*models.SyncResultDTO, error) {

	result := make([]models.ExpenseWithCreatedDTO, 0)

	rates, err := s.common.GetCurrencyRate()

	if err != nil {
		return nil, err
	}

	if syncData.LastSync == nil {
		expenses, err := s.db.GetAllExpensesAfter(nil, userID)

		if err != nil {
			return nil, err
		}

		for _, val := range expenses {
			expense := models.ExpenseWithCreatedDTO{}

			expense.Amount = val.Amount
			expense.CategoryID = &val.CategoryID
			expense.CategoryName = val.Category.Title
			expense.CloudID = &val.ID
			expense.Title = val.Title
			expense.CreatedAt = val.CreatedAt
			expense.Currency = val.Currency
			expense.Date = val.Date
			if val.DeletedAt != nil {
				expense.DeletedAt = &val.DeletedAt.Time
			}
			expense.UpdatedAt = val.UpdatedAt

			result = append(result, expense)
		}

		return &models.SyncResultDTO{
			Expenses:       result,
			UpdatedExpense: make([]models.ExpenseWithCreatedDTO, 0),
			Rates:          rates,
		}, nil

	} else {

		expensesLocal := syncData.Expenses

		for i, exp := range expensesLocal {

			if exp.CloudID == nil {

				createdExp, err := s.CreateExpense(models.ExpenseDTO{
					Title:        exp.Title,
					Date:         exp.Date,
					CategoryName: exp.CategoryName,
					Amount:       exp.Amount,
					Currency:     exp.Currency,
					CategoryID:   exp.CategoryID,
				}, userID)

				if err != nil {
					success := false
					expensesLocal[i].Success = &success
				} else {
					success := true
					expensesLocal[i].Success = &success
					expensesLocal[i].CloudID = &createdExp.ID
				}

			} else {
				existingExp, err := s.db.GetExpenseById(userID, *exp.CloudID)

				if err != nil {
					success := false
					expensesLocal[i].Success = &success
				} else {

					if exp.DeletedAt != nil {

						if existingExp.DeletedAt == nil {
							err := s.db.DeleteExpense(existingExp.ID, userID)
							if err != nil {
								success := false
								expensesLocal[i].Success = &success
							} else {
								success := true
								expensesLocal[i].Success = &success
							}
						} else {
							success := false
							expensesLocal[i].Success = &success
						}

					} else {
						lastUpdate := exp.UpdatedAt

						if existingExp.UpdatedAt.After(lastUpdate) {
							success := true
							expensesLocal[i].Success = &success
						} else {
							_, err := s.CreateExpense(models.ExpenseDTO{
								ID:           exp.CloudID,
								Title:        exp.Title,
								Date:         exp.Date,
								CategoryName: exp.CategoryName,
								Amount:       exp.Amount,
								Currency:     exp.Currency,
								CategoryID:   exp.CategoryID,
							}, userID)

							if err != nil {
								success := false
								expensesLocal[i].Success = &success
							} else {
								success := true
								expensesLocal[i].Success = &success
							}
						}

					}
				}

			}

		}

		expensesToSync, err := s.db.GetAllExpensesAfter(syncData.LastSync, userID)

		if err != nil {
			return nil, err
		}

		localExpensesMap := make(map[string]int)

		for _, localExp := range expensesLocal {
			if localExp.CloudID != nil {

				localExpensesMap[*localExp.CloudID] = *localExp.ID
			} else {
				print("no id found")
			}
		}

		for _, exp := range expensesToSync {

			expWithCreatedAt := models.ExpenseWithCreatedDTO{}

			id, ok := localExpensesMap[exp.ID]

			if ok {
				expWithCreatedAt.ID = &id
			}

			expWithCreatedAt.Amount = exp.Amount
			expWithCreatedAt.CategoryID = &exp.CategoryID
			expWithCreatedAt.CategoryName = exp.Category.Title
			expWithCreatedAt.CloudID = &exp.ID
			expWithCreatedAt.CreatedAt = exp.CreatedAt
			if exp.DeletedAt != nil {
				expWithCreatedAt.DeletedAt = &exp.DeletedAt.Time
			}

			expWithCreatedAt.Currency = exp.Currency
			expWithCreatedAt.Date = exp.Date
			expWithCreatedAt.Title = exp.Title
			expWithCreatedAt.UpdatedAt = exp.UpdatedAt

			result = append(result, expWithCreatedAt)

		}

		print(result)

		return &models.SyncResultDTO{
			Expenses:       result,
			UpdatedExpense: expensesLocal,
			Rates:          rates,
		}, nil
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
