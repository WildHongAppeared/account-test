package domain

// Struct for POST transaction
type Transaction struct {
	SourceID      string `json:"source_account_id"`
	DestinationID string `json:"destination_account_id"`
	Amount        string `json:"amount"`
}
