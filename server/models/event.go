package models

import (
	"time"
)

type Event struct {
	ID                 int64      `json:"id"`
	MachineID          string     `json:"machine_id"`
	FilePath           *string    `json:"file_path,omitempty"`
	FileHash           string     `json:"file_hash"`
	ExecutionTime      time.Time  `json:"execution_time"`
	Decision           *string    `json:"decision,omitempty"`
	ExecutingUser      *string    `json:"executing_user,omitempty"`
	CertSHA256         *string    `json:"cert_sha256,omitempty"`
	CertCN             *string    `json:"cert_cn,omitempty"`
	BundleID           *string    `json:"bundle_id,omitempty"`
	BundleName         *string    `json:"bundle_name,omitempty"`
	BundlePath         *string    `json:"bundle_path,omitempty"`
	SigningID          *string    `json:"signing_id,omitempty"`
	TeamID             *string    `json:"team_id,omitempty"`
	QuarantineDataURL  *string    `json:"quarantine_data_url,omitempty"`
	QuarantineTimestamp *time.Time `json:"quarantine_timestamp,omitempty"`
}

// SantaEvent represents an event in the Santa sync protocol format
type SantaEvent struct {
	FileSHA256          string    `json:"file_sha256"`
	FilePath            string    `json:"file_path"`
	FileName            string    `json:"file_name"`
	ExecutingUser       string    `json:"executing_user"`
	ExecutionTime       float64   `json:"execution_time"` // Unix timestamp
	Decision            string    `json:"decision"` // "ALLOW", "BLOCK", etc.
	LoggedInUsers       []string  `json:"logged_in_users,omitempty"`
	CurrentSessions     []string  `json:"current_sessions,omitempty"`
	CertificateSHA256   string    `json:"certificate_sha256,omitempty"`
	CertificateCN       string    `json:"certificate_cn,omitempty"`
	TeamID              string    `json:"team_id,omitempty"`
	SigningID           string    `json:"signing_id,omitempty"`
	CDHash              string    `json:"cdhash,omitempty"`
	BundleID            string    `json:"bundle_id,omitempty"`
	BundleName          string    `json:"bundle_name,omitempty"`
	BundlePath          string    `json:"bundle_path,omitempty"`
	BundleVersionString string    `json:"bundle_version_string,omitempty"`
	BundleVersion       string    `json:"bundle_version,omitempty"`
	QuarantineDataURL   string    `json:"quarantine_data_url,omitempty"`
	QuarantineTimestamp float64   `json:"quarantine_timestamp,omitempty"`
	PID                 int       `json:"pid,omitempty"`
	PPID                int       `json:"ppid,omitempty"`
	ParentName          string    `json:"parent_name,omitempty"`
}
