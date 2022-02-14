package message

const (
	EmployeeUpdated EventType = "EmployeeUpdated"
)

type EmployeeUpdatedEventType string

const (
	TypeEmployeeUpdated EmployeeUpdatedEventType = "EmployeeUpdated"
)

// EmployeeUpdatedEvent event captures updates to an employee.
type EmployeeUpdatedEvent struct {
	Metadata Metadata `json:"metadata"`
	Type     EmployeeUpdatedEventType
	Data     struct {
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
	} `json:"data"`
}
