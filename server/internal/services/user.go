// Package services ...
package services

// User ...
type User struct {
	Name string `json:"name"`
}

// UserService ...
type UserService interface {
	GetUser() (User, error)
}

// UserService ...
type userService struct{}

func NewUserService() *userService { //nolint:revive
	return &userService{}
}

// GetUser ...
func (s userService) GetUser() (User, error) {
	return User{
		Name: "John Doe",
	}, nil
}
