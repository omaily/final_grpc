package midleware

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

type Transfer struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type Exchange struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
}
