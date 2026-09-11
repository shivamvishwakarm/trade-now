package user

// Profile holds the public-facing user data returned by the /me endpoint.
type Profile struct {
	// TODO: add fields
	Email string `json:"email"`
	Name  string `json:"name"`
	Id    string `json:"id"`
}
