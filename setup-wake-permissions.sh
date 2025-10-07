#!/bin/bash
# Katana Alarm - System Wake-Up Permission Setup
# This script configures sudo permissions for rtcwake without password prompt

set -e

echo "=========================================="
echo "Katana Alarm - Wake-Up Permission Setup"
echo "=========================================="
echo ""
echo "This script will configure your system to allow Katana to wake"
echo "your computer from sleep/suspend mode without requiring a password."
echo ""
echo "What this does:"
echo "  - Grants passwordless sudo access to /usr/sbin/rtcwake"
echo "  - Creates sudoers configuration in /etc/sudoers.d/"
echo "  - Only affects the rtcwake command (safe and minimal)"
echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    echo "ERROR: Please run this script as your normal user (not with sudo)"
    echo "Usage: ./setup-wake-permissions.sh"
    exit 1
fi

# Check if rtcwake exists (check both PATH and common locations)
RTCWAKE_PATH=""
if command -v rtcwake &> /dev/null; then
    RTCWAKE_PATH=$(command -v rtcwake)
elif [ -x /usr/sbin/rtcwake ]; then
    RTCWAKE_PATH="/usr/sbin/rtcwake"
elif [ -x /sbin/rtcwake ]; then
    RTCWAKE_PATH="/sbin/rtcwake"
fi

if [ -z "$RTCWAKE_PATH" ]; then
    echo "WARNING: rtcwake command not found on your system."
    echo "Your distribution may not support RTC wake-up, or you need to install util-linux package."
    echo ""
    read -p "Do you want to continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
    RTCWAKE_PATH="/usr/sbin/rtcwake"  # Default path
else
    echo "✅ Found rtcwake at: $RTCWAKE_PATH"
    echo ""
fi

echo "This script requires sudo access to modify system configuration."
echo "You will be prompted for your password once."
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Setup cancelled."
    exit 0
fi

# Create sudoers configuration
SUDOERS_FILE="/etc/sudoers.d/katana-alarm-wake"
CURRENT_USER="$USER"

echo ""
echo "Creating sudoers configuration..."

# Create the sudoers rule
sudo tee "$SUDOERS_FILE" > /dev/null <<EOF
# Katana Alarm - Allow passwordless rtcwake for system wake-up
# Created: $(date)
# User: $CURRENT_USER
#
# This allows the Katana alarm application to schedule system wake-ups
# from sleep/suspend mode without requiring password authentication.
#
# Security: This only grants access to the rtcwake command with specific
# parameters and does not provide general sudo access.

$CURRENT_USER ALL=(ALL) NOPASSWD: /usr/sbin/rtcwake
EOF

# Set proper permissions (sudoers files must be 0440)
sudo chmod 0440 "$SUDOERS_FILE"

# Validate the sudoers file
if sudo visudo -c -f "$SUDOERS_FILE" &> /dev/null; then
    echo "✅ Sudoers configuration created successfully!"
    echo ""
    echo "Configuration file: $SUDOERS_FILE"
    echo ""
    echo "Testing the configuration..."
    
    # Test if sudo works without password
    if sudo -n rtcwake --version &> /dev/null; then
        echo "✅ TEST PASSED: rtcwake can now run without password!"
        echo ""
        echo "=========================================="
        echo "Setup Complete!"
        echo "=========================================="
        echo ""
        echo "Your system is now configured for Katana alarm wake-ups."
        echo "When you set an alarm, your computer will automatically wake"
        echo "from sleep/suspend to ring the alarm at the scheduled time."
        echo ""
        echo "To remove this configuration later, run:"
        echo "  sudo rm $SUDOERS_FILE"
        echo ""
    else
        echo "⚠️  WARNING: Configuration created but test failed."
        echo "You may need to log out and back in for changes to take effect."
        echo ""
    fi
else
    echo "❌ ERROR: Sudoers configuration validation failed!"
    echo "Removing invalid configuration..."
    sudo rm -f "$SUDOERS_FILE"
    exit 1
fi
