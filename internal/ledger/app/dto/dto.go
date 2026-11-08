package dto

import "time"

type CreateExpenseRequest struct {
	Title        string          `json:"title"`
	TotalAmount  int64           `json:"total_amount"`
	SplitType    string          `json:"split_type"`
	OccurredAt   time.Time       `json:"occurred_at"`
	Payers       []PayerRequest  `json:"payers"`
	Shares       []SharesRequest `json:"shares,omitempty"`
	Participants []string        `json:"participants,omitempty"`
}

type PayerRequest struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}

type SharesRequest struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}
