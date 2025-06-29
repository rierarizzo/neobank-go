package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rierarizzo/neobank-go/ledger-service/domain"
	"go.uber.org/zap"
)

type LedgerSvc interface {
	CreateAccount(ctx context.Context, currency domain.Currency) (*domain.LedgerAccount, error)
	PostTransfer(ctx context.Context, externalID string, fromAccountID, toAccountID, amount int64, currency domain.Currency) error
	GetBalanceCents(ctx context.Context, accountID int64) (int64, error)
}

type ledgerPostgresSvc struct {
	db *sql.DB
}

func NewLedgerPostgresService(db *sql.DB) LedgerSvc {
	return &ledgerPostgresSvc{db: db}
}

func (s *ledgerPostgresSvc) CreateAccount(ctx context.Context, currency domain.Currency) (*domain.LedgerAccount, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateAccount failed at begin transaction: %w", err)
	}
	defer func() {
		if rbe := tx.Rollback(); rbe != nil && !errors.Is(rbe, sql.ErrTxDone) {
			zap.L().Warn("CreateAccount rollback failed", zap.Error(rbe))
		}
	}()

	var newID int64
	var status string
	err = tx.QueryRowContext(ctx,
		"INSERT INTO ledger.ledger_accounts (currency) VALUES ($1) RETURNING id, status", currency,
	).Scan(&newID, &status)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount failed at insert ledger account: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("CreateAccount failed at commit transaction: %w", err)
	}

	zap.L().Debug("Account created successfully", zap.Int64("createdAccountId", newID))

	ledgerAccount := domain.LedgerAccount{
		ID:       newID,
		Currency: currency,
		Status:   domain.AccountStatus(status),
	}

	return &ledgerAccount, nil
}

func (s *ledgerPostgresSvc) PostTransfer(ctx context.Context, externalID string, fromAccountID, toAccountID, amount int64, currency domain.Currency) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("PostTransfer failed at begin transaction: %w", err)
	}
	defer func() {
		if rbe := tx.Rollback(); rbe != nil && !errors.Is(rbe, sql.ErrTxDone) {
			zap.L().Warn("PostTransfer rollback failed", zap.Error(rbe))
		}
	}()

	count := 0
	err = tx.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM ledger.transfers WHERE external_id = $1", externalID,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("PostTransfer failed: transfer already exists")
	}

	var fromCurr, toCurr string
	err = tx.QueryRowContext(ctx,
		"SELECT currency FROM ledger.ledger_accounts WHERE id = $1", fromAccountID,
	).Scan(&fromCurr)
	if err != nil {
		return err
	}
	err = tx.QueryRowContext(ctx,
		"SELECT currency FROM ledger.ledger_accounts WHERE id = $1", toAccountID,
	).Scan(&toCurr)
	if err != nil {
		return err
	}

	if fromCurr != string(currency) || toCurr != string(currency) {
		return fmt.Errorf("PostTransfer failed: currency mismatch, accounts are %s/%s but transfer is %s",
			fromCurr, toCurr, currency)
	}

	var transferID int64
	err = tx.QueryRowContext(ctx,
		"INSERT INTO ledger.transfers (external_id, from_account_id, to_account_id, currency) VALUES ($1, $2, $3, $4) RETURNING id",
		externalID,
		fromAccountID,
		toAccountID,
		currency,
	).Scan(&transferID)
	if err != nil {
		return fmt.Errorf("PostTransfer failed at insert transfer (header): %w", err)
	}

	insertEntryStmnt := "INSERT INTO ledger.entries (transfer_id, ledger_account_id, amount) VALUES ($1, $2, $3)"
	// Inserting credit
	_, err = tx.ExecContext(ctx, insertEntryStmnt, transferID, toAccountID, +amount)
	if err != nil {
		return fmt.Errorf("PostTransfer failed at insert entries: %w", err)
	}

	// Inserting debit
	_, err = tx.ExecContext(ctx, insertEntryStmnt, transferID, fromAccountID, -amount)
	if err != nil {
		return fmt.Errorf("PostTransfer failed at insert entries: %w", err)
	}

	// Verifies the balance
	var balance int64
	err = tx.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount),0) FROM ledger.entries WHERE transfer_id = $1", transferID,
	).Scan(&balance)
	if err != nil {
		return fmt.Errorf("PostTransfer failed at verify balance query: %w", err)
	}

	if balance != 0 {
		return fmt.Errorf("PostTransfer failed: balance is different from zero")
	}

	return tx.Commit()
}

func (s *ledgerPostgresSvc) GetBalanceCents(ctx context.Context, accountID int64) (int64, error) {
	var balance int64
	err := s.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM ledger.entries WHERE ledger_account_id = $1", accountID,
	).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("GetBalanceCents failed at verify balance query: %w", err)
	}

	return balance, nil
}
