package domain

type LedgerAccount struct {
	ID       int64         `json:"id" db:"id"`
	Currency Currency      `json:"currency" db:"currency"`
	Status   AccountStatus `json:"status" db:"status"`
}
