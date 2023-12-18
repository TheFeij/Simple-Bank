package requests

type CreateAccountRequest struct {
	Owner    string `json:"owner" validate:"required;max=50;min=2"`
	Currency string `json:"currency" validate:"required;min=2;max;5"`
}

type DepositRequest struct {
	AccountID uint64 `json:"accountID" validate:"required"`
	Amount    uint32 `json:"amount" validate:"required;gt=0"`
}

type WithdrawRequest struct {
	AccountID uint64 `json:"accountID" validate:"required"`
	Amount    uint32 `json:"amount" validate:"required;gt=0"`
}
