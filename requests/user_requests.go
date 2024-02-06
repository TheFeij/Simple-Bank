package requests

type CreateAccountRequest struct {
	Owner   string `json:"owner" binding:"required,alpha"`
	Balance uint64 `json:"balance" binding:"gte=0"`
}

type GetAccountRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

type DepositRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
	Amount    uint32 `json:"amount" binding:"required,gt=0"`
}

type WithdrawRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
	Amount    uint32 `json:"amount" binding:"required,gt=0"`
}
