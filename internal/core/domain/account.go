package domain

// Struct for POST account
type PostAccount struct {
	ID      string `json:"account_id"`
	Balance string `json:"initial_balance"`
}

// Struct for GET account
type Account struct {
	ID      string `json:"account_id" db:"id"`
	Balance string `json:"balance" db:"balance"`
}
