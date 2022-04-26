package expenses

import (
	"expense_tracker/models"
	"time"

	"gorm.io/gorm"
)

type ExpensesDAO struct {
	db *gorm.DB
}

func NewExpensesDAO(db *gorm.DB) *ExpensesDAO {
	return &ExpensesDAO{db: db}
}

func (p *ExpensesDAO) CreateCategory(title, userID string) (*models.Category, error) {

	existCategory := make([]models.Category, 0)

	if err := p.db.Table("categories").Where("title ILIKE ? and user_id = ?", title, userID).Find(&existCategory).Error; err != nil {
		return nil, err
	}

	if len(existCategory) != 0 {
		return &existCategory[0], nil
	}

	category := models.Category{}

	category.Title = title
	category.UserID = userID

	if err := p.db.Table("categories").Create(&category).Take(&category).Error; err != nil {

		return nil, err

	}

	return &category, nil

}

func (p *ExpensesDAO) DeleteCategory(id, userID string) error {

	if err := p.db.Table("categories").Where("id = ? AND user_id = ?", id, userID).Delete(&models.Category{}).Error; err != nil {

		return err

	}

	return nil

}

func (p *ExpensesDAO) GetCategoryList(query, userID string) (*[]models.Category, error) {

	categoryList := make([]models.Category, 0)

	if err := p.db.Table("categories").Where("title LIKE ? and user_id = ?", "%"+query+"%", userID).Find(&categoryList).Error; err != nil {
		return nil, err
	}

	return &categoryList, nil

}

func (p *ExpensesDAO) UpsertExpense(expense *models.Expense, userID string) error {

	if expense.ID == "" {
		expense.UserID = userID
		if err := p.db.Table("expenses").Preload("Category").Create(expense).Take(expense).Error; err != nil {
			return err
		}
	} else {
		if err := p.db.Table("expenses").Where("user_id = ?", userID).Preload("Category").Updates(expense).Take(expense).Error; err != nil {
			return err
		}
	}

	return nil

}

func (p *ExpensesDAO) DeleteExpense(id, userID string) error {

	if err := p.db.Table("expenses").Where("id = ? and user_id = ?", id, userID).Delete(&models.Expense{}).Error; err != nil {

		return err

	}

	return nil

}

type TotalDAO struct {
	Total float64 `json:"total" gorm:"total"`
}

func (p *ExpensesDAO) GetDayTotalExpenses(date time.Time, userID string) (*float64, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 24).Add(-1 * time.Second)

	total := TotalDAO{}
	total.Total = 0

	if err := p.db.Debug().Table("expenses").Where("date >= ? AND date <= ? and user_id = ?", start, end, userID).Select("sum(amount) as total").Take(&total).Error; err != nil {

		return nil, err

	}

	return &total.Total, nil

}

func (p *ExpensesDAO) GetDayExpenses(date time.Time, userID string) ([]models.Expense, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	expenses := make([]models.Expense, 0)

	if err := p.db.Debug().Table("expenses").Preload("Category").Where("date >= ? AND date <= ? and user_id = ?", start, end, userID).Order("date Desc").Find(&expenses).Error; err != nil {

		return nil, err

	}

	return expenses, nil

}

func (p *ExpensesDAO) GetMonthTotalExpenses(dto *models.MonthDTO, userID string) (*float64, error) {
	date := dto.Date
	category := dto.CategoryID

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-1 * time.Second)

	total := TotalDAO{}
	total.Total = 0

	if category != nil && *category != "" {
		if err := p.db.Debug().Table("expenses").Where("date >= ? AND date <= ? AND category_id = ? and user_id = ?", start, end, category, userID).Select("sum(amount) as total").Take(&total).Error; err != nil {
			return nil, err
		}
	} else {
		if err := p.db.Debug().Table("expenses").Where("date >= ? AND date <= ? and user_id = ?", start, end, userID).Select("sum(amount) as total").Take(&total).Error; err != nil {
			return nil, err
		}
	}

	return &total.Total, nil

}

func (p *ExpensesDAO) GetMonthExpenses(dto *models.MonthDTO, userID string) ([]models.Expense, error) {
	date := dto.Date
	categoryID := dto.CategoryID

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-1 * time.Second)

	expenses := make([]models.Expense, 0)

	if categoryID != nil && *categoryID != "" {
		if err := p.db.Debug().Table("expenses").Preload("Category").Where("date >= ? AND date <= ? AND category_id = ? and user_id = ?", start, end, categoryID, userID).Order("date Desc").Find(&expenses).Error; err != nil {

			return nil, err

		}
	} else {
		if err := p.db.Debug().Table("expenses").Preload("Category").Where("date >= ? AND date <= ? and user_id = ?", start, end, userID).Order("date Desc").Find(&expenses).Error; err != nil {

			return nil, err

		}
	}

	return expenses, nil

}

func (p *ExpensesDAO) GetMonthExpensesByCategory(date time.Time, userID string) ([]models.MonthCategory, error) {

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-1 * time.Second)

	categories := make([]models.MonthCategory, 0)

	if err := p.db.Debug().Table("categories").Select("categories.id as category_id, max(categories.title) as category_title, sum(e.amount) as amount").Joins("join expenses e on categories.id = e.category_id").Where("date >= ? AND date <= ? and user_id = ?", start, end, userID).Group("categories.id").Find(&categories).Error; err != nil {

		return nil, err

	}

	return categories, nil

}

func (p *ExpensesDAO) GetAllExpensesAfter(date *time.Time, userID string) ([]models.Expense, error) {

	expenses := make([]models.Expense, 0)

	args := []interface{}{userID}

	query := "user_id = ? "

	if date != nil {
		args = []interface{}{userID, date, date, date}
		query = query + " AND (updated_at>= ? OR deleted_at >= ? OR created_at >= ?)"
	}

	if err := p.db.Debug().Unscoped().Table("expenses").Preload("Category").Where(query, args...).Find(&expenses).Error; err != nil {
		return expenses, err
	}

	return expenses, nil

}

func (p *ExpensesDAO) GetExpenseById(userID, expenseID string) (*models.Expense, error) {

	expense := models.Expense{}

	if err := p.db.Debug().Unscoped().Table("expenses").Where("user_id = ? AND id = ?", userID, expenseID).Take(&expense).Error; err != nil {
		return nil, err
	}

	return &expense, nil

}
