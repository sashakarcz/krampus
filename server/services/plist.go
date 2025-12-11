package services

import (
	"fmt"
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
