# Contributing to Katana Multi-Timer

Thank you for your interest in contributing to Katana! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/Sudo-Omar-Khalaf/katana.git`
3. **Create** a feature branch: `git checkout -b feature-amazing-feature`
4. **Make** your changes
5. **Test** thoroughly
6. **Submit** a pull request

## 🛠️ Development Setup

### Prerequisites
- Go 1.19 or higher
- Git
- Linux development environment
- Audio system (ALSA/PulseAudio)

### Local Development
```bash
# Clone the repository
git clone https://github.com/Sudo-Omar-Khalaf/katana.git
cd katana

# Install dependencies
go mod download

# Run the application
go run main.go

# Build for testing
go build -o katana-dev
./katana-dev
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📋 Development Guidelines

### Code Style
- Follow standard Go formatting: `gofmt`
- Use `go vet` to check for common issues
- Keep functions small and focused
- Add comments for public functions and complex logic
- Use meaningful variable and function names

### Commit Messages
Use conventional commit format:
```
type(scope): description

feat(alarm): add recurring alarm functionality
fix(ui): resolve timer display issue
docs(readme): update installation instructions
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Pull Request Guidelines
1. **One feature per PR**: Keep changes focused and small
2. **Update documentation**: Include relevant documentation updates
3. **Add tests**: Ensure new features are tested
4. **Test on multiple distros**: Verify compatibility across Linux distributions
5. **Update changelog**: Add entry to CHANGELOG.md if applicable

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Operating System**: Distribution and version
2. **Go Version**: Output of `go version`
3. **Steps to Reproduce**: Detailed steps to recreate the issue
4. **Expected Behavior**: What you expected to happen
5. **Actual Behavior**: What actually happened
6. **Logs**: Any relevant error messages or logs

### Bug Report Template
```markdown
**Environment:**
- OS: Ubuntu 22.04
- Go Version: 1.19.5
- Katana Version: v1.0.0

**Steps to Reproduce:**
1. Open Katana
2. Set an alarm for 5 minutes
3. Close the application

**Expected:** Alarm should trigger after 5 minutes
**Actual:** No alarm sound or notification

**Additional Context:**
Audio works fine with other applications.
```

## ✨ Feature Requests

We welcome feature requests! Please:

1. **Check existing issues** to avoid duplicates
2. **Describe the problem** you're trying to solve
3. **Propose a solution** if you have ideas
4. **Consider the scope** - keep features focused and useful for most users

### Feature Request Template
```markdown
**Problem:** As a user, I want to... because...

**Proposed Solution:** Add a feature that...

**Alternative Solutions:** Other ways to solve this could be...

**Additional Context:** This would be useful for...
```

## 🏗️ Project Structure

```
katana/
├── main.go              # Application entry point
├── ui/                  # User interface components
│   └── mainui.go       # Main UI implementation
├── power/              # Power management (wake-up functionality)
│   └── power.go        # Cross-platform wake scheduling
├── sound/              # Audio playback
│   └── player.go       # Sound player implementation
├── storage/            # Data persistence
│   └── storage.go      # Database operations
├── tracker/            # Time tracking logic
│   └── session.go      # Session management
├── export/             # Data export functionality
│   └── export.go       # CSV/PDF export
└── assets/             # Static assets
    └── sounds/         # Alarm sound files
```

## 🎯 Areas for Contribution

### High Priority
- [ ] Unit tests for all modules
- [ ] Integration tests
- [ ] Performance optimizations
- [ ] Memory usage improvements
- [ ] Better error handling

### Features
- [ ] Custom alarm sounds upload
- [ ] Themes and customization
- [ ] Keyboard shortcuts
- [ ] Command-line interface
- [ ] Multiple timezone support
- [ ] Notification customization

### Platform Support
- [ ] Wayland compatibility testing
- [ ] ARM64 optimization
- [ ] Flatpak packaging
- [ ] AppImage creation
- [ ] Snap package

### Documentation
- [ ] API documentation
- [ ] User manual
- [ ] Video tutorials
- [ ] Troubleshooting guide
- [ ] Platform-specific installation guides

## 🧪 Testing

### Manual Testing Checklist
- [ ] Alarm functionality (sound, notifications)
- [ ] Stopwatch accuracy and lap timing
- [ ] Countdown timer with alerts
- [ ] Time tracker session management
- [ ] Data export (CSV/PDF)
- [ ] System wake-up functionality
- [ ] UI responsiveness across different screen sizes
- [ ] Audio compatibility (ALSA/PulseAudio)

### Automated Testing
- Write unit tests for new functions
- Test error handling paths
- Verify data integrity in storage operations
- Test cross-platform compatibility

## 📦 Release Process

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Update CHANGELOG.md with new features and fixes
3. **Testing**: Comprehensive testing on multiple distributions
4. **Documentation**: Update README and documentation
5. **Release**: Create GitHub release with binaries

## 🤝 Community

- **Discussions**: Use GitHub Discussions for questions and ideas
- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Code Review**: All contributions go through code review
- **Be Respectful**: Follow our code of conduct

## 📞 Getting Help

- **Documentation**: Check the README and wiki first
- **GitHub Discussions**: Ask questions in the community
- **Issues**: Create an issue for bugs or specific problems

## 🙏 Recognition

Contributors will be:
- Listed in the README contributors section
- Mentioned in release notes for significant contributions
- Given credit in commit messages and pull requests

Thank you for helping make Katana better for the entire Linux community! 🎉
