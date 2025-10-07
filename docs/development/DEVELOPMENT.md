# Katana Time Tracker - Development Guide

## Project Structure

```
katana/
├── main.go              # Application entry point
├── config.go            # Configuration management
├── README.md            # User documentation
├── CHANGELOG.md         # Version history
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── requirements.txt    # System dependencies
├── .gitignore          # Git ignore rules
├── image.png           # Application screenshot
├── data/               # Data storage directory
│   ├── sessions.db     # SQLite database
│   ├── sessions.json   # JSON fallback
│   └── config.json     # Application configuration
├── export/             # Export functionality
│   └── export.go       # CSV, JSON, PDF exporters
├── storage/            # Data persistence layer
│   └── storage.go      # SQLite with JSON fallback
├── tracker/            # Core session tracking
│   └── session.go      # Session data structures
└── ui/                 # User interface
    └── mainui.go       # Fyne-based terminal-style UI
```

## Architecture

### Core Components

1. **Main Application** (`main.go`)
   - Application initialization
   - Window setup and configuration
   - Graceful shutdown handling
   - Signal handling (SIGTERM, SIGINT)

2. **Configuration** (`config.go`)
   - JSON-based configuration system
   - Default values and validation
   - Automatic config file creation

3. **Session Tracking** (`tracker/session.go`)
   - Session data structure and validation
   - Activity and tag parsing
   - Duration calculation and formatting

4. **Data Storage** (`storage/storage.go`)
   - SQLite primary storage with JSON fallback
   - CRUD operations for sessions
   - Data integrity and error handling

5. **User Interface** (`ui/mainui.go`)
   - Terminal-style green-on-black theme
   - Custom widgets (TerminalButton, TerminalTab)
   - Real-time timer updates
   - Activity list with filtering
   - Visual analytics (daily/weekly/monthly grids)

6. **Export System** (`export/export.go`)
   - Multi-format export (CSV, JSON, PDF)
   - File save dialogs
   - Formatted output generation

## Development Setup

### Prerequisites

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y build-essential pkg-config libgl1-mesa-dev git golang

# Fedora/RHEL
sudo dnf install -y gcc pkg-config mesa-libGL-devel git golang

# Arch Linux
sudo pacman -S base-devel pkg-config mesa git go
```

### Getting Started

1. **Clone and Build**
   ```bash
   git clone <repository-url>
   cd katana
   go mod tidy
   go build -o katana
   ```

2. **Development Mode**
   ```bash
   go run main.go
   ```

3. **Testing**
   ```bash
   go test ./...
   go vet ./...
   ```

## Code Style Guidelines

### Go Code Standards

1. **Follow Go conventions**
   - Use `gofmt` for formatting
   - Follow effective Go practices
   - Use meaningful variable and function names

2. **Error Handling**
   - Always handle errors explicitly
   - Use proper error wrapping
   - Provide user-friendly error messages

3. **Documentation**
   - Document all exported functions and types
   - Use Go doc comments
   - Keep comments concise and helpful

### UI Development

1. **Custom Widgets**
   - Extend `widget.BaseWidget` for custom components
   - Implement proper renderer patterns
   - Handle mouse events and focus states

2. **Theme Consistency**
   - Use terminal green (`#00FF00`) for primary colors
   - Black backgrounds for terminal aesthetic
   - Monospace fonts for consistency

3. **Resource Management**
   - Properly dispose of resources
   - Use background goroutines carefully
   - Implement cleanup methods

## Adding Features

### New Export Format

1. Add function to `export/export.go`:
   ```go
   func ExportToXML(sessions []*tracker.Session, filename string) error {
       // Implementation
   }
   ```

2. Add UI button in `ui/mainui.go`:
   ```go
   exportXML := NewTerminalButton("Export XML", func() {
       // File dialog and export logic
   })
   ```

### New Storage Backend

1. Implement interface in `storage/storage.go`
2. Add initialization logic
3. Maintain backward compatibility

### UI Enhancements

1. Create new custom widget extending `widget.BaseWidget`
2. Implement `CreateRenderer()` method
3. Handle events (mouse, keyboard, focus)
4. Update main UI layout

## Testing

### Manual Testing Checklist

- [ ] Start/stop tracking works
- [ ] Sessions save correctly
- [ ] Activity list updates
- [ ] Tag filtering works
- [ ] Export functions work
- [ ] Application closes gracefully
- [ ] Data persists across restarts
- [ ] Error dialogs appear for invalid input

### Performance Testing

- Monitor memory usage during long sessions
- Test with large numbers of sessions
- Verify UI responsiveness
- Check database performance

## Building for Release

### Linux Binary

```bash
go build -ldflags="-s -w" -o katana
```

### Cross-compilation

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o katana.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o katana-mac
```

## Troubleshooting

### Common Issues

1. **CGO Compilation Errors**
   - Ensure C compiler is installed
   - Check pkg-config is available
   - Verify OpenGL libraries are present

2. **SQLite Issues**
   - Check file permissions in data/ directory
   - Verify SQLite driver compilation
   - Use JSON fallback if needed

3. **UI Problems**
   - Update graphics drivers
   - Check X11/Wayland compatibility
   - Verify Fyne dependencies

### Debugging

1. **Enable Logging**
   ```go
   log.SetOutput(os.Stderr) // Enable debug output
   ```

2. **Profile Performance**
   ```bash
   go build -race # Race condition detection
   go tool pprof # Performance profiling
   ```

## Contributing

1. Fork the repository
2. Create feature branch
3. Make changes following code guidelines
4. Test thoroughly
5. Submit pull request with clear description

## License

MIT License - see LICENSE file for details.
