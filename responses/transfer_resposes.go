package responses

import "time"

type TransferResponse struct {
	TransferID      uint64    `json:"transfer_id"`
	SrcAccountID    uint64    `json:"src_account_id"`
	DstAccountID    uint64    `json:"dst_account_id"`
	IncomingEntryID uint64    `json:"incoming_entry_id"`
	OutgoingEntryID uint64    `json:"out_going_entry_id"`
	CreatedAt       time.Time `json:"created_at"`
	Amount          int64     `json:"amount"`
}
