package requests

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required;max=50;min=2"`
	Currency string `json:"currency" binding:"required;min=2;max;5"`
}

type DepositRequest struct {
	AccountID uint64 `json:"accountID" binding:"required"`
	Amount    uint32 `json:"amount" binding:"required;gt=0"`
}

type WithdrawRequest struct {
	AccountID uint64 `json:"accountID" binding:"required"`
	Amount    uint32 `json:"amount" binding:"required;gt=0"`
}
