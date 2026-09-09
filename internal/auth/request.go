package auth

// RegisterRequest holds the fields required to create a new user account.
type RegisterRequest struct {
	Name     string `json:"name"     example:"Alice Smith"`
	Email    string `json:"email"    example:"alice@example.com"`
	Password string `json:"password" example:"s3cr3tP@ssw0rd"`
}

// RegisterResponse is returned on successful registration.
type RegisterResponse struct {
	Message string `json:"message" example:"user registered successfully"`
}

// LoginRequest holds credentials for authenticating a user.
type LoginRequest struct {
	Email    string `json:"email"    example:"alice@example.com"`
	Password string `json:"password" example:"s3cr3tP@ssw0rd"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	ID    int64  `json:"id"    example:"42"`
	Email string `json:"email" example:"alice@example.com"`
	Name  string `json:"name"  example:"Alice Smith"`
}

// ErrorResponse is the standard error envelope returned on failure.
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// MessageResponse is the standard success envelope for operations that return no data.
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}
