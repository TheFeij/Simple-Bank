package requests

type TransferRequest struct {
	FromAccountID uint64 `json:"from_account_id" validate:"required;nefield=ToAccountID"`
	ToAccountID   uint64 `json:"to_account_id" validate:"required"`
	Amount        uint32 `json:"amount" validate:"required;gt=0"`
}
