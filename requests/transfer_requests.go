package requests

type TransferRequest struct {
	FromAccountID uint64 `json:"from_account_id" binding:"required,nefield=ToAccountID"`
	ToAccountID   uint64 `json:"to_account_id" binding:"required"`
	Amount        uint32 `json:"amount" binding:"required,gt=0"`
}
