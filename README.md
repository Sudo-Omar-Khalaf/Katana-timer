# ⚔️ Katana - The Ultimate Linux Time Management Suite

> **Advanced Alarm Clock with System Wake-Up | Stopwatch | Countdown Timer | Time Tracker**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Build](https://github.com/Sudo-Omar-Khalaf/Katana-timer/actions/workflows/build.yml/badge.svg)](https://github.com/Sudo-Omar-Khalaf/Katana-timer/actions)
[![Linux](https://img.shields.io/badge/Linux-Compatible-FCC624?style=flat&logo=linux&logoColor=black)](https://www.linux.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Wake Support](https://img.shields.io/badge/Wake-From--Sleep-success?style=flat)](https://github.com/Sudo-Omar-Khalaf/katana)
[![Production Ready](https://img.shields.io/badge/Status-Production%20Ready-brightgreen)](https://github.com/Sudo-Omar-Khalaf/katana)

**Katana** is a powerful, production-ready desktop application for Linux that combines an **intelligent alarm clock with system wake-up**, **stopwatch**, **countdown timer**, and **time tracker** into one beautiful terminal-styled interface. 

### 🎯 What Makes Katana Special?

- **🌙 True Wake-From-Sleep**: Unlike other alarms, Katana can actually wake your computer from sleep/suspend mode
- **🔐 Secure & Passwordless**: One-time setup, then no sudo passwords needed
- **🚀 Fully Automated**: Zero configuration after installation - just create an alarm and it works
- **💻 Native Linux**: Built specifically for Linux with hardware RTC wake support
- **🎨 Beautiful UI**: Terminal-inspired dark theme with a professional, developer-friendly interface

Perfect for Ubuntu, Debian, Arch Linux, Fedora, Kali Linux, Parrot OS, Manjaro, and all other Linux distributions.

## 📑 Table of Contents

- [Installation](#-installation)
- [Features](#-features)
- [Usage](#-usage)
- [Advanced Features](#-advanced-features--configuration)
- [System Requirements](#️-system-requirements)
- [Troubleshooting](#-troubleshooting)
- [Screenshots & Demo](#-screenshots--demo)
- [Use Cases](#-use-cases)
- [Contributing](#-contributing)
- [Roadmap](#️-roadmap)
- [License](#-license)
- [Support](#-support--community)

## 🚀 Installation

### Automated Installation (Recommended) - 30 Seconds

The automated installer handles everything: dependencies, rtcwake setup, sudo configuration, and PATH setup.

```bash
# Clone the repository
git clone https://github.com/Sudo-Omar-Khalaf/katana.git
cd katana

# Run the automated installer (requires sudo once for setup)
chmod +x install-katana.sh
./install-katana.sh
```

**What the installer does:**
- ✅ Detects your Linux distribution automatically
- ✅ Installs `rtcwake` if not present (for wake-from-sleep)
- ✅ Configures passwordless sudo for rtcwake (secure, minimal permissions)
- ✅ Installs Go build dependencies if needed
- ✅ Builds optimized binary
- ✅ Installs to `~/.local/bin` (no root required)
- ✅ Creates desktop launcher
- ✅ Adds to your PATH automatically

After installation:
```bash
# Just run katana from anywhere!
katana
```

### Manual Installation

If you prefer to install manually:

#### Prerequisites
```bash
# Ubuntu/Debian/Kali/Parrot/Pop!_OS/Mint
sudo apt update && sudo apt install -y golang-go git util-linux

# Arch Linux/Manjaro/EndeavourOS
sudo pacman -S go git util-linux

# Fedora/RHEL/CentOS/Rocky/Alma
sudo dnf install -y golang git util-linux

# openSUSE/SLES
sudo zypper install -y go git util-linux
```

#### Build and Install
```bash
# Clone and build
git clone https://github.com/Sudo-Omar-Khalaf/katana.git
cd katana
go build -o katana

# Install to user directory (no sudo needed)
mkdir -p ~/.local/bin
mv katana ~/.local/bin/

# Add to PATH if not already (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Run katana
katana
```

#### Optional: Wake-From-Sleep Setup (Highly Recommended)
```bash
# Configure passwordless rtcwake for wake-from-sleep functionality
sudo tee /etc/sudoers.d/katana-alarm-wake > /dev/null <<EOF
$USER ALL=(ALL) NOPASSWD: $(command -v rtcwake)
EOF

sudo chmod 0440 /etc/sudoers.d/katana-alarm-wake
```

## ⚡ Features

### 🔔 Intelligent Alarm Clock with Wake-From-Sleep

**The star feature** - Katana doesn't just ring when your computer is awake, it actually **wakes your computer up**!

- **🌙 Hardware RTC Wake**: Uses Real-Time Clock wake timers to wake your system from sleep/suspend
- **🔐 Passwordless Operation**: One-time setup, then zero sudo prompts
- **⚙️ Fully Automatic**: Wake-up is scheduled automatically when you create/enable an alarm
- **🎵 15+ Alarm Sounds**: Choose from professional alarm sounds or add your own
- **🔄 Smart Scheduling**: Handles same-day vs next-day alarms intelligently
- **📱 Desktop Notifications**: Visual and audio alerts when alarms trigger
- **✨ Zero Configuration**: Just create an alarm and it works - no manual offset settings

**How It Works:**
1. Create an alarm for 7:00 AM
2. Katana automatically programs your hardware RTC wake timer
3. Sleep/suspend your computer normally
4. Computer wakes at 7:00 AM (hardware level)
5. Alarm rings! ⏰

**Platform Support:**
- ✅ Linux: Full support via `rtcwake` (all distributions)
- 🔧 Windows: Task Scheduler integration (prepared)
- 🔧 macOS: `pmset` integration (prepared)

### ⏱️ Professional Stopwatch
- **High Precision**: Millisecond accuracy timing
- **Lap Recording**: Capture and track unlimited lap times
- **Lap Differences**: Automatic delta calculation between laps
- **Export Data**: Save timing sessions for analysis
- **Running Average**: Track average lap times in real-time

### ⏰ Countdown Timer
- **Visual Progress**: Beautiful progress bar with color-coded urgency
- **Custom Durations**: Set hours, minutes, and seconds precisely
- **Multiple Alerts**: Sound, notification, and visual alerts
- **Auto-Reset**: Quick restart for repeated timing sessions
- **Background Operation**: Continue working while timer runs

### 📊 Advanced Time Tracker
- **Activity Logging**: Track work sessions with descriptive names
- **Tag System**: Organize sessions with custom tags (#work, #coding, etc.)
- **Live Tracking**: Real-time session duration display
- **Analytics Dashboard**: Daily, weekly, and monthly productivity reports
- **Export Options**: Generate CSV and PDF reports
- **Session Management**: Pause, resume, edit, or delete sessions
- **Historical Data**: Complete session history with search

### 🎨 Beautiful Terminal-Style UI
- **Dark Theme**: Easy on the eyes with professional terminal green aesthetics
- **Monospace Fonts**: Developer-friendly interface with clear readability
- **Responsive Layout**: Adapts seamlessly to any screen size
- **Intuitive Navigation**: Tab-based interface for quick feature switching
- **Visual Feedback**: Progress bars, status indicators, and real-time updates

## 📱 Usage

### Alarm Clock
1. **Set an Alarm**: Enter name, time (HH:MM), and select sound
2. **System Wake-Up**: Your computer will automatically wake up when the alarm triggers
3. **Manage Alarms**: Enable/disable or delete alarms as needed

### Stopwatch
1. **Start/Stop**: Click to begin timing
2. **Capture Laps**: Record lap times during activities
3. **Reset**: Clear all times and start fresh

### Countdown Timer
1. **Set Duration**: Enter hours, minutes, seconds
2. **Start**: Begin countdown with visual progress
3. **Alerts**: Get notified when timer reaches zero

### Time Tracker
1. **Start Session**: Enter activity name and optional tags
2. **Track Time**: Monitor active session duration
3. **Analytics**: View daily/weekly/monthly reports
4. **Export**: Generate CSV/PDF reports

## 🔧 Advanced Features & Configuration

### System Wake-Up Details

**How Wake-From-Sleep Works:**

Katana uses your computer's hardware Real-Time Clock (RTC) to schedule wake events. This is the same mechanism your BIOS/UEFI uses for scheduled tasks.

**Technical Details:**
- Uses `rtcwake` command on Linux systems
- Programs hardware RTC wake alarm at exact alarm time
- Works even when system is fully suspended/hibernated
- Minimal power draw while asleep
- No background processes needed

**Security:**
- Sudo access is limited to `rtcwake` command only (not general sudo)
- Configuration file is read-only (0440 permissions)
- Easily reversible: `sudo rm /etc/sudoers.d/katana-alarm-wake`
- No daemon or background service required

**Supported Wake States:**
- ✅ Suspend (Sleep to RAM)
- ✅ Hibernate (Sleep to Disk)
- ✅ Hybrid Sleep
- ✅ Manual wake testing

### Custom Alarm Sounds

**Add Your Own Sounds:**
```bash
# Sounds should be in WAV format
cp your-custom-alarm.wav katana/assets/sounds/

# Rebuild if you want to embed them
go build -o katana
```

**Built-in Sounds:**
- Classic alarm tones
- Emergency alerts
- Gentle wake-up sounds
- Retro game sounds
- Natural sounds (rooster, etc.)
- And more!

### Data Export & Analytics

**Export Time Tracking Data:**
- **CSV Format**: Import into Excel, Google Sheets, or any spreadsheet software
- **PDF Reports**: Professional formatted reports with charts
- **Date Ranges**: Export specific time periods
- **Tag Filtering**: Export only sessions with specific tags
- **Monthly Summaries**: Automatic monthly productivity reports

### Desktop Integration

**Desktop Launcher:**
The installer creates a desktop launcher at:
```
~/.local/share/applications/katana-timer.desktop
```

**Features:**
- Appears in application menu
- Click to launch Katana
- Proper icon and categorization
- Quick access from desktop environment

### Command-Line Options

```bash
# Run Katana
katana

# Future options (planned):
katana --version          # Show version info
katana --help            # Show help
katana --start-timer 25  # Quick 25-minute timer
```

## 🛠️ System Requirements

### Minimum Requirements
- **OS**: Any Linux distribution with kernel 2.6+ (Ubuntu, Debian, Arch, Fedora, Kali, Parrot, Manjaro, etc.)
- **Go**: Version 1.19 or higher (for building from source)
- **Memory**: 50MB RAM during operation
- **Disk**: 30MB for installation
- **Display**: X11 or Wayland desktop environment
- **Audio**: ALSA or PulseAudio (standard on most Linux systems)

### For Wake-From-Sleep Feature
- **util-linux**: Package containing `rtcwake` (auto-installed by installer)
- **RTC Support**: Hardware Real-Time Clock (present in 99.9% of computers)
- **Sudo Access**: One-time setup for passwordless `rtcwake` execution

### Tested Distributions
✅ Ubuntu 18.04+ / Pop!_OS / Linux Mint  
✅ Debian 10+ / Kali Linux / Parrot OS  
✅ Arch Linux / Manjaro / EndeavourOS  
✅ Fedora 30+ / RHEL 8+ / CentOS Stream  
✅ openSUSE Leap / Tumbleweed  

### Architecture Support
- ✅ x86_64 (AMD64)
- ✅ ARM64 (aarch64)
- ✅ ARM (armv7l)
- ✅ i386 (32-bit)

## 🔍 Troubleshooting

### Common Issues

#### Audio Not Playing

**Problem**: Alarm sounds don't play

**Solutions:**
```bash
# Install audio dependencies
sudo apt install -y alsa-utils pulseaudio-utils  # Ubuntu/Debian
sudo pacman -S alsa-utils pulseaudio             # Arch Linux
sudo dnf install -y alsa-utils pulseaudio        # Fedora

# Test audio system
aplay /usr/share/sounds/alsa/Front_Center.wav

# Check PulseAudio status
pulseaudio --check -v
```

#### Wake-Up Not Working

**Problem**: Computer doesn't wake up for alarms

**Check RTC Support:**
```bash
# Test if rtcwake is installed and working
sudo rtcwake -m no -s 10  # Should schedule a wake in 10 seconds

# Check RTC device
ls -l /dev/rtc*

# Test wake-up scheduling
cat /sys/class/rtc/rtc0/wakealarm
```

**Verify Sudo Permissions:**
```bash
# Test passwordless rtcwake
sudo -n rtcwake --version

# If it asks for password, reconfigure:
./install-katana.sh  # Re-run installer
```

**Common Causes:**
- BIOS/UEFI wake-on-RTC disabled (check BIOS settings)
- Secure Boot restrictions (rare)
- Virtual machine without RTC passthrough

#### Build Errors

**Problem**: `go build` fails

**Solutions:**
```bash
# Update Go modules
go mod download
go mod tidy

# Clear cache and rebuild
go clean -cache -modcache
go build -o katana

# Verify Go version (needs 1.19+, recommended 1.23+)
go version

# If you see "invalid go version" error in go.mod:
# The go.mod file should use format "go 1.23" not "go 1.23.5"
# Fix with: sed -i 's/go 1\.[0-9]*\.[0-9]*/go 1.23/' go.mod
```

#### Binary Not Found After Install

**Problem**: `katana: command not found`

**Solutions:**
```bash
# Check if installed
ls -lh ~/.local/bin/katana

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Reload shell configuration
source ~/.bashrc  # or source ~/.zshrc

# Verify
which katana
```

#### Alarm Doesn't Ring at Set Time

**Problem**: Alarm created but doesn't trigger

**Checks:**
```bash
# Verify alarm is enabled (toggle switch should be green)
# Check system time is correct
timedatectl status

# View Katana logs (if running from terminal)
katana  # Check terminal output for errors
```

### Getting Help

If you encounter issues not covered here:

1. **Check Logs**: Run `katana` from terminal to see detailed logs
2. **Test Wake Functionality**: Use `test-wake-functionality.sh` script
3. **GitHub Issues**: [Report a bug](https://github.com/Sudo-Omar-Khalaf/katana/issues)
4. **Documentation**: See `IMPLEMENTATION_ANALYSIS.md` for technical details

## 🎬 Screenshots & Demo

### Main Interface
```
┌─────────────────────────────────────────────────────┐
│  ⚔️  KATANA TIMER                    ⏰ 13:37:42    │
├─────────────────────────────────────────────────────┤
│  [Alarm]  [Stopwatch]  [Timer]  [Tracker]          │
│                                                      │
│  🔔 Active Alarms                                   │
│  ┌──────────────────────────────────────────────┐  │
│  │ ⏰ Morning Alarm      07:00    [ON]  [Edit]  │  │
│  │ 🎵 Classic Alarm                      [Del]   │  │
│  ├──────────────────────────────────────────────┤  │
│  │ ⏰ Meeting Reminder   14:30   [OFF]  [Edit]  │  │
│  │ 🎵 Emergency Alert                    [Del]   │  │
│  └──────────────────────────────────────────────┘  │
│                                                      │
│  [+ Create New Alarm]                               │
│                                                      │
│  💡 System wake-up is enabled                       │
│  ✅ Computer will wake from sleep for alarms        │
└─────────────────────────────────────────────────────┘
```

### Time Tracker Dashboard
```
┌─────────────────────────────────────────────────────┐
│  📊 Time Tracker                                     │
├─────────────────────────────────────────────────────┤
│  🟢 Active: Coding Session          [02:34:15]      │
│  Tags: #work #golang #katana                        │
│                                                      │
│  [Pause]  [Stop]  [Add Note]                        │
│                                                      │
│  Today's Summary                                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 6h 15m              │
│  • Coding           4h 30m   [████████░░] 72%       │
│  • Meetings         1h 15m   [███░░░░░░░] 20%       │
│  • Documentation      30m   [█░░░░░░░░░]  8%        │
│                                                      │
│  [Export CSV]  [Generate PDF Report]                │
└─────────────────────────────────────────────────────┘
```

## 🎯 Use Cases

### For Developers
- ⏰ Wake up for early morning stand-ups (even if laptop was suspended)
- ⏱️ Time your coding sessions and track productivity
- 📊 Generate weekly reports for client billing
- ⏰ Set reminders for deployment windows

### For Students
- ⏰ Never miss early morning classes (reliable wake-up even from sleep)
- ⏱️ Time your study sessions with the Pomodoro technique
- 📊 Track time spent on different subjects
- ⏰ Set exam reminders

### For System Administrators
- ⏰ Schedule maintenance window reminders
- ⏱️ Time server operations and deployments
- 📊 Log time for client projects
- ⏰ Alert for monitoring events

### For Freelancers
- 📊 Track billable hours accurately
- ⏰ Never miss client meetings (reliable wake-up)
- ⏱️ Time different project tasks
- 📄 Export professional time reports for invoicing

## 🤝 Contributing

We welcome contributions! Katana is open source and community-driven.

### How to Contribute

1. **Fork** the repository on GitHub
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/katana.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes and test thoroughly
5. **Commit** your changes: `git commit -m 'Add amazing feature'`
6. **Push** to your fork: `git push origin feature/amazing-feature`
7. **Submit** a Pull Request with a clear description

### Development Setup
```bash
# Clone the repository
git clone https://github.com/Sudo-Omar-Khalaf/katana.git
cd katana

# Install dependencies
go mod download

# Run in development mode
go run main.go

# Build for testing
go build -o katana

# Run tests (when available)
go test ./...
```

### Contribution Ideas

**Features:**
- [ ] Windows support for wake-from-sleep
- [ ] macOS support for wake-from-sleep
- [ ] Recurring alarms (weekly, daily patterns)
- [ ] Snooze functionality
- [ ] More alarm sound options
- [ ] Themes and customization
- [ ] Mobile companion app
- [ ] Cloud sync for settings

**Code Quality:**
- [ ] Unit tests
- [ ] Integration tests
- [ ] Performance optimizations
- [ ] Code documentation
- [ ] Refactoring opportunities

**Documentation:**
- [ ] Video tutorials
- [ ] Usage examples
- [ ] Translation to other languages
- [ ] Wiki articles

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add comments for complex logic
- Keep functions focused and small
- Write meaningful commit messages

### Reporting Bugs

Found a bug? Please [open an issue](https://github.com/Sudo-Omar-Khalaf/katana/issues) with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- System information (OS, distribution, kernel version)
- Error logs if available

## 🗺️ Roadmap

### Version 2.0 (Planned)
- [ ] Recurring alarm patterns (daily, weekly, weekdays)
- [ ] Snooze functionality with custom intervals
- [ ] Gradual volume increase for gentle wake-up
- [ ] Weather-based alarm adjustments
- [ ] Calendar integration

### Version 2.5 (Future)
- [ ] Windows wake-from-sleep support
- [ ] macOS wake-from-sleep support
- [ ] Mobile companion app (Android/iOS)
- [ ] Cloud sync for settings and data
- [ ] Team time tracking features

### Community Requested
- [ ] Themes and color customization
- [ ] Keyboard shortcuts configuration
- [ ] Multiple language support
- [ ] Plugin system for extensions
- [ ] REST API for automation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**TL;DR**: Free to use, modify, distribute, and sell. No warranty provided.

## 🌟 Star History

If you find Katana useful, please consider giving it a star! ⭐

Stars help others discover this project and motivate continued development.

[![Star History Chart](https://api.star-history.com/svg?repos=Sudo-Omar-Khalaf/katana&type=Date)](https://star-history.com/#Sudo-Omar-Khalaf/katana&Date)

## 📞 Support & Community

- **🐛 Issues**: [GitHub Issues](https://github.com/Sudo-Omar-Khalaf/katana/issues) - Report bugs or request features
- **💬 Discussions**: [GitHub Discussions](https://github.com/Sudo-Omar-Khalaf/katana/discussions) - Ask questions, share ideas
- **📚 Documentation**: [Wiki](https://github.com/Sudo-Omar-Khalaf/katana/wiki) - Comprehensive guides
- **📧 Email**: For private inquiries or security issues

## 🏆 Acknowledgments

Special thanks to:
- The Go community for excellent tooling
- Fyne.io for the beautiful UI framework
- Linux kernel developers for RTC wake support
- All contributors and users who provide feedback

## 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/Sudo-Omar-Khalaf/katana?style=social)
![GitHub forks](https://img.shields.io/github/forks/Sudo-Omar-Khalaf/katana?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/Sudo-Omar-Khalaf/katana?style=social)
![GitHub issues](https://img.shields.io/github/issues/Sudo-Omar-Khalaf/katana)
![GitHub pull requests](https://img.shields.io/github/issues-pr/Sudo-Omar-Khalaf/katana)
![GitHub last commit](https://img.shields.io/github/last-commit/Sudo-Omar-Khalaf/katana)

## 🔗 Related Projects

- **Fyne** - Cross-platform GUI toolkit for Go
- **util-linux** - Essential Linux utilities including rtcwake
- **ALSA** - Advanced Linux Sound Architecture

---

<div align="center">

### Made with ❤️ for the Linux Community

**Perfect for developers, system administrators, students, and anyone who needs reliable time management on Linux**

*Because your computer should wake up when you need it to, not just when it feels like it.* 😴 ➡️ ⏰

[⭐ Star](https://github.com/Sudo-Omar-Khalaf/katana) · [🐛 Report Bug](https://github.com/Sudo-Omar-Khalaf/katana/issues) · [✨ Request Feature](https://github.com/Sudo-Omar-Khalaf/katana/issues) · [💬 Discuss](https://github.com/Sudo-Omar-Khalaf/katana/discussions)

</div>
