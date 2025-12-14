#!/bin/bash

# Santa Krampus Server Configuration Script
# Downloads and installs configuration profile from Krampus server for NorthPole Security Santa

set -euo pipefail

# Configuration variables
KRAMPUS_SERVER="https://krampus.starnix.net"
MACHINE_ID="finnmac"
CLIENT_MODE="MONITOR"          # MONITOR or LOCKDOWN
FULL_SYNC_INTERVAL=600         # Sync interval in seconds (10 minutes)
ORGANIZATION_NAME="Starnix"    # Optional organization name

# Authentication (if your Krampus server requires auth for /api/machines endpoint)
# Leave empty if not using authentication or if the endpoint is public
AUTH_TOKEN=""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
    fi
}

# Check if Santa is installed
check_santa_installed() {
    if ! command -v santactl >/dev/null 2>&1; then
        log_error "santactl command not found. Please install NorthPole Security Santa first."
    fi
    log_info "Santa installation detected: $(santactl version 2>/dev/null | head -1 || echo 'version unknown')"
}

# Download mobileconfig from Krampus server
download_mobileconfig() {
    local config_file="/tmp/${MACHINE_ID}.mobileconfig"

    log_info "Downloading configuration profile from Krampus server..."

    # Build curl command with optional auth header
    local curl_cmd="curl -s -S -f"
    if [ -n "${AUTH_TOKEN}" ]; then
        curl_cmd="${curl_cmd} -H 'Authorization: Bearer ${AUTH_TOKEN}'"
    fi

    # Download the mobileconfig
    local payload=$(cat <<EOF
{
    "client_mode": "${CLIENT_MODE}",
    "upload_interval": ${FULL_SYNC_INTERVAL},
    "organization_name": "${ORGANIZATION_NAME}"
}
EOF
)

    if eval "${curl_cmd} -X POST -H 'Content-Type: application/json' \
        -d '${payload}' \
        '${KRAMPUS_SERVER}/api/machines/${MACHINE_ID}/mobileconfig' \
        -o '${config_file}'"; then
        log_info "Configuration profile downloaded successfully"
        echo "${config_file}"
    else
        log_error "Failed to download configuration profile from ${KRAMPUS_SERVER}"
    fi
}

# Install configuration profile
install_profile() {
    local config_file="$1"

    log_info "Installing configuration profile..."

    # Install the profile using profiles command
    if profiles install -path="${config_file}" -type=configuration; then
        log_info "Configuration profile installed successfully"

        # Clean up the downloaded file
        rm -f "${config_file}"
    else
        log_error "Failed to install configuration profile"
    fi
}

# Restart Santa services to pick up new configuration
restart_santa_services() {
    log_info "Restarting Santa services to apply new configuration..."

    # Unload services
    launchctl bootout system /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist 2>/dev/null || \
        launchctl unload /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist 2>/dev/null || true

    launchctl bootout system /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist 2>/dev/null || \
        launchctl unload /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist 2>/dev/null || true

    sleep 2

    # Load services
    launchctl bootstrap system /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist 2>/dev/null || \
        launchctl load /Library/LaunchDaemons/com.northpolesec.santa.daemon.plist 2>/dev/null || true

    launchctl bootstrap system /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist 2>/dev/null || \
        launchctl load /Library/LaunchDaemons/com.northpolesec.santa.syncservice.plist 2>/dev/null || true

    log_info "Santa services restarted"
    sleep 3
}

# Attempt initial sync
attempt_initial_sync() {
    log_info "Attempting initial sync with Krampus server..."

    if santactl sync --debug 2>&1; then
        log_info "✓ Initial sync successful!"
        return 0
    else
        log_warn "✗ Initial sync failed - this may be normal on first run"
        log_warn "The server will register the machine on the next automatic sync"
        return 1
    fi
}

# Display Santa status
show_status() {
    log_info "Current Santa status:"
    echo ""
    santactl status 2>&1 || log_warn "Failed to get Santa status"
    echo ""
}

# Main function
main() {
    log_info "=== Krampus Santa Sync Configuration ==="
    echo ""

    check_root
    check_santa_installed

    echo ""
    local config_file=$(download_mobileconfig)

    echo ""
    install_profile "${config_file}"

    echo ""
    restart_santa_services

    echo ""
    attempt_initial_sync || true

    echo ""
    show_status

    echo ""
    log_info "=== Configuration Complete ==="
    log_info "Santa is now configured to sync with ${KRAMPUS_SERVER}"
    log_info "Machine ID: ${MACHINE_ID}"
    log_info "Mode: ${CLIENT_MODE}"
    echo ""
    log_info "Next steps:"
    echo "  1. Check the Krampus web UI to verify the machine registered"
    echo "  2. Monitor logs: sudo log stream --predicate 'processImagePath contains \"santa\"' --level debug"
    echo "  3. Force sync anytime: sudo santactl sync"
    echo "  4. View installed profiles: sudo profiles list"
    echo ""

    if [ "${CLIENT_MODE}" = "MONITOR" ]; then
        log_info "Running in MONITOR mode - all binaries will be allowed but logged"
    else
        log_warn "Running in LOCKDOWN mode - unapproved binaries will be BLOCKED!"
    fi
}

# Run main function
main "$@"
