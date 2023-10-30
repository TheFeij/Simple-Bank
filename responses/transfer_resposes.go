package responses

import "time"

type TransferResponse struct {
	TransferID   uint64    `json:"transferID"`
	SrcAccountID uint64    `json:"srcAccountID"`
	DstAccountID uint64    `json:"dstAccountID"`
	Time         time.Time `json:"time"`
	Currency     string    `json:"currency"`
	Amount       int64     `json:"amount"`
}
