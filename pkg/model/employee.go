package model

// Employee model.
type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Job       string `json:"job"`
}
