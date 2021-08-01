package structure

type User struct {
	// ID        string    `json:"id" sql:"id"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Bio      string `json:"bio"`
	// CreatedAt time.Time `json:"createdat"`
	// UpdatedAt time.Time `json:"updatedat"`
	// UserID    	string    `json:"userID"`
}

// type newBio struct {
// 	Bio string `json:"bio"`
// }

type BioUpdate struct {
	Bio   string `json:"bio"`
	Email string `json:"email"`
}

// type respBodyBio struct {

// }
