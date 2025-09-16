#!/bin/bash

# Katana Multi-Timer Installation Script
# Supports Ubuntu, Debian, Arch Linux, Kali, Parrot OS, and other Linux distributions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        DISTRO=$ID
    else
        print_error "Cannot detect Linux distribution"
        exit 1
    fi
}

# Install dependencies based on distribution
install_dependencies() {
    print_status "Installing dependencies for $DISTRO..."
    
    case $DISTRO in
        ubuntu|debian|kali|parrot)
            sudo apt update
            sudo apt install -y golang-go git alsa-utils
            ;;
        arch|manjaro)
            sudo pacman -Sy --noconfirm go git alsa-utils
            ;;
        fedora|rhel|centos)
            sudo dnf install -y golang git alsa-utils
            ;;
        opensuse*)
            sudo zypper install -y go git alsa-utils
            ;;
        *)
            print_warning "Unknown distribution: $DISTRO"
            print_status "Please install: golang, git, and alsa-utils manually"
            read -p "Continue anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
            ;;
    esac
}

# Check if Go is installed and version
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        return 1
    fi
    
    GO_VERSION=$(go version | grep -oP 'go\d+\.\d+' | sed 's/go//')
    GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
    GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
    
    if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 19 ]); then
        print_warning "Go version $GO_VERSION detected. Recommended: 1.19+"
    else
        print_success "Go version $GO_VERSION is compatible"
    fi
}

# Install Katana
install_katana() {
    print_status "Downloading and building Katana..."
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd $TEMP_DIR
    
    # Clone repository
    git clone https://github.com/Sudo-Omar-Khalaf/katana.git
    cd katana
    
    # Build application
    print_status "Building Katana..."
    go mod download
    go build -o katana
    
    # Install to system
    print_status "Installing to /usr/local/bin..."
    sudo mv katana /usr/local/bin/
    sudo chmod +x /usr/local/bin/katana
    
    # Create application directory
    mkdir -p ~/.katana/sounds
    
    # Copy sound files
    if [ -d "assets/sounds" ]; then
        cp assets/sounds/* ~/.katana/sounds/ 2>/dev/null || true
    fi
    
    # Cleanup
    cd /
    rm -rf $TEMP_DIR
    
    print_success "Katana installed successfully!"
}

# Setup wake-up permissions (optional)
setup_wake_permissions() {
    print_status "Setting up system wake-up permissions..."
    
    if command -v rtcwake &> /dev/null; then
        # Create sudoers rule for rtcwake
        echo "$USER ALL=(ALL) NOPASSWD: /usr/sbin/rtcwake" | sudo tee /etc/sudoers.d/katana-wake > /dev/null
        print_success "Wake-up permissions configured"
    else
        print_warning "rtcwake not found. Wake-up functionality may be limited"
    fi
}

# Create desktop entry
create_desktop_entry() {
    print_status "Creating desktop entry..."
    
    cat > ~/.local/share/applications/katana.desktop << EOF
[Desktop Entry]
Name=Katana Multi-Timer
Comment=Alarm Clock, Stopwatch, Timer & Time Tracker
Exec=katana
Icon=alarm-clock
Terminal=false
Type=Application
Categories=Utility;Clock;Office;
Keywords=alarm;timer;stopwatch;clock;tracker;productivity;
EOF
    
    print_success "Desktop entry created"
}

# Main installation process
main() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════╗"
    echo "║        Katana Multi-Timer            ║"
    echo "║     Installation Script v1.0         ║"
    echo "╚══════════════════════════════════════╝"
    echo -e "${NC}"
    
    print_status "Starting Katana installation..."
    
    # Detect distribution
    detect_distro
    print_status "Detected distribution: $DISTRO"
    
    # Install dependencies
    install_dependencies
    
    # Check Go installation
    check_go
    
    # Install Katana
    install_katana
    
    # Setup wake permissions
    read -p "Setup system wake-up permissions? (requires sudo) [Y/n]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        setup_wake_permissions
    fi
    
    # Create desktop entry
    create_desktop_entry
    
    echo
    print_success "Installation completed successfully!"
    echo
    echo -e "${GREEN}You can now run Katana by typing: ${YELLOW}katana${NC}"
    echo -e "${GREEN}Or find it in your applications menu${NC}"
    echo
    print_status "For help and documentation, visit: https://github.com/Sudo-Omar-Khalaf/katana"
    echo
}

# Run main function
main "$@"
