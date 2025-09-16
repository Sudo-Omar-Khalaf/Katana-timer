# Katana Multi-Timer - Linux Alarm Clock & Time Tracker

> **The Ultimate Linux Desktop Alarm Clock, Stopwatch, Countdown Timer & Time Tracker**

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Linux](https://img.shields.io/badge/Linux-Compatible-FCC624?style=flat&logo=linux&logoColor=black)](https://www.linux.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Katana** is a powerful, lightweight desktop application for Linux that combines an **alarm clock**, **stopwatch**, **countdown timer**, and **time tracker** into one beautiful terminal-styled interface. Perfect for Ubuntu, Debian, Arch Linux, Kali Linux, Parrot OS, and all other Linux distributions.

## 🚀 Quick Installation (30 seconds)

### One-Line Install (Recommended)
```bash
# Download, build, and install in one command
curl -sSL https://raw.githubusercontent.com/Sudo-Omar-Khalaf/katana/main/install.sh | bash
```

### Manual Installation

#### Prerequisites
```bash
# Ubuntu/Debian/Kali/Parrot
sudo apt update && sudo apt install -y golang-go git

# Arch Linux/Manjaro
sudo pacman -S go git

# Fedora/RHEL/CentOS
sudo dnf install -y golang git

# OpenSUSE
sudo zypper install -y go git
```

#### Install Katana
```bash
# Clone and build
git clone https://github.com/Sudo-Omar-Khalaf/katana.git
cd katana
go build -o katana
sudo mv katana /usr/local/bin/

# Run from anywhere
katana
```

## ⚡ Features

### 🔔 Smart Alarm Clock
- **System Wake-Up**: Automatically wakes your computer from sleep/suspend
- **Custom Sounds**: Choose from 15+ built-in alarm sounds
- **One-time & Recurring**: Set alarms for specific times or daily schedules
- **Visual Notifications**: Desktop notifications with sound alerts

### ⏱️ Professional Stopwatch
- **High Precision**: Millisecond accuracy timing
- **Lap Times**: Capture and track multiple lap times
- **Time Differences**: Automatic calculation between laps
- **Export Data**: Save timing data for analysis

### ⏰ Countdown Timer
- **Visual Progress**: Progress bar with color-coded urgency
- **Custom Durations**: Set hours, minutes, and seconds
- **Alert System**: Multiple notification methods when timer expires
- **Auto-Reset**: Quick restart functionality

### 📊 Time Tracker
- **Activity Logging**: Track work sessions with custom activities
- **Tag System**: Organize sessions with custom tags
- **Analytics**: Daily, weekly, and monthly time reports
- **Export Options**: CSV and PDF reports for productivity analysis
- **Session Management**: Pause, resume, and manage active sessions

### 🎨 Terminal-Style UI
- **Dark Theme**: Easy on the eyes with terminal green aesthetics
- **Monospace Fonts**: Professional developer-friendly interface
- **Responsive Layout**: Adapts to any screen size
- **Keyboard Shortcuts**: Efficient navigation and control

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

## 🔧 Advanced Features

### System Wake-Up (Linux Power Management)
Katana can wake your computer from sleep/suspend mode when alarms are scheduled:

**Supported Methods:**
- **Linux**: `rtcwake` (requires sudo) or `at` command
- **Ubuntu/Debian**: Native systemd integration
- **Arch Linux**: Hardware RTC wake support

**Setup Wake-Up (Optional):**
```bash
# Grant sudo permissions for wake functionality (one-time setup)
echo "$USER ALL=(ALL) NOPASSWD: /usr/sbin/rtcwake" | sudo tee /etc/sudoers.d/katana-wake
```

### Custom Alarm Sounds
Add your own alarm sounds:
```bash
# Copy .wav files to the sounds directory
cp your-alarm.wav ~/.katana/sounds/
```

### Data Export
Export your time tracking data:
- **CSV Format**: For spreadsheet analysis
- **PDF Reports**: Professional formatted reports
- **Monthly Summaries**: Comprehensive productivity insights

## 🛠️ System Requirements

- **OS**: Any Linux distribution (Ubuntu, Debian, Arch, Kali, Parrot, etc.)
- **Go**: Version 1.19 or higher
- **Dependencies**: Automatically handled by Go modules
- **Audio**: ALSA or PulseAudio (standard on most Linux systems)
- **Display**: X11 or Wayland desktop environment

## 🔍 Troubleshooting

### Common Issues

**Audio not working:**
```bash
# Install audio dependencies
sudo apt install -y alsa-utils pulseaudio-utils  # Ubuntu/Debian
sudo pacman -S alsa-utils pulseaudio             # Arch Linux
```

**Wake-up not working:**
```bash
# Test rtcwake functionality
sudo rtcwake -m no -s 10  # Should schedule a wake in 10 seconds
```

**Build errors:**
```bash
# Update Go and dependencies
go mod tidy
go clean -cache
```

## 🎯 SEO Keywords
Linux alarm clock, Ubuntu alarm, Debian timer, Arch Linux stopwatch, Kali Linux alarm, Parrot OS timer, Linux time tracker, desktop alarm Linux, system wake alarm, Linux productivity timer, terminal alarm clock, Go alarm application, Linux desktop timer, Ubuntu wake-up alarm, open source alarm Linux

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/Sudo-Omar-Khalaf/katana.git`
3. **Create** a feature branch: `git checkout -b feature-name`
4. **Make** your changes and test thoroughly
5. **Submit** a pull request

### Development Setup
```bash
git clone https://github.com/Sudo-Omar-Khalaf/katana.git
cd katana
go mod download
go run main.go
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Star History

If you find Katana useful, please consider giving it a star! ⭐

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/Sudo-Omar-Khalaf/katana/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Sudo-Omar-Khalaf/katana/discussions)
- **Documentation**: [Wiki](https://github.com/Sudo-Omar-Khalaf/katana/wiki)

---

**Made with ❤️ for the Linux community**

*Perfect for developers, system administrators, and anyone who needs reliable time management on Linux systems.*
