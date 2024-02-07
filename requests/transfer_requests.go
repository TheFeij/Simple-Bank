package requests

type TransferRequest struct {
	FromAccountID uint64 `json:"from_account_id" binding:"required,min=1,nefield=ToAccountID"`
	ToAccountID   uint64 `json:"to_account_id" binding:"required,min=1"`
	Amount        uint32 `json:"amount" binding:"required,gt=0"`
}
