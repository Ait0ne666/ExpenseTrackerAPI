package models

import "time"

type Category struct {
	ID    string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title string `json:"title" gorm:"title;uniqueIndex"`
}

type Expense struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title      *string   `json:"title" gorm:"title"`
	Date       time.Time `json:"date" gorm:"date"`
	Category   Category  `json:"category"`
	Amount     float64   `json:"amount" gorm:"amount"`
	CategoryID string    `json:"category_id" gorm:"category_id"`
}

type ExpenseDTO struct {
	ID           *string   `json:"id"`
	Title        *string   `json:"title"`
	Date         time.Time `json:"date"`
	CategoryName string    `json:"category_name"`
	CategoryID   *string   `json:"category_id"`
	Amount       float64   `json:"amount"`
}

type DayTotalResulDTO struct {
	DayTotal   float64 `json:"day_total"`
	MonthTotal float64 `json:"month_total"`
}

type MonthExpensesDTO struct {
	MonthTotal float64   `json:"month_total"`
	Expenses   []Expense `json:"expenses"`
}

type MonthCategory struct {
	CategoryID    string  `json:"category_id" gorm:"category_id"`
	CategoryTitle string  `json:"category_title" gorm:"category_title"`
	Amount        float64 `json:"amount" gorm:"amount"`
}

type MonthExpensesByCategoryDTO struct {
	MonthTotal float64         `json:"month_total"`
	Categories []MonthCategory `json:"categories"`
}

type MonthDTO struct {
	Date       time.Time `json:"date"`
	CategoryID *string   `json:"category_id"`
}
