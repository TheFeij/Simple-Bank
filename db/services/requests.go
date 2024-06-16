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

// TransferRequest represents a request to transfer money from a source account to another account
type TransferRequest struct {
	// Owner is the username of the owner of the account with id = FromAccountID
	Owner string
	// FromAccountID is the id of the source account
	FromAccountID int64
	// ToAccountID is the id of the destination account
	ToAccountID int64
	// Amount is the amount of money to be transferred from FromAccountID to ToAccountID
	Amount int32
}

// DepositRequest represents a request to deposit money
type DepositRequest struct {
	// Owner of the account
	Owner string
	// AccountID is the id of the account
	AccountID int64
	// Amount is the amount of money to be deposited
	Amount int32
}

// WithdrawRequest represents a request to withdraw money
type WithdrawRequest struct {
	// Owner of the account
	Owner string
	// AccountID is the id of the account
	AccountID int64
	// Amount is the amount of money to be withdrawn
	Amount int32
}
