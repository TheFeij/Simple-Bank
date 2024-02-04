package responses

import "time"

type TransferResponse struct {
	TransferID   uint64    `json:"transfer_id"`
	SrcAccountID uint64    `json:"src_account_id"`
	DstAccountID uint64    `json:"dst_account_id"`
	Time         time.Time `json:"time"`
	Amount       int64     `json:"amount"`
}
