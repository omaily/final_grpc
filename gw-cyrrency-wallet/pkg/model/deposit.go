package model

import "fmt"

type Deposit struct {
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

func (d Deposit) Validate() error {
	if !enumCurrency[d.Currency] {
		return fmt.Errorf("invalid amount or currency")
	}
	return nil
}
