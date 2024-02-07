package requests

type TransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1,nefield=ToAccountID"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int32 `json:"amount" binding:"required,gt=0"`
}
