package model

import (
	"github.com/go-playground/validator/v10"
)

var employeeValidator *validator.Validate

func init() {
	v := NewValidator()
	employeeValidator = v
}

// Employee represents an employee.
type Employee struct {
	ID                    string  `db:"employee_id" json:"id" validate:"required"`
	Auth0ID               string  `db:"auth0_id" json:"auth0Id"`
	EmailAddress          string  `db:"email_address" json:"emailAddress" validate:"required"`
	FirstName             string  `db:"first_name" json:"firstName" validate:"required"`
	MiddleName            *string `db:"middle_name" json:"middleName" validate:"required"`
	LastName              string  `db:"last_name" json:"lastName" validate:"required"`
	PhoneNumber           string  `db:"phone_number" json:"phoneNumber" validate:"required"`
	BirthDate             string  `db:"birth_date" json:"birthDate" validate:"required"`
	HireDate              string  `db:"hire_date" json:"hireDate" validate:"required"`
	Picture               *string `db:"picture" json:"picture"`
	Language              string  `db:"language" json:"language" validate:"required"`
	Country               string  `db:"country" json:"country" validate:"required"`
	City                  string  `db:"city" json:"city" validate:"required"`
	Zipcode               string  `db:"zipcode" json:"zipcode" validate:"required"`
	Salary                int     `db:"salary" json:"salary" validate:"gte=0,lte=500000"`
	Position              string  `db:"position" json:"position" validate:"required"`
	EmergencyContactName  *string `db:"emergency_contact_name" json:"emergencyContactName"`
	EmergencyContactEmail *string `db:"emergency_contact_email" json:"emergencyContactEmail"`
	EmergencyContactPhone *string `db:"emergency_contact_phone" json:"emergencyContactPhone"`
	Deleted               bool    `db:"deleted" json:"deleted"`
	UpdatedAt             string  `db:"updated_at" json:"updatedAt" validate:"required"`
	CreatedAt             string  `db:"created_at" json:"createdAt" validate:"required"`
}

// Validate employee model.
func (e *Employee) Validate() error {
	return employeeValidator.Struct(e)
}

// NewEmployee represents a new employee.
type NewEmployee struct {
	EmailAddress string  `db:"email_address" json:"emailAddress" validate:"required"`
	FirstName    string  `db:"first_name" json:"firstName" validate:"required"`
	MiddleName   *string `db:"middle_name" json:"middleName" validate:"required"`
	LastName     string  `db:"last_name" json:"lastName" validate:"required"`
	PhoneNumber  string  `db:"phone_number" json:"phoneNumber" validate:"required"`
	BirthDate    string  `db:"birth_date" json:"birthDate" validate:"required"`
	HireDate     string  `db:"hire_date" json:"hireDate" validate:"required"`
	Language     string  `db:"language" json:"language" validate:"required"`
	Country      string  `db:"country" json:"country" validate:"required"`
	City         string  `db:"city" json:"city" validate:"required"`
	Zipcode      string  `db:"zipcode" json:"zipcode" validate:"required"`
	Salary       int     `db:"salary" json:"salary" validate:"required"`
	Position     string  `db:"position" json:"position" validate:"required"`
}

// Validate new employee request.
func (e *NewEmployee) Validate() error {
	return employeeValidator.Struct(e)
}

// UpdateEmployee represents an employee update.
type UpdateEmployee struct {
	EmailAddress          *string `db:"email_address"`
	FirstName             *string `db:"first_name"`
	MiddleName            *string `db:"middle_name"`
	LastName              *string `db:"last_name"`
	PhoneNumber           *string `db:"phone_number"`
	BirthDate             *string `db:"birth_date"`
	HireDate              *string `db:"hire_date"`
	Picture               *string `db:"picture"`
	Language              *string `db:"language"`
	Country               *string `db:"country"`
	City                  *string `db:"city"`
	Zipcode               *string `db:"zipcode"`
	Salary                *int    `db:"salary" validate:"gte=0,lte=500000"`
	Position              *string `db:"position"`
	EmergencyContactName  *string `db:"emergency_contact_name"`
	EmergencyContactEmail *string `db:"emergency_contact_email"`
	EmergencyContactPhone *string `db:"emergency_contact_phone"`
	Deleted               *bool   `db:"deleted"`
	UpdatedAt             string  `db:"updated_at" validate:"required"`
}

// Validate employee update request.
func (e *UpdateEmployee) Validate() error {
	return employeeValidator.Struct(e)
}
