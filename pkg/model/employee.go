package model

import "time"

// Employee represents an employee.
type Employee struct {
	ID                    string    `db:"employee_id" json:"id"`
	Auth0ID               string    `db:"auth0_id" json:"auth0Id"`
	EmailAddress          string    `db:"email_address" json:"emailAddress"`
	FirstName             string    `db:"first_name" json:"firstName"`
	MiddleName            *string   `db:"middle_name" json:"middleName"`
	LastName              string    `db:"last_name" json:"lastName"`
	PhoneNumber           string    `db:"phone_number" json:"phoneNumber"`
	BirthDate             string    `db:"birth_date" json:"birthDate"`
	HireDate              string    `db:"hire_date" json:"hireDate"`
	Picture               *string   `db:"picture" json:"picture"`
	Language              string    `db:"language" json:"language"`
	Country               string    `db:"country" json:"country"`
	City                  string    `db:"city" json:"city"`
	Zipcode               string    `db:"zipcode" json:"zipcode"`
	Salary                string    `db:"salary" json:"salary"`
	Position              string    `db:"position" json:"position"`
	EmergencyContactName  *string   `db:"emergency_contact_name" json:"emergencyContactName"`
	EmergencyContactEmail *string   `db:"emergency_contact_email" json:"emergencyContactEmail"`
	EmergencyContactPhone *string   `db:"emergency_contact_phone" json:"emergencyContactPhone"`
	Deleted               bool      `db:"deleted" json:"deleted"`
	UpdatedAt             time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt             time.Time `db:"created_at" json:"createdAt"`
}
