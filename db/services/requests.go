package services

// UpdateUserRequest represents a request to update user information.
type UpdateUserRequest struct {
	// Username of the user to update.
	Username string
	// Full name of the user (optional).
	Fullname *string
	// New password for the user (optional).
	Password *string
	// New email address for the user (optional).
	Email *string
}
