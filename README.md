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

## Installation (Linux)

### Prerequisites
- **Go 1.18 or newer** ([Download Go](https://golang.org/dl/))
- **C compiler and build tools** (required for Fyne and SQLite)
- **pkg-config** and **libgl1-mesa-dev** (for Fyne/OpenGL)
- **git**

Install all dependencies on Ubuntu/Debian:
```sh
sudo apt update
sudo apt install -y build-essential pkg-config libgl1-mesa-dev git golang
```

If you use another Linux distribution, install the equivalent packages for your system.

### 1. Clone the Repository
```sh
git clone https://github.com/Sudo-Omar-Khalaf/Katana-timer.git
cd Katana-timer
```

### 2. Install Go Dependencies
All dependencies are managed via Go modules:
```sh
go mod tidy
```

### 3. Build and Run
To run the app directly:
```sh
go run main.go
```

Or to build a standalone binary:
```sh
go build -o katana
./katana
```

## Quick Install (One Command)

After cloning the repo, you can install all dependencies and launch Katana with:

```sh
while read dep; do sudo apt install -y "$dep" || pip install "$dep" || go install "$dep"; done < requirements.txt && go mod tidy && go run main.go
```

This command will:
- Install each dependency from `requirements.txt` using apt, pip, or go install as appropriate
- Run `go mod tidy` to fetch Go modules
- Launch the Katana app

---

## Usage
- Enter an activity name and (optionally) tags/categories, then click **Start** to begin tracking.
- Click **Stop** to end the session. Your activity will appear in the list.
- Use the tabs to view your tracked time by day, week, or month.
- Use the export buttons to save your data as CSV, JSON, or PDF.
- Filter activities by tag using the filter box.
- Analytics are shown for today, this week, and this month.
- The app will notify you if a session runs over 2 hours.

## Quick Links & Dependency Installation

**[View Katana on GitHub](https://github.com/Sudo-Omar-Khalaf/Katana-timer)**

If you want to install all requirements from a `requirements.txt` file (for Python or extra tools), run:

```sh
while read dep; do sudo apt install -y "$dep" || pip install "$dep" || go install "$dep"; done < requirements.txt
```

This command will attempt to install each dependency listed in `requirements.txt` using `apt`, `pip`, or `go install` as appropriate. Make sure to review the file and adjust as needed for your environment.

---

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
