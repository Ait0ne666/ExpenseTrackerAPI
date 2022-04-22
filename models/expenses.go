package models

import (
	"time"

	"gorm.io/gorm"
)

type Currency string

const (
	EUR Currency = "eur"
	USD Currency = "usd"
	TBH Currency = "tbh"
	RUB Currency = "rub"
)

type Category struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at" swaggertype:"string" format:"date-time"`
	Title     string         `json:"title" gorm:"title;"`
	UserID    string         `json:"user_id" gorm:"user_id;default:be9e19d1-7150-4a6d-85f0-1cee1c3ac7bc"`
}

type Expense struct {
	ID         string          `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt  time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at" swaggertype:"string" format:"date-time"`
	Title      *string         `json:"title" gorm:"title"`
	Date       time.Time       `json:"date" gorm:"date"`
	Category   Category        `json:"category"`
	Amount     float64         `json:"amount" gorm:"amount"`
	CategoryID string          `json:"category_id" gorm:"category_id"`
	Currency   Currency        `json:"currency" gorm:"currency;default:tbh"`
	UserID     string          `json:"user_id" gorm:"user_id;default:be9e19d1-7150-4a6d-85f0-1cee1c3ac7bc"`
}

type ExpenseDTO struct {
	ID           *string   `json:"id"`
	Title        *string   `json:"title"`
	Date         time.Time `json:"date"`
	CategoryName string    `json:"category_name"`
	CategoryID   *string   `json:"category_id"`
	Amount       float64   `json:"amount"`
	Currency     Currency  `json:"currency"`
}

type ExpenseWithCreatedDTO struct {
	ID           *string    `json:"id"`
	CloudID      *string    `json:"cloud_id"`
	Title        *string    `json:"title"`
	Date         time.Time  `json:"date"`
	CategoryName string     `json:"category_name"`
	CategoryID   *string    `json:"category_id"`
	Amount       float64    `json:"amount"`
	Currency     Currency   `json:"currency"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	Success      *bool      `json:"success"`
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
	Currency   Currency        `json:"currency"`
	Categories []MonthCategory `json:"categories"`
}

type MonthDTO struct {
	Date       time.Time `json:"date"`
	CategoryID *string   `json:"category_id"`
	Currency   Currency  `json:"currency"`
}

type CurrencyRate struct {
	ID   string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Date time.Time `json:"date" gorm:"date"`
	EUR  float64   `json:"eur" gorm:"eur"`
	USD  float64   `json:"usd" gorm:"usd"`
	TBH  float64   `json:"tbh" gorm:"tbh"`
}

type User struct {
	ID       string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Active   string `json:"active" gorm:"active;default:true"`
	Login    string `json:"login" gorm:"login"`
	Password string `json:"password" gorm:"password"`
}

type LoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SyncDTO struct {
	Expenses []ExpenseWithCreatedDTO `json:"expenses"`
	LastSync *time.Time              `json:"last_sync"`
}

type SyncResultDTO struct {
	Expenses       []ExpenseWithCreatedDTO `json:"expenses"`
	UpdatedExpense []ExpenseWithCreatedDTO `json:"updated_expenses"`
}
