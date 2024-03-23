package requests

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type GetAccountsListRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int8  `form:"page_size" binding:"required,min=5,max=10"`
}

type DepositRequest struct {
	AccountID int64 `json:"account_id" binding:"required"`
	Amount    int32 `json:"amount" binding:"required,gt=0"`
}

type WithdrawRequest struct {
	AccountID int64 `json:"account_id" binding:"required"`
	Amount    int32 `json:"amount" binding:"required,gt=0"`
}
