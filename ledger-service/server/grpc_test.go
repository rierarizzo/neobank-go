package server

import (
	"context"
	"errors"
	"testing"

	"github.com/rierarizzo/neobank-go/ledger-service/domain"
	pb "github.com/rierarizzo/neobank-go/ledger-service/proto/ledger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLedgerRepository struct {
	mock.Mock
}

func (m *MockLedgerRepository) CreateAccount(ctx context.Context, currency domain.Currency) (*domain.LedgerAccount, error) {
	args := m.Called(mock.Anything, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LedgerAccount), args.Error(1)
}

func (m *MockLedgerRepository) PostTransfer(ctx context.Context, externalID string, fromAccountID, toAccountID, amount int64, currency domain.Currency) error {
	args := m.Called(mock.Anything, externalID, fromAccountID, toAccountID, amount, currency)
	return args.Error(0)
}

func (m *MockLedgerRepository) GetBalanceCents(ctx context.Context, accountID int64) (int64, error) {
	args := m.Called(mock.Anything, accountID)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockLedgerRepository)
	svc := NewGRPCServer(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CreateAccount", mock.Anything, mock.Anything).Return(
			&domain.LedgerAccount{ID: 1, Currency: "USD", Status: "active"}, nil).Once()

		resp, err := svc.CreateAccount(ctx, &pb.CreateAccountRequest{Currency: "USD"})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.Id)
		assert.Equal(t, "USD", resp.Currency)
		assert.Equal(t, "active", resp.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("CreateAccount", mock.Anything, mock.Anything).Return(
			nil, errors.New("database error")).Once()

		resp, err := svc.CreateAccount(ctx, &pb.CreateAccountRequest{Currency: "USD"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}

func TestPostTransfer(t *testing.T) {
	mockRepo := new(MockLedgerRepository)
	svc := NewGRPCServer(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("PostTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := svc.PostTransfer(ctx, &pb.PostTransferRequest{
			ExternalID:    "ext-1",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        100,
			Currency:      "USD",
		})

		assert.NoError(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("PostTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("insufficient funds")).Once()

		resp, err := svc.PostTransfer(ctx, &pb.PostTransferRequest{
			ExternalID:    "ext-1",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        100,
			Currency:      "USD",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBalanceCents(t *testing.T) {
	mockRepo := new(MockLedgerRepository)
	svc := NewGRPCServer(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetBalanceCents", mock.Anything, mock.Anything).Return(int64(500), nil).Once()

		resp, err := svc.GetBalanceCents(ctx, &pb.GetBalanceCentsRequest{AccountID: 1})

		assert.NoError(t, err)
		assert.Equal(t, int64(500), resp.BalanceCents)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetBalanceCents", mock.Anything, mock.Anything).Return(int64(0), errors.New("account not found")).Once()

		resp, err := svc.GetBalanceCents(ctx, &pb.GetBalanceCentsRequest{AccountID: 1})

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}
