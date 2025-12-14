package models

import (
	"time"
)

type Vote struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ProposalID int64     `json:"proposal_id"`
	VoteType   string    `json:"vote_type"` // "ALLOWLIST" or "BLOCKLIST"
	CreatedAt  time.Time `json:"created_at"`
}

type VoteType string

const (
	VoteTypeAllowlist VoteType = "ALLOWLIST"
	VoteTypeBlocklist VoteType = "BLOCKLIST"
)
