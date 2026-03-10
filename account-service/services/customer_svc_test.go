package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rierarizzo/neobank-go/account-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) CreateCustomer(ctx context.Context, customer domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func TestCreateCustomer(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	svc := NewCustomerService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		customer := domain.Customer{
			ID:             uuid.New(),
			IdentityNumber: "123456789",
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@example.com",
			PhoneNumber:    "+1234567890",
			Nationality:    "US",
			AddressLine1:   "123 Main St",
			City:           "New York",
			Country:        "USA",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.On("CreateCustomer", mock.Anything, customer).Return(nil).Once()

		err := svc.CreateCustomer(ctx, customer)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		customer := domain.Customer{
			ID:             uuid.New(),
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

		mockRepo.On("CreateCustomer", mock.Anything, customer).Return(errors.New("database error")).Once()

		err := svc.CreateCustomer(ctx, customer)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
