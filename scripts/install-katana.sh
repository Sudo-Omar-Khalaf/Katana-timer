#!/bin/bash
# Katana Multi-Timer - Automated Installation Script
# Handles: rtcwake installation, path detection, sudo setup, and build

set -e

INSTALL_DIR="$HOME/.local/bin"
KATANA_BINARY="katana"
SUDOERS_FILE="/etc/sudoers.d/katana-alarm-wake"

echo "=========================================="
echo "🥷 Katana Multi-Timer Installation"
echo "=========================================="
echo ""

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "ℹ️  $1"
}

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID"
    elif [ -f /etc/lsb-release ]; then
        . /etc/lsb-release
        echo "$DISTRIB_ID" | tr '[:upper:]' '[:lower:]'
    else
        echo "unknown"
    fi
}

DISTRO=$(detect_distro)
print_info "Detected distribution: $DISTRO"
echo ""

# Step 1: Check and install rtcwake if needed
echo "Step 1: Checking for rtcwake..."
echo "================================"

RTCWAKE_PATH=""
if command -v rtcwake &> /dev/null; then
    RTCWAKE_PATH=$(command -v rtcwake)
    print_success "rtcwake found at: $RTCWAKE_PATH"
elif [ -x /usr/sbin/rtcwake ]; then
    RTCWAKE_PATH="/usr/sbin/rtcwake"
    print_success "rtcwake found at: $RTCWAKE_PATH"
elif [ -x /sbin/rtcwake ]; then
    RTCWAKE_PATH="/sbin/rtcwake"
    print_success "rtcwake found at: $RTCWAKE_PATH"
else
    print_warning "rtcwake not found. Installing util-linux package..."
    
    case "$DISTRO" in
        ubuntu|debian|kali|parrot|pop|mint|elementary)
            print_info "Installing via apt..."
            sudo apt-get update -qq
            sudo apt-get install -y util-linux
            ;;
        fedora|rhel|centos|rocky|alma)
            print_info "Installing via dnf/yum..."
            sudo dnf install -y util-linux || sudo yum install -y util-linux
            ;;
        arch|manjaro|endeavouros)
            print_info "Installing via pacman..."
            sudo pacman -S --noconfirm util-linux
            ;;
        opensuse*|sles)
            print_info "Installing via zypper..."
            sudo zypper install -y util-linux
            ;;
        *)
            print_warning "Unknown distribution. Attempting generic installation..."
            if command -v apt-get &> /dev/null; then
                sudo apt-get update -qq && sudo apt-get install -y util-linux
            elif command -v dnf &> /dev/null; then
                sudo dnf install -y util-linux
            elif command -v yum &> /dev/null; then
                sudo yum install -y util-linux
            elif command -v pacman &> /dev/null; then
                sudo pacman -S --noconfirm util-linux
            else
                print_error "Could not install util-linux automatically."
                print_info "Please install it manually and run this script again."
                exit 1
            fi
            ;;
    esac
    
    # Check again after installation
    if command -v rtcwake &> /dev/null; then
        RTCWAKE_PATH=$(command -v rtcwake)
        print_success "rtcwake installed successfully at: $RTCWAKE_PATH"
    elif [ -x /usr/sbin/rtcwake ]; then
        RTCWAKE_PATH="/usr/sbin/rtcwake"
        print_success "rtcwake installed successfully at: $RTCWAKE_PATH"
    else
        print_error "Failed to install rtcwake. Wake-up feature may not work."
        RTCWAKE_PATH="/usr/sbin/rtcwake"  # Default fallback
    fi
fi
echo ""

# Step 2: Configure sudo permissions for rtcwake
echo "Step 2: Configuring wake-up permissions..."
echo "==========================================="

if sudo -n rtcwake --version &> /dev/null 2>&1; then
    print_success "rtcwake sudo permissions already configured!"
else
    print_info "Setting up passwordless sudo for rtcwake..."
    print_info "You will be prompted for your sudo password once."
    echo ""
    
    # Create sudoers configuration
    sudo tee "$SUDOERS_FILE" > /dev/null <<EOF
# Katana Alarm - Allow passwordless rtcwake for system wake-up
# Created: $(date)
# User: $USER
#
# This allows the Katana alarm application to schedule system wake-ups
# from sleep/suspend mode without requiring password authentication.
#
# Security: This only grants access to the rtcwake command and does not
# provide general sudo access.

$USER ALL=(ALL) NOPASSWD: $RTCWAKE_PATH
EOF

    # Set proper permissions
    sudo chmod 0440 "$SUDOERS_FILE"
    
    # Validate sudoers file
    if sudo visudo -c -f "$SUDOERS_FILE" &> /dev/null; then
        print_success "Sudo permissions configured successfully!"
        
        # Test the configuration
        if sudo -n $RTCWAKE_PATH --version &> /dev/null 2>&1; then
            print_success "Wake-up permissions verified!"
        else
            print_warning "Permissions configured but verification failed."
            print_info "You may need to restart your terminal session."
        fi
    else
        print_error "Sudoers configuration validation failed!"
        sudo rm -f "$SUDOERS_FILE"
        print_warning "Continuing without wake-up permissions..."
        print_info "Alarms will work when system is awake."
    fi
fi
echo ""

# Step 3: Check Go installation
echo "Step 3: Checking Go installation..."
echo "===================================="

if ! command -v go &> /dev/null; then
    print_error "Go is not installed!"
    print_info "Please install Go 1.19 or later from: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_success "Go $GO_VERSION found"
echo ""

# Step 4: Check build dependencies
echo "Step 4: Checking build dependencies..."
echo "======================================="

MISSING_DEPS=()

# Check for required packages
if ! pkg-config --exists gl 2>/dev/null; then
    MISSING_DEPS+=("OpenGL")
fi

if ! pkg-config --exists alsa 2>/dev/null; then
    MISSING_DEPS+=("ALSA")
fi

if ! command -v gcc &> /dev/null; then
    MISSING_DEPS+=("GCC")
fi

if [ ${#MISSING_DEPS[@]} -gt 0 ]; then
    print_warning "Missing dependencies: ${MISSING_DEPS[*]}"
    print_info "Installing build dependencies..."
    
    case "$DISTRO" in
        ubuntu|debian|kali|parrot|pop|mint|elementary)
            sudo apt-get install -y build-essential pkg-config libgl1-mesa-dev libasound2-dev
            ;;
        fedora|rhel|centos|rocky|alma)
            sudo dnf install -y gcc pkg-config mesa-libGL-devel alsa-lib-devel || \
            sudo yum install -y gcc pkg-config mesa-libGL-devel alsa-lib-devel
            ;;
        arch|manjaro|endeavouros)
            sudo pacman -S --noconfirm base-devel pkg-config mesa alsa-lib
            ;;
        opensuse*|sles)
            sudo zypper install -y gcc pkg-config Mesa-libGL-devel alsa-devel
            ;;
    esac
    
    print_success "Dependencies installed!"
else
    print_success "All dependencies present!"
fi
echo ""

# Step 5: Build Katana
echo "Step 5: Building Katana..."
echo "=========================="

print_info "Downloading Go dependencies..."
go mod download

print_info "Building application..."
if go build -o "$KATANA_BINARY" -ldflags="-s -w" .; then
    print_success "Build successful!"
else
    print_error "Build failed!"
    exit 1
fi
echo ""

# Step 6: Install binary
echo "Step 6: Installing Katana..."
echo "============================"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Copy binary
cp "$KATANA_BINARY" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$KATANA_BINARY"

print_success "Katana installed to: $INSTALL_DIR/$KATANA_BINARY"

# Check if install dir is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_warning "$INSTALL_DIR is not in your PATH"
    print_info "Adding to PATH in your shell configuration..."
    
    # Detect shell and add to appropriate config file
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        SHELL_RC="$HOME/.bashrc"
    else
        SHELL_RC="$HOME/.profile"
    fi
    
    if ! grep -q "$INSTALL_DIR" "$SHELL_RC" 2>/dev/null; then
        echo "" >> "$SHELL_RC"
        echo "# Added by Katana installer" >> "$SHELL_RC"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
        print_success "PATH updated in $SHELL_RC"
        print_info "Run: source $SHELL_RC (or restart your terminal)"
    fi
else
    print_success "Installation directory is in PATH"
fi
echo ""

# Step 7: Create desktop entry (optional)
echo "Step 7: Creating desktop launcher..."
echo "====================================="

DESKTOP_DIR="$HOME/.local/share/applications"
DESKTOP_FILE="$DESKTOP_DIR/katana-timer.desktop"

mkdir -p "$DESKTOP_DIR"

cat > "$DESKTOP_FILE" <<EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Katana Multi-Timer
Comment=Advanced time tracking, alarms, and productivity tool
Exec=$INSTALL_DIR/$KATANA_BINARY
Icon=preferences-system-time
Terminal=false
Categories=Utility;Clock;
Keywords=timer;alarm;stopwatch;countdown;productivity;
EOF

chmod +x "$DESKTOP_FILE"
print_success "Desktop launcher created!"
echo ""

# Installation complete!
echo "=========================================="
echo "🎉 Installation Complete!"
echo "=========================================="
echo ""
print_success "Katana Multi-Timer has been installed successfully!"
echo ""
echo "📍 Installation Location: $INSTALL_DIR/$KATANA_BINARY"
echo "🔧 Wake-up Config: $SUDOERS_FILE"
echo "🖥️  Desktop Launcher: $DESKTOP_FILE"
echo ""
echo "🚀 Usage:"
echo "   Run from terminal: katana"
echo "   Or search 'Katana' in your application menu"
echo ""
echo "⚡ Features Enabled:"
echo "   ✅ Time Tracker with analytics"
echo "   ✅ Stopwatch with lap times"
echo "   ✅ Countdown timer"
echo "   ✅ Smart alarms with 15+ sounds"
if sudo -n $RTCWAKE_PATH --version &> /dev/null 2>&1; then
    echo "   ✅ System wake-up from sleep (configured)"
else
    echo "   ⚠️  System wake-up (restart terminal to activate)"
fi
echo ""
echo "📚 Documentation:"
echo "   README.md - User guide"
echo "   WAKE_SETUP_GUIDE.md - Wake-up troubleshooting"
echo "   BUGFIXES.md - Known issues and fixes"
echo ""
echo "🔄 To uninstall:"
echo "   rm -f $INSTALL_DIR/$KATANA_BINARY"
echo "   sudo rm -f $SUDOERS_FILE"
echo "   rm -f $DESKTOP_FILE"
echo ""
echo "Happy timing! 🥷⏰"
echo ""
