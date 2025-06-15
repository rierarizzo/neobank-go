package domain

type Entry struct {
	ID              int64 `json:"id" db:"id"`
	TransferID      int64 `json:"transfer_id" db:"transfer_id"`
	LedgerAccountID int64 `json:"ledger_account_id" db:"ledger_account_id"`
	Amount          int   `json:"amount" db:"amount"`
}
