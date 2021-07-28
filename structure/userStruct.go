package structure

import "time"

type User struct {
	// ID        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	Username  string    `json:"username" validate:"min=3`
	TokenHash string    `json:"tokenhash"`
	RefToken  string    `json:"reftoken"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
