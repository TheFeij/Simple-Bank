package responses

import "time"

type TransferResponse struct {
	TransferID      int64     `json:"transfer_id"`
	SrcAccountID    int64     `json:"src_account_id"`
	DstAccountID    int64     `json:"dst_account_id"`
	IncomingEntryID int64     `json:"incoming_entry_id"`
	OutgoingEntryID int64     `json:"out_going_entry_id"`
	CreatedAt       time.Time `json:"created_at"`
	Amount          int32     `json:"amount"`
}
