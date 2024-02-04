package requests

type CreateAccountRequest struct {
	Owner   string `json:"owner" validate:"required;max=50;min=2"`
	Balance uint64 `json:"balance" validate:"gt=0"`
}

type DepositRequest struct {
	AccountID uint64 `json:"account_id" validate:"required"`
	Amount    uint32 `json:"amount" validate:"required;gt=0"`
}

type WithdrawRequest struct {
	AccountID uint64 `json:"account_id" validate:"required"`
	Amount    uint32 `json:"amount" validate:"required;gt=0"`
}
