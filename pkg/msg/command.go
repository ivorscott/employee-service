package message

const (
	CreateEmployee CommandType = "CreateEmployee"
)

type CreateEmployeeCommandType string

const (
	TypeCreateEmployee CreateEmployeeCommandType = "CreateEmployee"
)

// CreateEmployeeCommand command captures updates to an employee.
type CreateEmployeeCommand struct {
	Metadata Metadata `json:"metadata"`
	Type     CreateEmployeeCommandType
	Data     struct {
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
	} `json:"data"`
}
