# Katana Time Tracker

Katana is a lightweight, visually appealing, and flexible time tracking desktop application built in Go using the Fyne GUI toolkit. It is designed for minimal resource usage, a classic but modern UI, and robust time tracking features for daily, weekly, and monthly productivity.

## Features
- Start/stop time tracking with activity input
- Tag and category support for activities
- Daily, weekly, and monthly viewers with analytics
- Export tracked data to CSV, JSON, or PDF
- Tag/category filtering and analytics
- Desktop notifications for long sessions
- Responsive, resizable, and modern UI with pastel color palette
- Minimal terminal/log output for a silent experience

## Installation

### Prerequisites
- Go 1.18 or newer ([Download Go](https://golang.org/dl/))
- A working C compiler (required for Fyne and SQLite)
- Linux, Windows, or macOS

### 1. Clone the Repository
```
git clone https://github.com/yourusername/katana.git
cd katana
```

### 2. Install Dependencies
All dependencies are managed via Go modules. To install them:
```
go mod tidy
```

### 3. Build and Run
To run the app directly:
```
go run main.go
```

Or to build a standalone binary:
```
go build -o katana
./katana
```

## Usage
- Enter an activity name and (optionally) tags/categories, then click **Start** to begin tracking.
- Click **Stop** to end the session. Your activity will appear in the list.
- Use the tabs to view your tracked time by day, week, or month.
- Use the export buttons to save your data as CSV, JSON, or PDF.
- Filter activities by tag using the filter box.
- Analytics are shown for today, this week, and this month.
- The app will notify you if a session runs over 2 hours.

## Data Storage
- By default, Katana stores data in a local SQLite database (`data/sessions.db`).
- If SQLite is unavailable, it will fallback to a JSON file (`katana_export.json`).

## Troubleshooting
- If you see errors about missing C libraries, ensure you have a C compiler installed (e.g., `build-essential` on Ubuntu/Debian).
- If the UI does not appear, make sure your system supports OpenGL and you have the latest graphics drivers.
- All non-critical log output is suppressed by default. Only critical errors will appear in the terminal.

## Customization
- The UI uses a modern pastel color palette and classic layout. You can adjust colors in `ui/mainui.go`.
- To change export formats or analytics, see the `export/` and `ui/` packages.

## License
MIT License. See [LICENSE](LICENSE) for details.

---

Enjoy tracking your time with Katana! For issues or feature requests, please open an issue on GitHub.
