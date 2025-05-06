package accounts

import "time"

type User struct {
	ID string `json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Phone *string `json:"phone"`

	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`

	Role string `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
