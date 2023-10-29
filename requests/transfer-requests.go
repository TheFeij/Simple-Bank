package requests

type TransferRequest struct {
	FromAccountID uint64 `json:"fromAccountID" binding:"required;nefield=ToAccountID"`
	ToAccountID   uint64 `json:"toAccountID" binding:"required"`
	Amount        uint32 `json:"amount" binding:"required;gt=0"`
}
