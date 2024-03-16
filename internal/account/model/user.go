package model

// Account represents a user in the system.
type Account struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
