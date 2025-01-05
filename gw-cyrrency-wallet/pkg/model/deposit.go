package model

import "fmt"

type Transfer struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

var (
	enumCurrency = map[string]bool{
		"USD": true,
		"RUB": true,
		"EUR": true,
	}
)

func (d Transfer) Validate() error {
	if !enumCurrency[d.Currency] || d.Amount <= 0 {
		return fmt.Errorf("invalid amount or currency")
	}
	return nil
}
