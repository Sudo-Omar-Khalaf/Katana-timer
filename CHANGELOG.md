# Katana Time Tracker - Change Log

## Recent Enhancements (v1.1.0)

### 🔧 Code Quality Improvements

1. **Deprecated Package Fixes**
   - Replaced deprecated `ioutil` package with modern `os` package functions
   - Updated `ioutil.ReadFile()` → `os.ReadFile()`
   - Updated `ioutil.WriteFile()` → `os.WriteFile()`

2. **Better Resource Management**
   - Added `Close()` method to Storage for proper database cleanup
   - Added `Cleanup()` method to MainUI for resource management
   - Implemented graceful shutdown handling in main.go
   - Added signal handling for SIGTERM and SIGINT

3. **Enhanced Error Handling**
   - Added comprehensive input validation with user-friendly error dialogs
   - Activity name validation (max 100 characters, non-empty)
   - Tag validation (max 20 characters each, max 5 tags total)
   - Session validation before starting and saving
   - Better error reporting for database operations

4. **New Storage Features**
   - Added `GetAllSessions()` method for complete data access
   - Improved error handling in JSON fallback mode
   - Better null handling in database queries

5. **Session Enhancements**
   - Added `Validate()` method to Session struct
   - Added `GetFormattedDuration()` method for human-readable durations
   - Better validation for session data integrity

6. **Configuration System**
   - Added `config.go` with configurable application settings
   - JSON-based configuration file (`data/config.json`)
   - Configurable notification thresholds, update intervals, and limits

### 🚀 Build & Development Improvements

1. **Enhanced .gitignore**
   - Comprehensive exclusion of build artifacts
   - Export files properly ignored
   - IDE files and OS-specific files excluded
   - Optional data directory exclusion

2. **Better Build Process**
   - No compilation warnings
   - Clean dependency management
   - Improved VS Code task configuration

### 🐛 Bug Fixes

1. **Memory Management**
   - Fixed potential memory leaks in background goroutines
   - Proper cleanup of database connections
   - Resource cleanup on application exit

2. **Data Integrity**
   - Added validation to prevent invalid session data
   - Better error handling for corrupt data files
   - Graceful fallback to defaults when configuration is invalid

3. **UI Stability**
   - Better error handling in UI operations
   - Validation before user actions
   - Improved feedback for user errors

### 📦 Dependencies

All dependencies are up-to-date and secure:
- Go 1.23.5+
- Fyne v2.6.3
- SQLite driver v1.14.32
- PDF generation library (jung-kurt/gofpdf)
- System notifications (gen2brain/beeep)

### 🔮 Future Enhancements

Planned improvements for future versions:
- Configuration UI panel
- Keyboard shortcuts
- Data import/export improvements
- Advanced analytics and reporting
- Plugin system for custom exporters
- Multiple workspace support

---

For installation and usage instructions, see the main [README.md](README.md).
