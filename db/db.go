package db

import (
	"expense_tracker/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {

	name, retrieved_name := os.LookupEnv("db_name")
	user, retrieved_user := os.LookupEnv("db_user")
	host, retrieved_host := os.LookupEnv("db_host")
	pass, retrieved_pass := os.LookupEnv("db_pass")

	dsn := "host=localhost user=postgres password=fibonachi0%0 dbname=gorm port=5432 sslmode=disable"
	if retrieved_name && retrieved_host && retrieved_pass && retrieved_user {
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require password=%s", host, user, name, pass)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Category{}, &models.Expense{}, &models.CurrencyRate{}, &models.User{})

	DB = db
	return db, nil
}
