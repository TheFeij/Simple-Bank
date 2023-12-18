package requests

type TransferRequest struct {
	FromAccountID uint64 `json:"fromAccountID" validate:"required;nefield=ToAccountID"`
	ToAccountID   uint64 `json:"toAccountID" validate:"required"`
	Amount        uint32 `json:"amount" validate:"required;gt=0"`
}
