package services

import (
	"context"
	"github.com/rierarizzo/neobank-go/account-service/domain"
	"github.com/rierarizzo/neobank-go/account-service/repositories"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer domain.Customer) error
}

type customerDefaultService struct {
	customerRepo repositories.CustomerRepository
}

func NewCustomerService(customerRepo repositories.CustomerRepository) CustomerService {
	return &customerDefaultService{customerRepo: customerRepo}
}

func (s *customerDefaultService) CreateCustomer(ctx context.Context, customer domain.Customer) error {
	return s.customerRepo.CreateCustomer(ctx, customer)
}
