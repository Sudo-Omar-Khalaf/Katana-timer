package tracker

import (
	"fmt"
	"strings"
	"time"
)

// Session represents a single time tracking session
type Session struct {
	ID        int64
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Activity  string
	Category  string
	Tags      []string
}

// NewSession creates a new tracking session
func NewSession(activity string) *Session {
	// Parse activity string for category and tags (e.g., "study:math #important")
	category := ""
	tags := []string{}

	// Extract category (e.g., "study:math" -> category="study")
	if idx := strings.Index(activity, ":"); idx != -1 {
		category = activity[:idx]
		activity = activity[idx+1:]
	}

	// Extract tags (words starting with #)
	words := strings.Fields(activity)
	cleanActivity := []string{}

	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tags = append(tags, word[1:])
		} else {
			cleanActivity = append(cleanActivity, word)
		}
	}

	return &Session{
		StartTime: time.Now(),
		Activity:  strings.Join(cleanActivity, " "),
		Category:  category,
		Tags:      tags,
	}
}

// Stop ends the current tracking session
func (s *Session) Stop() {
	s.EndTime = time.Now()
	s.Duration = s.EndTime.Sub(s.StartTime)
}

// Validate checks if the session data is valid
func (s *Session) Validate() error {
	if s.Activity == "" {
		return fmt.Errorf("activity cannot be empty")
	}
	if s.StartTime.IsZero() {
		return fmt.Errorf("start time cannot be zero")
	}
	if !s.EndTime.IsZero() && s.EndTime.Before(s.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}
	return nil
}

// GetFormattedDuration returns a human-readable duration string
func (s *Session) GetFormattedDuration() string {
	if s.Duration == 0 {
		return "0m"
	}
	
	hours := int(s.Duration.Hours())
	minutes := int(s.Duration.Minutes()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
