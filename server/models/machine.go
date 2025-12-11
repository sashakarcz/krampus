package models

import (
	"time"
)

type Machine struct {
	ID                int64      `json:"id"`
	MachineID         string     `json:"machine_id"`
	SerialNumber      *string    `json:"serial_number,omitempty"`
	Hostname          *string    `json:"hostname,omitempty"`
	OSVersion         *string    `json:"os_version,omitempty"`
	OSBuild           *string    `json:"os_build,omitempty"`
	SantaVersion      *string    `json:"santa_version,omitempty"`
	ClientMode        *string    `json:"client_mode,omitempty"` // "MONITOR" or "LOCKDOWN"
	EnrolledAt        time.Time  `json:"enrolled_at"`
	LastSync          *time.Time `json:"last_sync,omitempty"`
	LastPreflightSync *time.Time `json:"last_preflight_sync,omitempty"`
}

type ClientMode string

const (
	ClientModeMonitor  ClientMode = "MONITOR"
	ClientModeLockdown ClientMode = "LOCKDOWN"
)
