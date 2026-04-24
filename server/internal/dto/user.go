package dto

// CreateUser ...
type CreateUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"-"`
}
