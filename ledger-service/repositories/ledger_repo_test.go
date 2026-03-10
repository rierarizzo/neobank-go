package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rierarizzo/neobank-go/ledger-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewLedgerPostgresService(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO ledger.ledger_accounts").
			WithArgs("USD").
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).
				AddRow(1, "active"))
		mock.ExpectCommit()

		account, err := repo.CreateAccount(ctx, "USD")

		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, int64(1), account.ID)
		assert.Equal(t, domain.Currency("USD"), account.Currency)
		assert.Equal(t, domain.AccountStatus("active"), account.Status)
	})

	t.Run("begin transaction error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		account, err := repo.CreateAccount(ctx, "USD")

		assert.Error(t, err)
		assert.Nil(t, account)
	})

	t.Run("insert error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO ledger.ledger_accounts").
			WithArgs("USD").
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		account, err := repo.CreateAccount(ctx, "USD")

		assert.Error(t, err)
		assert.Nil(t, account)
	})
}

func TestPostTransfer(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewLedgerPostgresService(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT COUNT\\(1\\) FROM ledger.transfers").
			WithArgs("ext-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery("SELECT currency FROM ledger.ledger_accounts").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("USD"))
		mock.ExpectQuery("SELECT currency FROM ledger.ledger_accounts").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("USD"))
		mock.ExpectQuery("INSERT INTO ledger.transfers").
			WithArgs("ext-1", 1, 2, "USD").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO ledger.entries").
			WithArgs(1, 2, int64(100)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO ledger.entries").
			WithArgs(1, 1, int64(-100)).
			WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\),0\\) FROM ledger.entries").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(0))
		mock.ExpectCommit()

		err := repo.PostTransfer(ctx, "ext-1", 1, 2, 100, "USD")

		assert.NoError(t, err)
	})

	t.Run("duplicate external ID", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT COUNT\\(1\\) FROM ledger.transfers").
			WithArgs("ext-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectRollback()

		err := repo.PostTransfer(ctx, "ext-1", 1, 2, 100, "USD")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("currency mismatch", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT COUNT\\(1\\) FROM ledger.transfers").
			WithArgs("ext-2").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery("SELECT currency FROM ledger.ledger_accounts").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("USD"))
		mock.ExpectQuery("SELECT currency FROM ledger.ledger_accounts").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("EUR"))
		mock.ExpectRollback()

		err := repo.PostTransfer(ctx, "ext-2", 1, 2, 100, "USD")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "currency mismatch")
	})

	t.Run("balance verification fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT COUNT\\(1\\) FROM ledger.transfers").
			WithArgs("ext-3").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery("SELECT currency FROM ledger.ledger_accounts").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("USD"))
		mock.ExpectQuery("SELECT currency FROM ledger.ledger_accounts").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("USD"))
		mock.ExpectQuery("INSERT INTO ledger.transfers").
			WithArgs("ext-3", 1, 2, "USD").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO ledger.entries").
			WithArgs(1, 2, int64(100)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO ledger.entries").
			WithArgs(1, 1, int64(-100)).
			WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\),0\\) FROM ledger.entries").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(50))
		mock.ExpectRollback()

		err := repo.PostTransfer(ctx, "ext-3", 1, 2, 100, "USD")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "balance is different from zero")
	})
}

func TestGetBalanceCents(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewLedgerPostgresService(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\), 0\\) FROM ledger.entries").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(500))

		balance, err := repo.GetBalanceCents(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(500), balance)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\), 0\\) FROM ledger.entries").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		balance, err := repo.GetBalanceCents(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, int64(0), balance)
	})
}
