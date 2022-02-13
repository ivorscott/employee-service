// Package event stores all the events sent or received by the service.
package event

// EmployeeUpdated event captures updates to an employee.
type EmployeeUpdated struct {
	EmailAddress          *string `json:"emailAddress"`
	FirstName             *string `json:"firstName"`
	MiddleName            *string `json:"middleName"`
	LastName              *string `json:"lastName"`
	PhoneNumber           *string `json:"phoneNumber"`
	Picture               *string `json:"picture"`
	Language              *string `json:"language"`
	Country               *string `json:"country"`
	City                  *string `json:"city"`
	Zipcode               *string `json:"zipcode"`
	Salary                *int    `json:"salary"`
	Position              *string `json:"position"`
	EmergencyContactName  *string `json:"emergencyContactName"`
	EmergencyContactEmail *string `json:"emergencyContactEmail"`
	EmergencyContactPhone *string `json:"emergencyContactPhone"`
	Deleted               *bool   `json:"deleted"`
	UpdatedAt             string  `json:"updatedAt"`
}
