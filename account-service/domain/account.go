package domain

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	CustomerID      uuid.UUID  `db:"customer_id" json:"customerId"`
	LedgerAccountID int64      `db:"ledger_account_id" json:"ledgerAccountId"`
	AccountType     string     `db:"account_type" json:"accountType"`
	Status          string     `db:"status" json:"status"`
	OpenedAt        time.Time  `db:"opened_at" json:"openedAt"`
	ClosedAt        *time.Time `db:"closed_at" json:"closedAt,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
}
