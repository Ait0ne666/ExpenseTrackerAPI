package common

import (
	"expense_tracker/models"
	"time"

	"gorm.io/gorm"
)

type CommonDAO struct {
	db *gorm.DB
}

func NewCommonDAO(db *gorm.DB) *CommonDAO {
	return &CommonDAO{db: db}
}

func (r *CommonDAO) InsertCurrencyRate(currency *models.CurrencyRate) error {

	if err := r.db.Debug().Table("currency_rates").Create(currency).Take(currency).Error; err != nil {
		return err
	}

	return nil
}

func (r *CommonDAO) GetCurrencyRate(date time.Time) (*models.CurrencyRate, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	rate := make([]models.CurrencyRate, 0)

	if err := r.db.Debug().Table("currency_rates").Where("date >= ? and date <= ?", start, end).Find(&rate).Error; err != nil {
		return nil, err
	}

	if len(rate) == 0 {
		return nil, nil
	}

	return &rate[0], nil

}
