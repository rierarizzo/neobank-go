package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCustomer(t *testing.T) {
	t.Run("creates customer with all fields", func(t *testing.T) {
		id := uuid.New()
		now := time.Now()
		dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

		c := Customer{
			ID:             id,
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
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		assert.Equal(t, id, c.ID)
		assert.Equal(t, "123456789", c.IdentityNumber)
		assert.Equal(t, "John", c.FirstName)
		assert.Equal(t, "Doe", c.LastName)
		assert.Equal(t, "john@example.com", c.Email)
		assert.Equal(t, "+1234567890", c.PhoneNumber)
		assert.Equal(t, &dob, c.DateOfBirth)
		assert.Equal(t, "US", c.Nationality)
		assert.Equal(t, "123 Main St", c.AddressLine1)
		assert.Equal(t, "Apt 4B", c.AddressLine2)
		assert.Equal(t, "New York", c.City)
		assert.Equal(t, "NY", c.State)
		assert.Equal(t, "10001", c.PostalCode)
		assert.Equal(t, "USA", c.Country)
		assert.Equal(t, now, c.CreatedAt)
		assert.Equal(t, now, c.UpdatedAt)
	})

	t.Run("creates customer with optional fields nil", func(t *testing.T) {
		id := uuid.New()
		now := time.Now()

		c := Customer{
			ID:             id,
			IdentityNumber: "987654321",
			FirstName:      "Jane",
			LastName:       "Smith",
			Email:          "jane@example.com",
			PhoneNumber:    "+0987654321",
			Nationality:    "UK",
			AddressLine1:   "456 High St",
			City:           "London",
			Country:        "UK",
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		assert.Equal(t, id, c.ID)
		assert.Nil(t, c.DateOfBirth)
		assert.Empty(t, c.AddressLine2)
		assert.Empty(t, c.State)
		assert.Empty(t, c.PostalCode)
	})
}

func TestAccount(t *testing.T) {
	t.Run("creates account with all fields", func(t *testing.T) {
		customerID := uuid.New()
		id := uuid.New()
		now := time.Now()

		a := Account{
			ID:              id,
			CustomerID:      customerID,
			LedgerAccountID: 12345,
			AccountType:     "checking",
			Status:          "active",
			OpenedAt:        now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		assert.Equal(t, id, a.ID)
		assert.Equal(t, customerID, a.CustomerID)
		assert.Equal(t, int64(12345), a.LedgerAccountID)
		assert.Equal(t, "checking", a.AccountType)
		assert.Equal(t, "active", a.Status)
		assert.Nil(t, a.ClosedAt)
	})

	t.Run("creates closed account", func(t *testing.T) {
		customerID := uuid.New()
		id := uuid.New()
		now := time.Now()
		closedAt := now.Add(24 * time.Hour)

		a := Account{
			ID:              id,
			CustomerID:      customerID,
			LedgerAccountID: 67890,
			AccountType:     "savings",
			Status:          "closed",
			OpenedAt:        now,
			ClosedAt:        &closedAt,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		assert.Equal(t, "closed", a.Status)
		assert.NotNil(t, a.ClosedAt)
		assert.Equal(t, &closedAt, a.ClosedAt)
	})
}
