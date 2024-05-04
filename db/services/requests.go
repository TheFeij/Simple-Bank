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

// ListAccountsRequest represents a request to get a list of a user's accounts
type ListAccountsRequest struct {
	// Owner is the username of the owner of the accounts
	Owner string
	// PageSize represents number of accounts in a page
	PageSize int
	// PageNumber page number
	PageNumber int
}
