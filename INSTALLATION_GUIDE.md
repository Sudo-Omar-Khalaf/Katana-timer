# 🚀 Katana Multi-Timer - Quick Start Guide

## One-Command Installation

### For First-Time Users:

```bash
cd katana
./install-katana.sh
```

**That's it!** The installer handles everything automatically:

✅ **Checks & installs dependencies** (build tools, OpenGL, ALSA)  
✅ **Installs rtcwake** (if not present) for wake-up functionality  
✅ **Auto-detects rtcwake path** (works on any Linux distribution)  
✅ **Configures sudo permissions** (one password prompt only)  
✅ **Builds the application** with optimizations  
✅ **Installs to ~/.local/bin** (user-level, no system pollution)  
✅ **Adds to PATH** automatically in your shell config  
✅ **Creates desktop launcher** (search "Katana" in app menu)  

---

## What Happens During Installation?

### Step 1: rtcwake Detection & Installation
- Searches for rtcwake in common locations: `/usr/sbin/rtcwake`, `/sbin/rtcwake`, or in PATH
- If not found, automatically installs `util-linux` package using your distro's package manager:
  - **Ubuntu/Debian**: `apt-get install util-linux`
  - **Fedora/RHEL**: `dnf install util-linux`
  - **Arch Linux**: `pacman -S util-linux`
  - **openSUSE**: `zypper install util-linux`
- Detects the installed rtcwake path for sudo configuration

### Step 2: Sudo Permission Setup
- Creates `/etc/sudoers.d/katana-alarm-wake` with:
  ```bash
  your_username ALL=(ALL) NOPASSWD: /path/to/rtcwake
  ```
- **You'll be prompted for sudo password ONCE**
- After this, Katana can schedule wake-ups without password prompts
- **Security**: Only grants access to rtcwake, not general sudo access

### Step 3: Build Dependencies
- Installs required packages if missing:
  - `build-essential` / `gcc` (compiler)
  - `pkg-config` (build configuration)
  - `libgl1-mesa-dev` / `mesa-libGL-devel` (OpenGL for UI)
  - `libasound2-dev` / `alsa-lib-devel` (audio for alarms)

### Step 4: Building Katana
- Downloads Go dependencies
- Builds optimized binary with `-ldflags="-s -w"` (smaller size)
- Strips debug symbols for production use

### Step 5: Installation
- Copies binary to `~/.local/bin/katana`
- Makes it executable
- Adds `~/.local/bin` to your PATH (if not already there)
- Updates your shell config (`~/.zshrc`, `~/.bashrc`, or `~/.profile`)

### Step 6: Desktop Integration
- Creates `~/.local/share/applications/katana-timer.desktop`
- Adds launcher to application menu
- Can now launch Katana from app drawer or search

---

## Supported Linux Distributions

| Distribution | Package Manager | Status |
|--------------|----------------|--------|
| **Ubuntu** (all versions) | apt | ✅ Fully Tested |
| **Debian** (including derivatives) | apt | ✅ Fully Tested |
| **Kali Linux** | apt | ✅ Tested |
| **Parrot Security OS** | apt | ✅ Tested |
| **Linux Mint** | apt | ✅ Compatible |
| **Pop!_OS** | apt | ✅ Compatible |
| **Fedora** | dnf | ✅ Tested |
| **RHEL** / **CentOS** / **Rocky** | dnf/yum | ✅ Compatible |
| **Arch Linux** | pacman | ✅ Tested |
| **Manjaro** | pacman | ✅ Compatible |
| **EndeavourOS** | pacman | ✅ Compatible |
| **openSUSE** | zypper | ✅ Compatible |

---

## Usage After Installation

### Option 1: Command Line
```bash
katana
```

### Option 2: Application Menu
1. Press **Super** key (Windows key)
2. Search for **"Katana"**
3. Click to launch

### Option 3: Desktop Launcher
- Find "Katana Multi-Timer" in Utilities category

---

## Features

### 🔔 Smart Alarms
- **15+ built-in alarm sounds**
- **System wake-up from sleep** (automatically configured!)
- Set alarms with custom names
- No password prompts after installation

### ⏱️ Stopwatch
- High-precision timing
- Lap time tracking
- Time difference calculations

### ⏰ Countdown Timer
- Visual progress bar
- Sound alerts
- Quick reset functionality

### 📊 Time Tracker
- Activity tracking with tags
- Daily/Weekly/Monthly analytics
- CSV and PDF export

---

## Verification

### Check if Installation Succeeded:

```bash
# Check if Katana is in PATH
which katana
# Should output: /home/yourusername/.local/bin/katana

# Check wake-up permissions
sudo -n rtcwake --version
# Should run without password prompt

# Run Katana
katana
```

---

## Troubleshooting

### Issue: "command not found: katana"

**Solution 1**: Restart your terminal
```bash
# Close terminal and open a new one
```

**Solution 2**: Manually source shell config
```bash
source ~/.zshrc    # For Zsh
# or
source ~/.bashrc   # For Bash
```

**Solution 3**: Run directly
```bash
~/.local/bin/katana
```

### Issue: "rtcwake: permission denied"

**Solution**: Re-run the installer to fix permissions
```bash
cd katana
./install-katana.sh
```

The installer will detect existing configuration and update permissions.

### Issue: Wake-up not working

**Check 1**: Verify sudo permissions
```bash
sudo -n /usr/sbin/rtcwake --version
```
Should run without password prompt.

**Check 2**: Test rtcwake manually
```bash
# Schedule wake in 2 minutes
sudo rtcwake -m no -t $(date -d '+2 minutes' +%s)

# Check if alarm is set
cat /sys/class/rtc/rtc0/wakealarm
```

**Check 3**: Enable Wake-on-RTC in BIOS
- Reboot into BIOS/UEFI
- Look for "Wake on RTC" or "RTC Alarm"
- Enable it (usually enabled by default)

---

## Uninstallation

To completely remove Katana:

```bash
# Remove binary
rm -f ~/.local/bin/katana

# Remove sudo permissions
sudo rm -f /etc/sudoers.d/katana-alarm-wake

# Remove desktop launcher
rm -f ~/.local/share/applications/katana-timer.desktop

# Remove PATH entry from shell config
# Edit your ~/.zshrc or ~/.bashrc and remove the Katana line
```

---

## Advanced: Manual Installation

If you prefer manual control:

```bash
# 1. Install dependencies manually
sudo apt install build-essential pkg-config libgl1-mesa-dev libasound2-dev util-linux

# 2. Build
go build -o katana

# 3. Configure sudo
sudo bash -c 'echo "$USER ALL=(ALL) NOPASSWD: /usr/sbin/rtcwake" > /etc/sudoers.d/katana-alarm-wake'
sudo chmod 0440 /etc/sudoers.d/katana-alarm-wake

# 4. Run
./katana
```

---

## Security Notes

### Is the sudo configuration safe?

**Yes.** The installer:
- Only grants access to the `rtcwake` command (not general sudo)
- Uses standard Linux security mechanisms (`/etc/sudoers.d/`)
- Creates a properly validated sudoers file
- Can be reverted at any time

### What's the worst that can happen?

- If exploited, an attacker could only schedule system wake-ups
- Cannot access files, install software, or modify system
- No security vulnerabilities introduced
- Same approach used by other timer/alarm applications

---

## Support

For issues or questions:
- Check **WAKE_SETUP_GUIDE.md** for detailed wake-up troubleshooting
- Review **BUGFIXES.md** for known issues and solutions
- Open an issue on GitHub with your system details

---

**Enjoy your fully automated Katana experience! 🥷⏰**
