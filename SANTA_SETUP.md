# Santa Client Configuration Guide

This guide explains how to configure NorthPole Security Santa clients to sync with your Krampus server.

## Understanding NorthPole Security Santa Configuration

NorthPole Security Santa (the successor to Google Santa) uses **Configuration Profiles** (`.mobileconfig` files) for deployment, not direct plist manipulation. This is the officially supported method.

## Method 1: Web UI Download (Recommended)

1. **Log into Krampus Web UI**
   - Navigate to `https://krampus.starnix.net`
   - Log in with your OIDC credentials

2. **Register Your Machine**
   - Go to the "Machines" page
   - Click "Register New Machine"
   - Enter Machine ID (e.g., `finnmac`)
   - Enter Serial Number (optional)
   - Click "Register"

3. **Download Configuration Profile**
   - Find your registered machine in the list
   - Click "Download mobileconfig"
   - Choose settings:
     - **Client Mode**: MONITOR (logs only) or LOCKDOWN (enforce blocking)
     - **Upload Interval**: 600 seconds (10 minutes) recommended
     - **Organization Name**: Your company name (optional)
   - Save the `.mobileconfig` file

4. **Install Configuration Profile on Mac**
   ```bash
   # Copy the downloaded .mobileconfig to your Mac
   scp finnmac.mobileconfig user@finnmac:/tmp/

   # On the Mac, install it with sudo
   sudo profiles install -path=/tmp/finnmac.mobileconfig -type=configuration

   # Verify installation
   sudo profiles list | grep northpole
   ```

5. **Restart Santa Services**
   ```bash
   # Modern macOS (Big Sur+)
   sudo launchctl bootout system /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist
   sudo launchctl bootout system /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist
   sleep 2
   sudo launchctl bootstrap system /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist
   sudo launchctl bootstrap system /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist

   # Or on older macOS
   sudo launchctl unload /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist
   sudo launchctl unload /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist
   sleep 2
   sudo launchctl load /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist
   sudo launchctl load /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist
   ```

6. **Verify Configuration**
   ```bash
   # Check Santa status - Sync should show "Enabled | Yes"
   santactl status

   # Manually trigger a sync
   sudo santactl sync

   # Monitor Santa logs
   sudo log stream --predicate 'processImagePath contains "santa"' --level debug
   ```

## Method 2: Automated Script (Requires Auth Token)

If you need to automate deployment, use the included `santa-config.sh` script:

1. **Get Authentication Token**
   - Log into Krampus web UI
   - Open browser developer tools (F12)
   - Go to Application/Storage → Cookies
   - Copy the `jwt` cookie value

2. **Edit Configuration Script**
   ```bash
   vi santa-config.sh

   # Update these variables:
   KRAMPUS_SERVER="https://krampus.starnix.net"
   MACHINE_ID="finnmac"  # Change to your machine ID
   CLIENT_MODE="MONITOR"  # or "LOCKDOWN"
   AUTH_TOKEN="your-jwt-token-here"  # Paste JWT token
   ```

3. **Run Script**
   ```bash
   sudo bash santa-config.sh
   ```

## Method 3: Manual Profile Creation

If you can't use the web UI, you can manually create a configuration profile:

1. **Create Configuration File**
   ```bash
   sudo vi /tmp/santa-config.mobileconfig
   ```

2. **Paste This Template** (replace values as needed):
   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
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
               <string>com.northpolesec.santa.config</string>
               <key>PayloadOrganization</key>
               <string>Your Organization</string>
               <key>PayloadType</key>
               <string>com.northpolesec.santa</string>
               <key>PayloadUUID</key>
               <string>CHANGE-THIS-UUID-1234-5678</string>
               <key>PayloadVersion</key>
               <integer>1</integer>
               <key>SyncBaseURL</key>
               <string>https://krampus.starnix.net</string>
               <key>ClientMode</key>
               <integer>1</integer>
               <key>MachineID</key>
               <string>finnmac</string>
               <key>FullSyncInterval</key>
               <integer>600</integer>
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
       <string>Configures NorthPole Security Santa endpoint security</string>
       <key>PayloadDisplayName</key>
       <string>Santa Configuration</string>
       <key>PayloadIdentifier</key>
       <string>com.northpolesec.santa.config</string>
       <key>PayloadOrganization</key>
       <string>Your Organization</string>
       <key>PayloadRemovalDisallowed</key>
       <false/>
       <key>PayloadScope</key>
       <string>System</string>
       <key>PayloadType</key>
       <string>Configuration</string>
       <key>PayloadUUID</key>
       <string>CHANGE-THIS-UUID-MAIN-PROFILE</string>
       <key>PayloadVersion</key>
       <integer>1</integer>
   </dict>
   </plist>
   ```

3. **Key Values to Update**:
   - `SyncBaseURL`: Your Krampus server URL
   - `MachineID`: Unique identifier for this Mac
   - `ClientMode`: 1 = MONITOR (logging only), 2 = LOCKDOWN (enforcement)
   - `PayloadUUID`: Generate unique UUIDs (use `uuidgen` command)

4. **Install the Profile**
   ```bash
   sudo profiles install -path=/tmp/santa-config.mobileconfig -type=configuration
   ```

## Troubleshooting

### Sync Shows "Enabled | No"

This means the configuration profile wasn't applied correctly. Check:

1. **Verify Profile Installation**
   ```bash
   sudo profiles list | grep -A 20 northpole
   ```

2. **Check for Conflicting Configurations**
   ```bash
   # Remove old preference files if they exist
   sudo rm -f /Library/Preferences/com.northpolesec.santa.plist
   sudo rm -f /Library/Managed\ Preferences/com.northpolesec.santa.plist
   ```

3. **Restart Santa Completely**
   ```bash
   sudo launchctl bootout system /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist
   sudo launchctl bootout system /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist
   sleep 5
   sudo launchctl bootstrap system /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist
   sudo launchctl bootstrap system /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist
   sleep 5
   santactl status
   ```

### Connection Refused / Sync Fails

1. **Check Network Connectivity**
   ```bash
   curl -v https://krampus.starnix.net/ping
   ```

2. **Verify Machine is Registered**
   - Check Krampus web UI → Machines page
   - Machine must be registered before first sync

3. **Check Santa Logs**
   ```bash
   sudo log stream --predicate 'processImagePath contains "santa"' --level debug
   ```

### Certificate Errors

If using HTTPS with self-signed certificates, you may need to add the CA certificate to the configuration profile using the `ServerAuthRootsData` or `ServerAuthRootsFile` keys.

## Client Mode Explanation

- **MONITOR (1)**: All binaries are allowed to execute. Santa logs execution events and sends them to the Krampus server. Use this for initial deployment and testing.

- **LOCKDOWN (2)**: Only explicitly allowed binaries can execute. Unknown binaries are blocked. Use this for production security enforcement.

## References

- [NorthPole Security Santa Documentation](https://northpole.dev/)
- [Configuration Keys](https://northpole.dev/configuration/keys/)
- [Profile Configuration](https://northpole.dev/deployment/profile-configuration/)
- [Sync Servers](https://northpole.dev/features/sync/)

## Support

For issues:
1. Check Santa status: `santactl status`
2. Review logs: `sudo log stream --predicate 'processImagePath contains "santa"' --level debug`
3. Verify Krampus server is running: `curl https://krampus.starnix.net/ping`
4. Check Krampus server logs for sync errors
