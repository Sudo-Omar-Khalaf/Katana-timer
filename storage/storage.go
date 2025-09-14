package storage

import (
	"database/sql"
	"katana/tracker"
	"os"
	"path/filepath"
	"time"
	"encoding/json"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
)

// Storage abstracts session storage (SQLite or JSON fallback)
type Storage struct {
	db *sql.DB
	jsonPath string
	useSQLite bool
}

// NewStorage initializes the storage (SQLite or JSON fallback)
func NewStorage() (*Storage, error) {
	dbPath := filepath.Join("data", "sessions.db")
	os.MkdirAll("data", 0755)
	db, err := sql.Open("sqlite3", dbPath)
	if err == nil {
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
		if err == nil {
			return &Storage{db: db, useSQLite: true}, nil
		}
	}
	// Fallback to JSON
	jsonPath := filepath.Join("data", "sessions.json")
	return &Storage{db: nil, jsonPath: jsonPath, useSQLite: false}, nil
}

// SaveSession saves a session to the database or JSON
func (s *Storage) SaveSession(sess *tracker.Session) error {
	if s.useSQLite {
		tagsJSON, _ := json.Marshal(sess.Tags)
		_, err := s.db.Exec(`INSERT INTO sessions (start_time, end_time, duration, activity, category, tags) VALUES (?, ?, ?, ?, ?, ?)`,
			sess.StartTime.Format(time.RFC3339),
			sess.EndTime.Format(time.RFC3339),
			sess.Duration.Milliseconds(),
			sess.Activity,
			sess.Category,
			string(tagsJSON),
		)
		return err
	}
	// Fallback: append to JSON file
	var sessions []*tracker.Session
	b, err := ioutil.ReadFile(s.jsonPath)
	if err == nil {
		json.Unmarshal(b, &sessions)
	}
	sessions = append(sessions, sess)
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.jsonPath, data, 0644)
}

// LoadSessionsForDay loads all sessions for a given day (used for daily/weekly/monthly viewers)
func (s *Storage) LoadSessionsForDay(day time.Time) ([]*tracker.Session, error) {
	var sessions []*tracker.Session
	if s.useSQLite {
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		end := start.Add(24 * time.Hour)
		rows, err := s.db.Query(`SELECT id, start_time, end_time, duration, activity, category, tags FROM sessions WHERE start_time >= ? AND start_time < ?`, start.Format(time.RFC3339), end.Format(time.RFC3339))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var sess tracker.Session
			var startStr, endStr, tagsStr string
			var duration int64
			if err := rows.Scan(&sess.ID, &startStr, &endStr, &duration, &sess.Activity, &sess.Category, &tagsStr); err != nil {
				continue
			}
			sess.StartTime, _ = time.Parse(time.RFC3339, startStr)
			sess.EndTime, _ = time.Parse(time.RFC3339, endStr)
			sess.Duration = time.Duration(duration) * time.Millisecond
			json.Unmarshal([]byte(tagsStr), &sess.Tags)
			sessions = append(sessions, &sess)
		}
		return sessions, nil
	}
	// Fallback: load from JSON file
	b, err := ioutil.ReadFile(s.jsonPath)
	if err != nil {
		return nil, nil
	}
	json.Unmarshal(b, &sessions)
	var filtered []*tracker.Session
	for _, sess := range sessions {
		if sameDay(sess.StartTime, day) {
			filtered = append(filtered, sess)
		}
	}
	return filtered, nil
}

func sameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
