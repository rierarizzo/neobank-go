package domain

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	IdentityNumber string     `db:"identity_number" json:"identityNumber"`
	FirstName      string     `db:"first_name" json:"firstName"`
	LastName       string     `db:"last_name" json:"lastName"`
	Email          string     `db:"email" json:"email"`
	PhoneNumber    string     `db:"phone_number" json:"phoneNumber"`
	DateOfBirth    *time.Time `db:"date_of_birth" json:"dateOfBirth,omitempty"`
	Nationality    string     `db:"nationality" json:"nationality"`
	AddressLine1   string     `db:"address_line1" json:"addressLine1"`
	AddressLine2   string     `db:"address_line2" json:"addressLine2,omitempty"`
	City           string     `db:"city" json:"city"`
	State          string     `db:"state" json:"state,omitempty"`
	PostalCode     string     `db:"postal_code" json:"postalCode,omitempty"`
	Country        string     `db:"country" json:"country,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}
