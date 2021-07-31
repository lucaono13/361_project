package structure

import (
	"time"
)

type User struct {
	// ID        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
	// UserID    	string    `json:"userID"`
}
