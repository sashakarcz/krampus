package models

import (
	"time"
)

type Proposal struct {
	ID             int64      `json:"id"`
	Identifier     string     `json:"identifier"`
	RuleType       string     `json:"rule_type"` // "BINARY", "CERTIFICATE", "SIGNINGID", "TEAMID", "CDHASH"
	ProposedPolicy string     `json:"proposed_policy"` // "ALLOWLIST" or "BLOCKLIST"
	CustomMessage  *string    `json:"custom_message,omitempty"`
	CreatedBy      int64      `json:"created_by"`
	Status         string     `json:"status"` // "PENDING", "APPROVED", "REJECTED"
	AllowlistVotes int        `json:"allowlist_votes"`
	BlocklistVotes int        `json:"blocklist_votes"`
	CreatedAt      time.Time  `json:"created_at"`
	FinalizedAt    *time.Time `json:"finalized_at,omitempty"`
}

type ProposalStatus string

const (
	ProposalStatusPending  ProposalStatus = "PENDING"
	ProposalStatusApproved ProposalStatus = "APPROVED"
	ProposalStatusRejected ProposalStatus = "REJECTED"
)

// ProposalWithCreator includes creator information
type ProposalWithCreator struct {
	Proposal
	CreatorUsername string  `json:"creator_username"`
	CreatorEmail    *string `json:"creator_email,omitempty"`
}
