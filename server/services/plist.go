package services

import (
	"fmt"
	"github.com/google/uuid"
)

// GeneratePlist generates a Santa configuration plist file
func GeneratePlist(machineID, clientMode, syncBaseURL string, uploadInterval int) string {
	// Convert client mode to integer
	// 1 = MONITOR, 2 = LOCKDOWN
	mode := 1
	if clientMode == "LOCKDOWN" {
		mode = 2
	}

	plistTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>SyncBaseURL</key>
	<string>%s</string>
	<key>ClientMode</key>
	<integer>%d</integer>
	<key>MachineID</key>
	<string>%s</string>
	<key>FullSyncInterval</key>
	<integer>%d</integer>
	<key>EventLogType</key>
	<string>protobuf</string>
	<key>EventLogPath</key>
	<string>/var/db/santa/events.pb</string>
	<key>EnablePageZeroProtection</key>
	<true/>
	<key>MachineOwner</key>
	<string></string>
	<key>EnableBundles</key>
	<true/>
	<key>EnableTransitiveRules</key>
	<false/>
</dict>
</plist>`

	return fmt.Sprintf(plistTemplate, syncBaseURL, mode, machineID, uploadInterval)
}

// GenerateMobileConfig generates an Apple Configuration Profile (.mobileconfig)
// for easy Santa configuration deployment
func GenerateMobileConfig(machineID, clientMode, syncBaseURL, organizationName string, uploadInterval int) string {
	// Convert client mode to integer
	mode := 1
	if clientMode == "LOCKDOWN" {
		mode = 2
	}

	// Generate UUIDs for the profile and payload
	profileUUID := uuid.New().String()
	payloadUUID := uuid.New().String()

	if organizationName == "" {
		organizationName = "Krampus Santa Sync"
	}

	mobileConfigTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>PayloadDescription</key>
			<string>Santa Configuration</string>
			<key>PayloadDisplayName</key>
			<string>Santa</string>
			<key>PayloadIdentifier</key>
			<string>com.northpolesec.santa.%s</string>
			<key>PayloadOrganization</key>
			<string>%s</string>
			<key>PayloadType</key>
			<string>com.northpolesec.santa</string>
			<key>PayloadUUID</key>
			<string>%s</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
			<key>SyncBaseURL</key>
			<string>%s</string>
			<key>ClientMode</key>
			<integer>%d</integer>
			<key>MachineID</key>
			<string>%s</string>
			<key>FullSyncInterval</key>
			<integer>%d</integer>
			<key>EventLogType</key>
			<string>protobuf</string>
			<key>EventLogPath</key>
			<string>/var/db/santa/events.pb</string>
			<key>EnablePageZeroProtection</key>
			<true/>
			<key>EnableBundles</key>
			<true/>
			<key>EnableTransitiveRules</key>
			<false/>
		</dict>
	</array>
	<key>PayloadDescription</key>
	<string>Configures NorthPole Security Santa endpoint security for %s</string>
	<key>PayloadDisplayName</key>
	<string>Santa Configuration - %s</string>
	<key>PayloadIdentifier</key>
	<string>com.northpolesec.santa.config.%s</string>
	<key>PayloadOrganization</key>
	<string>%s</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
	<key>PayloadScope</key>
	<string>System</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>%s</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
</dict>
</plist>`

	return fmt.Sprintf(mobileConfigTemplate,
		payloadUUID, organizationName, payloadUUID,
		syncBaseURL, mode, machineID, uploadInterval,
		machineID, machineID, profileUUID, organizationName, profileUUID)
}
