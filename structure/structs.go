package structure

type User struct {
	// ID        string    `json:"id" sql:"id"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	// CreatedAt time.Time `json:"createdat"`
	// UpdatedAt time.Time `json:"updatedat"`
	// UserID    	string    `json:"userID"`
}
