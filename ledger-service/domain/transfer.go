package domain

import "time"

type Transfer struct {
	ID            int64     `json:"id" db:"id"`
	ExternalID    string    `json:"external_id" db:"external_id"`
	FromAccountID int64     `json:"from_account_id" db:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id" db:"to_account_id"`
	Currency      Currency  `json:"currency" db:"currency"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
