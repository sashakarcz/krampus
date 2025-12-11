package handlers

import (
	"krampus/server/database"
	"krampus/server/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Preflight handles the Santa preflight request
func Preflight(c *gin.Context) {
	machineID := c.Param("machine_id")

	// Santa can send either JSON or URL-encoded forms
	var input struct {
		Hostname      string `json:"primary_user" form:"primary_user"`
		OSVersion     string `json:"os_version" form:"os_version"`
		OSBuild       string `json:"os_build" form:"os_build"`
		SantaVersion  string `json:"santa_version" form:"santa_version"`
		SerialNumber  string `json:"serial_num" form:"serial_num"`
		ClientMode    string `json:"client_mode" form:"client_mode"`
		ModelID       string `json:"model_identifier" form:"model_identifier"`
	}

	// Try to bind as form first, then JSON
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Preflight parse error for machine %s: %v", machineID, err)
		log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))
		// Continue anyway - Santa might send minimal data
	}

	// Update or insert machine information
	_, err := database.DB.Exec(
		`INSERT INTO machines (machine_id, serial_number, hostname, os_version, os_build, santa_version, client_mode, last_preflight_sync)
		 VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
		 ON CONFLICT(machine_id) DO UPDATE SET
		   serial_number = excluded.serial_number,
		   hostname = excluded.hostname,
		   os_version = excluded.os_version,
		   os_build = excluded.os_build,
		   santa_version = excluded.santa_version,
		   client_mode = excluded.client_mode,
		   last_preflight_sync = datetime('now')`,
		machineID, input.SerialNumber, input.Hostname, input.OSVersion,
		input.OSBuild, input.SantaVersion, input.ClientMode,
	)
	if err != nil {
		log.Printf("Failed to update machine: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update machine"})
		return
	}

	// Return sync configuration
	c.JSON(http.StatusOK, gin.H{
		"client_mode":               "LOCKDOWN",
		"batch_size":                100,
		"upload_logs_url":           "",
		"clean_sync":                false,
		"enable_bundles":            true,
		"enable_transitive_rules":   false,
		"blocked_path_regex":        "",
		"allowed_path_regex":        "",
		"enable_all_event_upload":   false,
	})
}

// EventUpload handles Santa event upload
func EventUpload(c *gin.Context) {
	machineID := c.Param("machine_id")

	var input struct {
		Events []models.SantaEvent `json:"events"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("EventUpload parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert events into database
	for _, event := range input.Events {
		execTime := time.Unix(int64(event.ExecutionTime), 0)

		_, err := database.DB.Exec(
			`INSERT INTO events (machine_id, file_path, file_hash, execution_time, decision,
			                     executing_user, cert_sha256, cert_cn, bundle_id, bundle_name,
			                     bundle_path, signing_id, team_id, quarantine_data_url)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			machineID, event.FilePath, event.FileSHA256, execTime, event.Decision,
			event.ExecutingUser, event.CertificateSHA256, event.CertificateCN,
			event.BundleID, event.BundleName, event.BundlePath,
			event.SigningID, event.TeamID, event.QuarantineDataURL,
		)
		if err != nil {
			log.Printf("Failed to insert event: %v", err)
			continue
		}
	}

	// Return acknowledgement
	c.JSON(http.StatusOK, gin.H{
		"event_upload_bundle_binaries": []string{},
	})
}

// RuleDownload handles Santa rule download requests
func RuleDownload(c *gin.Context) {
	machineID := c.Param("machine_id")

	var input struct {
		Cursor string `json:"cursor"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Cursor is optional, ignore error
	}

	// Parse cursor (if provided, it's the last rule ID)
	var startID int64 = 0
	if input.Cursor != "" {
		// In a production system, you might want to properly parse/validate the cursor
		// For now, we'll use a simple approach
	}

	// Fetch rules with pagination
	batchSize := 100
	rows, err := database.DB.Query(
		`SELECT id, identifier, policy, rule_type, custom_message
		 FROM rules
		 WHERE id > ?
		 ORDER BY id
		 LIMIT ?`,
		startID, batchSize,
	)
	if err != nil {
		log.Printf("Failed to query rules: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules"})
		return
	}
	defer rows.Close()

	santaRules := []models.SantaRule{}
	var lastID int64

	for rows.Next() {
		var id int64
		var r models.SantaRule
		err := rows.Scan(&id, &r.Identifier, &r.Policy, &r.RuleType, &r.CustomMsg)
		if err != nil {
			log.Printf("Failed to scan rule: %v", err)
			continue
		}
		santaRules = append(santaRules, r)
		lastID = id
	}

	// Prepare response
	response := gin.H{
		"rules": santaRules,
	}

	// If we got a full batch, there might be more rules
	if len(santaRules) == batchSize {
		response["cursor"] = lastID
	}

	// Update last sync time
	_, _ = database.DB.Exec(
		`UPDATE machines SET last_sync = datetime('now') WHERE machine_id = ?`,
		machineID,
	)

	c.JSON(http.StatusOK, response)
}

// Postflight handles Santa postflight requests
func Postflight(c *gin.Context) {
	machineID := c.Param("machine_id")

	// Log the completion of sync
	_, err := database.DB.Exec(
		`UPDATE machines SET last_sync = datetime('now') WHERE machine_id = ?`,
		machineID,
	)
	if err != nil {
		log.Printf("Failed to update machine last_sync: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{})
}
