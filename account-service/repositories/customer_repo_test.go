package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rierarizzo/neobank-go/account-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		customerID := uuid.New()
		dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

		mock.ExpectExec("INSERT INTO account.customers").
			WithArgs(
				customerID, "123456789", "John", "Doe", "john@example.com", "+1234567890",
				&dob, "US", "123 Main St", "Apt 4B", "New York", "NY", "10001", "USA",
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		customer := domain.Customer{
			ID:             customerID,
			IdentityNumber: "123456789",
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@example.com",
			PhoneNumber:    "+1234567890",
			DateOfBirth:    &dob,
			Nationality:    "US",
			AddressLine1:   "123 Main St",
			AddressLine2:   "Apt 4B",
			City:           "New York",
			State:          "NY",
			PostalCode:     "10001",
			Country:        "USA",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err := repo.CreateCustomer(ctx, customer)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		customerID := uuid.New()

		mock.ExpectExec("INSERT INTO account.customers").
			WithArgs(
				customerID, "987654321", "Jane", "Smith", "jane@example.com", "+0987654321",
				(*time.Time)(nil), "UK", "456 High St", "", "London", "", "", "UK",
			).
			WillReturnError(sql.ErrConnDone)

		customer := domain.Customer{
			ID:             customerID,
			IdentityNumber: "987654321",
			FirstName:      "Jane",
			LastName:       "Smith",
			Email:          "jane@example.com",
			PhoneNumber:    "+0987654321",
			Nationality:    "UK",
			AddressLine1:   "456 High St",
			City:           "London",
			Country:        "UK",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err := repo.CreateCustomer(ctx, customer)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CreateCustomer failed")
	})
}
