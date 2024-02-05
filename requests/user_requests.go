package requests

type CreateAccountRequest struct {
	Owner   string `json:"owner" binding:"required;max=50;min=2"`
	Balance uint64 `json:"balance" binding:"gt=0"`
}

type DepositRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
	Amount    uint32 `json:"amount" binding:"required;gt=0"`
}

type WithdrawRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
	Amount    uint32 `json:"amount" binding:"required;gt=0"`
}
