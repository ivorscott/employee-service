package msg

import "encoding/json"

// UnmarshalEmployeeUpdatedEvent parses the JSON-encoded data and returns EmployeeUpdatedEvent.
func UnmarshalEmployeeUpdatedEvent(data []byte) (EmployeeUpdatedEvent, error) {
	var e EmployeeUpdatedEvent
	err := json.Unmarshal(data, &e)
	return e, err
}

// Marshal JSON encodes EmployeeUpdatedEvent.
func (r *EmployeeUpdatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

const (
	// EmployeeUpdated is a MessageType representing an Employee updated event.
	EmployeeUpdated MessageType = "EmployeeUpdated"
)

// EmployeeUpdatedEventType represents an Employee updated event type.
type EmployeeUpdatedEventType string

const (
	// TypeEmployeeUpdated represents an EmployeeUpdated event.
	TypeEmployeeUpdated EmployeeUpdatedEventType = "EmployeeUpdated"
)

// EmployeeUpdatedEvent event captures updates to an employee.
type EmployeeUpdatedEvent struct {
	Metadata Metadata `json:"metadata"`
	Type     EmployeeUpdatedEventType
	Data     struct {
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
		Salary                *int    `db:"salary"`
		Position              *string `db:"position"`
		EmergencyContactName  *string `db:"emergency_contact_name"`
		EmergencyContactEmail *string `db:"emergency_contact_email"`
		EmergencyContactPhone *string `db:"emergency_contact_phone"`
		Deleted               *bool   `db:"deleted"`
		UpdatedAt             string  `db:"updated_at"`
	} `json:"data"`
}
