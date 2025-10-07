# 📁 Katana Project Structure

**Last Updated**: October 7, 2025  
**Version**: 1.0.0

---

## 📂 Directory Organization

The Katana project follows a clean, organized structure for better maintainability and ease of contribution.

### Root Directory
```
katana/
├── README.md              # Main project documentation
├── LICENSE                # MIT License
├── CHANGELOG.md           # Version history and changes
├── CONTRIBUTING.md        # Contribution guidelines
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── main.go                # Application entry point
├── config.go              # Configuration management
├── katana                 # Compiled binary (not in git)
├── katana.code-workspace  # VS Code workspace settings
└── install.sh             # Symlink to scripts/install-katana.sh
```

### Source Code (`/`)
```
├── power/                 # Power management and wake-up scheduling
│   └── power.go          # RTC wake implementation
├── ui/                    # User interface
│   └── mainui.go         # Main UI with Fyne framework
├── sound/                 # Audio playback
│   └── player.go         # Sound player implementation
├── tracker/               # Time tracking functionality
│   └── session.go        # Session management
├── storage/               # Data persistence
│   └── storage.go        # SQLite database operations
└── export/                # Data export functionality
    └── export.go         # CSV and PDF export
```

### Assets (`assets/`)
```
assets/
└── sounds/                # Built-in alarm sounds (15 wav files)
    ├── mixkit-classic-alarm-995.wav
    ├── mixkit-emergency-alert-alarm-1007.wav
    ├── mixkit-rooster-crowing-in-the-morning-2462.wav
    └── ... (12 more alarm sounds)
```

### Documentation (`docs/`)
```
docs/
├── IMPLEMENTATION_ANALYSIS.md    # Complete technical analysis
├── COMPLETE_SOLUTION.md          # Technical solution details
├── FINAL_SUMMARY.md              # Quick reference summary
├── DEMO.md                       # Demo and usage examples
│
├── guides/                       # User guides
│   ├── INSTALLATION_GUIDE.md    # Installation instructions
│   ├── WAKE_SETUP_GUIDE.md      # Wake-from-sleep setup
│   └── QUICK_REFERENCE.txt      # Command reference
│
├── development/                  # Development documentation
│   ├── BUGFIXES.md              # Bug fix history
│   ├── DEVELOPMENT.md           # Development notes
│   ├── PROJECT_STATUS.md        # Project status tracking
│   └── GITHUB_PUSH_GUIDE.md     # Git workflow guide
│
└── images/                       # Screenshots and images
    └── image.png                # Application screenshot
```

### Scripts (`scripts/`)
```
scripts/
├── install-katana.sh            # Automated installation script
├── setup-wake-permissions.sh    # Wake permission setup
└── test-wake-functionality.sh   # Wake functionality tests
```

### Data (`data/`)
```
data/
└── sessions.db                  # SQLite database (time tracking data)
```

---

## 🎯 Key Files

### For Users
- **README.md** - Start here! Complete project overview and usage
- **docs/guides/INSTALLATION_GUIDE.md** - Detailed installation steps
- **docs/guides/WAKE_SETUP_GUIDE.md** - Setting up wake-from-sleep
- **scripts/install-katana.sh** - Automated installer
- **install.sh** - Convenient symlink to installer

### For Developers
- **main.go** - Application entry point
- **docs/IMPLEMENTATION_ANALYSIS.md** - Technical deep dive
- **docs/development/** - Development notes and project tracking
- **CONTRIBUTING.md** - How to contribute

### For Installation
- **scripts/install-katana.sh** - Main installation script
  - Auto-detects Linux distribution
  - Installs dependencies (rtcwake, Go)
  - Configures sudo permissions
  - Builds and installs binary
  - Sets up PATH and desktop launcher

---

## 📦 Build Artifacts (Not in Git)

These files are generated during build/runtime and excluded via `.gitignore`:

```
katana                    # Compiled binary (25MB)
katana_export.csv         # Exported time tracking data
katana_export.pdf         # PDF reports
*.test                    # Test binaries
*.out                     # Coverage files
```

---

## 🔧 Configuration Files

### Go Configuration
- **go.mod** - Module dependencies (Fyne, SQLite, etc.)
- **go.sum** - Dependency checksums for verification

### IDE Configuration
- **katana.code-workspace** - VS Code workspace settings
- **.vscode/** - VS Code specific settings (gitignored)

### Version Control
- **.gitignore** - Files excluded from git
  - Binaries
  - Export files
  - IDE files
  - OS generated files
  - Build artifacts

---

## 🚀 Quick Navigation

### Installation
```bash
# Run the installer
./install.sh
# OR
./scripts/install-katana.sh
```

### Building from Source
```bash
# Install dependencies
go mod download

# Build binary
go build -o katana

# Run without building
go run main.go
```

### Testing
```bash
# Test wake functionality
./scripts/test-wake-functionality.sh
```

### Documentation Paths
```bash
# User guides
cat docs/guides/QUICK_REFERENCE.txt
open docs/guides/INSTALLATION_GUIDE.md

# Technical docs
open docs/IMPLEMENTATION_ANALYSIS.md

# Development notes
ls docs/development/
```

---

## 📊 Project Statistics

- **Total Directories**: 15
- **Total Files**: 48+
- **Source Code**: ~2,000 lines of Go
- **Documentation**: 14 markdown files
- **Scripts**: 3 shell scripts
- **Alarm Sounds**: 15 WAV files

---

## 🎨 Design Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Clean Root**: Root directory contains only essential files
3. **Organized Docs**: Documentation grouped by audience (user/developer)
4. **Easy Navigation**: Clear folder names and structure
5. **Consistent Naming**: Descriptive, kebab-case filenames

---

## 🔄 Recent Changes

**October 7, 2025 - Major Reorganization**
- Created `docs/` folder structure
- Moved all documentation to appropriate subdirectories
- Created `scripts/` folder for installation scripts
- Removed redundant and backup files
- Improved code formatting and imports
- Created symlink for convenient installation access

---

## 📝 Future Structure Enhancements

Planned additions to project structure:

- [ ] `tests/` - Unit and integration tests
- [ ] `cmd/` - Multiple command-line tools
- [ ] `internal/` - Private application packages
- [ ] `pkg/` - Public library packages
- [ ] `api/` - REST API (if implemented)
- [ ] `web/` - Web interface (if implemented)
- [ ] `docs/api/` - API documentation
- [ ] `examples/` - Usage examples

---

**Maintained by**: Katana Development Team  
**License**: MIT  
**Repository**: https://github.com/Sudo-Omar-Khalaf/Katana-timer
