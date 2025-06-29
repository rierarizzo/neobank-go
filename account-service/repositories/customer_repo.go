package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rierarizzo/neobank-go/account-service/domain"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer domain.Customer) error
}

type customerPostgresRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerPostgresRepository{db: db}
}

func (r *customerPostgresRepository) CreateCustomer(ctx context.Context, customer domain.Customer) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO account.customers 
    		(id, identity_number, first_name, last_name, email, phone_number, date_of_birth, nationality, address_line1, 
    		 address_line2, city, state, postal_code, country) 
    		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		customer.ID, customer.IdentityNumber, customer.FirstName, customer.LastName, customer.Email, customer.PhoneNumber,
		customer.DateOfBirth, customer.Nationality, customer.AddressLine1, customer.AddressLine2, customer.City,
		customer.State, customer.PostalCode, customer.Country)
	if err != nil {
		return fmt.Errorf("CreateCustomer failed at insert customer: %w", err)
	}

	return nil
}
