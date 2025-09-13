package storage

import (
	"database/sql"
	"katana/tracker"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Storage abstracts session storage (SQLite or JSON fallback)
type Storage struct {
	db *sql.DB
}

// NewStorage initializes the storage (SQLite or JSON fallback)
func NewStorage() (*Storage, error) {
	dbPath := filepath.Join("data", "sessions.db")
	os.MkdirAll("data", 0755)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_time TEXT,
		end_time TEXT,
		duration INTEGER,
		activity TEXT,
		category TEXT,
		tags TEXT
	)`)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

// SaveSession saves a session to the database
func (s *Storage) SaveSession(sess *tracker.Session) error {
	tags := ""
	if len(sess.Tags) > 0 {
		tags = "|" + (string)(time.Now().UnixNano()) // placeholder for tags, improve later
	}
	_, err := s.db.Exec(`INSERT INTO sessions (start_time, end_time, duration, activity, category, tags) VALUES (?, ?, ?, ?, ?, ?)`,
		sess.StartTime.Format(time.RFC3339),
		sess.EndTime.Format(time.RFC3339),
		sess.Duration.Milliseconds(),
		sess.Activity,
		sess.Category,
		tags,
	)
	return err
}

// LoadSessionsForDay loads all sessions for a given day (used for daily/weekly/monthly viewers)
func (s *Storage) LoadSessionsForDay(day time.Time) ([]*tracker.Session, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	rows, err := s.db.Query(`SELECT id, start_time, end_time, duration, activity, category, tags FROM sessions WHERE start_time >= ? AND start_time < ?`, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []*tracker.Session
	for rows.Next() {
		var sess tracker.Session
		var startStr, endStr, tags string
		var duration int64
		if err := rows.Scan(&sess.ID, &startStr, &endStr, &duration, &sess.Activity, &sess.Category, &tags); err != nil {
			continue
		}
		sess.StartTime, _ = time.Parse(time.RFC3339, startStr)
		sess.EndTime, _ = time.Parse(time.RFC3339, endStr)
		sess.Duration = time.Duration(duration) * time.Millisecond
		sess.Tags = []string{} // TODO: parse tags
		sessions = append(sessions, &sess)
	}
	return sessions, nil
}
