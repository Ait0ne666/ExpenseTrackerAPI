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

func (p *ExpensesDAO) CreateCategory(title string) (*models.Category, error) {

	existCategory := make([]models.Category, 0)

	if err := p.db.Table("categories").Where("title = ?", title).Find(&existCategory).Error; err != nil {
		return nil, err
	}

	if len(existCategory) != 0 {
		return &existCategory[0], nil
	}

	category := models.Category{}

	category.Title = title

	if err := p.db.Table("categories").Create(&category).Take(&category).Error; err != nil {

		return nil, err

	}

	return &category, nil

}

func (p *ExpensesDAO) DeleteCategory(id string) error {

	if err := p.db.Table("categories").Where("id = ?", id).Delete(&models.Category{}).Error; err != nil {

		return err

	}

	return nil

}

func (p *ExpensesDAO) GetCategoryList(query string) (*[]models.Category, error) {

	categoryList := make([]models.Category, 0)

	if err := p.db.Table("categories").Where("title LIKE ?", "%"+query+"%").Find(&categoryList).Error; err != nil {
		return nil, err
	}

	return &categoryList, nil

}

func (p *ExpensesDAO) UpsertExpense(expense *models.Expense) error {

	if expense.ID == "" {
		if err := p.db.Table("expenses").Preload("Category").Create(expense).Take(expense).Error; err != nil {
			return err
		}
	} else {
		if err := p.db.Table("expenses").Preload("Category").Updates(expense).Take(expense).Error; err != nil {
			return err
		}
	}

	return nil

}

func (p *ExpensesDAO) DeleteExpense(id string) error {

	if err := p.db.Table("expenses").Where("id = ?", id).Delete(&models.Expense{}).Error; err != nil {

		return err

	}

	return nil

}

type TotalDAO struct {
	Total float64 `json:"total" gorm:"total"`
}

func (p *ExpensesDAO) GetDayTotalExpenses(date time.Time) (*float64, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 24).Add(-1 * time.Second)

	total := TotalDAO{}
	total.Total = 0

	if err := p.db.Debug().Table("expenses").Where("date >= ? AND date <= ?", start, end).Select("sum(amount) as total").Take(&total).Error; err != nil {

		return nil, err

	}

	return &total.Total, nil

}

func (p *ExpensesDAO) GetMonthTotalExpenses(dto *models.MonthDTO) (*float64, error) {
	date := dto.Date
	category := dto.CategoryID

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-1 * time.Second)

	total := TotalDAO{}
	total.Total = 0

	if category != nil && *category != "" {
		if err := p.db.Debug().Table("expenses").Where("date >= ? AND date <= ? AND category_id = ?", start, end, category).Select("sum(amount) as total").Take(&total).Error; err != nil {
			return nil, err
		}
	} else {
		if err := p.db.Debug().Table("expenses").Where("date >= ? AND date <= ?", start, end).Select("sum(amount) as total").Take(&total).Error; err != nil {
			return nil, err
		}
	}

	return &total.Total, nil

}

func (p *ExpensesDAO) GetMonthExpenses(dto *models.MonthDTO) ([]models.Expense, error) {
	date := dto.Date
	categoryID := dto.CategoryID

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-1 * time.Second)

	expenses := make([]models.Expense, 0)

	if categoryID != nil && *categoryID != "" {
		if err := p.db.Debug().Table("expenses").Preload("Category").Where("date >= ? AND date <= ? AND category_id = ?", start, end, categoryID).Order("date Desc").Find(&expenses).Error; err != nil {

			return nil, err

		}
	} else {
		if err := p.db.Debug().Table("expenses").Preload("Category").Where("date >= ? AND date <= ?", start, end).Order("date Desc").Find(&expenses).Error; err != nil {

			return nil, err

		}
	}

	return expenses, nil

}

func (p *ExpensesDAO) GetMonthExpensesByCategory(date time.Time) ([]models.MonthCategory, error) {

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-1 * time.Second)

	categories := make([]models.MonthCategory, 0)

	if err := p.db.Debug().Table("categories").Select("categories.id as category_id, max(categories.title) as category_title, sum(e.amount) as amount").Joins("join expenses e on categories.id = e.category_id").Where("date >= ? AND date <= ?", start, end).Group("categories.id").Find(&categories).Error; err != nil {

		return nil, err

	}

	return categories, nil

}
