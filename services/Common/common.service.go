package common

import (
	"encoding/json"
	"expense_tracker/models"
	repository "expense_tracker/repositories"
	"io/ioutil"
	"net/http"
	"time"
)

type CommonService struct {
	db repository.CommonDAO
}

func NewCommonService(
	db repository.CommonDAO,
) *CommonService {
	return &CommonService{
		db: db,
	}
}

func (s *CommonService) GetCurrencyRate() (*models.CurrencyRate, error) {

	rate, err := s.db.GetCurrencyRate(time.Now())

	if err != nil {
		return nil, err
	}

	if rate == nil {
		rate, err = getRateFromAPi()

		if err != nil {
			return nil, err
		}

		err = s.db.InsertCurrencyRate(rate)

		if err != nil {
			return nil, err
		}
	}

	return rate, nil
}

type APIRates struct {
	USD float64 `json:"USD"`
	RUB float64 `json:"RUB"`
	THB float64 `json:"THB"`
}

type CurrencyResponse struct {
	Rates APIRates `json:"rates"`
}

func getRateFromAPi() (*models.CurrencyRate, error) {
	url := "http://api.exchangeratesapi.io/v1/latest?access_key=fc71495558512551e12571433cd53512&symbols=USD,RUB,THB"

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	rates := CurrencyResponse{}

	if err := json.Unmarshal(bytes, &rates); err != nil {
		return nil, err
	}

	result := models.CurrencyRate{}

	result.Date = time.Now()
	result.EUR = rates.Rates.RUB
	result.USD = rates.Rates.RUB / rates.Rates.USD
	result.TBH = rates.Rates.RUB / rates.Rates.THB

	return &result, nil
}
