package midleware

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

type Deposit struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}
