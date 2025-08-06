#!/usr/bin/env bash
# Simple uninstallation script for ctop

# Color definitions
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Output functions
function output() { echo -e "${GREEN}[ctop-uninstall]${NC} $*"; }
function log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
function log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check if command exists
function command_exists() {
    command -v "$@" > /dev/null 2>&1
}

# Determine if we need sudo/su
sh_c='sh -c'
if [[ $EUID -ne 0 ]]; then
    if command_exists sudo; then
        log_warning "sudo is required to uninstall ctop"
        sh_c='sudo -E sh -c'
    elif command_exists su; then
        log_warning "su is required to uninstall ctop"
        sh_c='su -c'
    else
        log_error "This uninstaller needs root privileges. Neither 'sudo' nor 'su' were found."
        exit 1
    fi
fi

# Default installation path
CTOP_PATH="/usr/local/bin/ctop"

output "Checking for ctop installation..."
if [ -f "$CTOP_PATH" ]; then
    output "Found ctop at $CTOP_PATH"

    # Get version info before uninstalling
    version=$("$CTOP_PATH" -v 2>/dev/null || echo "unknown version")
    
    output "Removing ctop ($version)..."
    if $sh_c "rm -f $CTOP_PATH"; then
        output "ctop was successfully uninstalled"
    else
        log_error "Failed to remove ctop"
        exit 1
    fi
else
    output "ctop is not installed at $CTOP_PATH"
    
    # Check if it's installed elsewhere
    alternative_path=$(command -v ctop 2>/dev/null)
    if [ -n "$alternative_path" ]; then
        output "Found ctop at $alternative_path"
        read -p "Do you want to remove it? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            if $sh_c "rm -f $alternative_path"; then
                output "ctop was successfully uninstalled"
            else
                log_error "Failed to remove ctop"
                exit 1
            fi
        else
            output "Uninstallation cancelled"
        fi
    else
        output "No ctop installation found"
    fi
fi

output "Uninstallation complete"
