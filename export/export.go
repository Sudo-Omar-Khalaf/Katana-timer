package export

import (
	"encoding/csv"
	"encoding/json"
	"katana/tracker"
	"os"
	"time"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// ExportToCSV exports sessions to a CSV file
func ExportToCSV(sessions []*tracker.Session, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"Start Time", "End Time", "Duration (min)", "Activity", "Category", "Tags"})
	for _, s := range sessions {
		tags := ""
		if len(s.Tags) > 0 {
			tags = jsonTags(s.Tags)
		}
		w.Write([]string{
			s.StartTime.Format(time.RFC3339),
			s.EndTime.Format(time.RFC3339),
			formatMinutes(s.Duration),
			s.Activity,
			s.Category,
			tags,
		})
	}
	return nil
}

// ExportToPDF exports sessions to a PDF file
func ExportToPDF(sessions []*tracker.Session, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Tracked Sessions")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 10)
	for _, s := range sessions {
		row := s.StartTime.Format("2006-01-02 15:04") + " - " + s.EndTime.Format("15:04") +
			" | " + s.Activity + " | " + s.Category + " | " + jsonTags(s.Tags) + " | " + formatMinutes(s.Duration) + " min"
		pdf.Cell(0, 8, row)
		pdf.Ln(8)
	}
	return pdf.OutputFileAndClose(filename)
}

// ExportMonthlyToCSV exports all sessions for current month grouped by day
func ExportMonthlyToCSV(storage interface{ LoadSessionsForMonth(int, time.Month) ([]*tracker.Session, error) }, filename string) error {
	now := time.Now()
	sessions, err := storage.LoadSessionsForMonth(now.Year(), now.Month())
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	// Write header
	w.Write([]string{"Date", "Start Time", "End Time", "Duration (min)", "Activity", "Category", "Tags", "Daily Total (min)"})

	// Group sessions by day
	dailySessions := make(map[string][]*tracker.Session)
	dailyTotals := make(map[string]float64)
	
	for _, s := range sessions {
		dayKey := s.StartTime.Format("2006-01-02")
		dailySessions[dayKey] = append(dailySessions[dayKey], s)
		dailyTotals[dayKey] += s.Duration.Minutes()
	}

	// Get all days in current month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	
	for d := firstDay; d.Day() <= lastDay.Day(); d = d.AddDate(0, 0, 1) {
		dayKey := d.Format("2006-01-02")
		daySessions := dailySessions[dayKey]
		
		if len(daySessions) == 0 {
			// No sessions for this day
			w.Write([]string{dayKey, "", "", "", "No activity", "", "", "0.0"})
		} else {
			// Write sessions for this day
			for i, s := range daySessions {
				tags := ""
				if len(s.Tags) > 0 {
					tags = jsonTags(s.Tags)
				}
				dailyTotal := ""
				if i == 0 { // Only show daily total on first row of the day
					dailyTotal = fmt.Sprintf("%.1f", dailyTotals[dayKey])
				}
				w.Write([]string{
					dayKey,
					s.StartTime.Format("15:04"),
					s.EndTime.Format("15:04"),
					formatMinutes(s.Duration),
					s.Activity,
					s.Category,
					tags,
					dailyTotal,
				})
			}
		}
	}
	return nil
}

// ExportMonthlyToPDF exports all sessions for current month grouped by day
func ExportMonthlyToPDF(storage interface{ LoadSessionsForMonth(int, time.Month) ([]*tracker.Session, error) }, filename string) error {
	now := time.Now()
	sessions, err := storage.LoadSessionsForMonth(now.Year(), now.Month())
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Time Tracking Report - %s %d", now.Month().String(), now.Year()))
	pdf.Ln(15)

	// Group sessions by day
	dailySessions := make(map[string][]*tracker.Session)
	dailyTotals := make(map[string]float64)
	
	for _, s := range sessions {
		dayKey := s.StartTime.Format("2006-01-02")
		dailySessions[dayKey] = append(dailySessions[dayKey], s)
		dailyTotals[dayKey] += s.Duration.Minutes()
	}

	// Get all days in current month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	
	monthlyTotal := 0.0
	
	for d := firstDay; d.Day() <= lastDay.Day(); d = d.AddDate(0, 0, 1) {
		dayKey := d.Format("2006-01-02")
		daySessions := dailySessions[dayKey]
		dayTotal := dailyTotals[dayKey]
		monthlyTotal += dayTotal
		
		// Day header
		pdf.SetFont("Arial", "B", 12)
		if dayTotal > 0 {
			pdf.Cell(0, 8, fmt.Sprintf("%s (%s) - %.1f minutes", dayKey, d.Format("Monday"), dayTotal))
		} else {
			pdf.Cell(0, 8, fmt.Sprintf("%s (%s) - No activity", dayKey, d.Format("Monday")))
		}
		pdf.Ln(10)
		
		if len(daySessions) > 0 {
			pdf.SetFont("Arial", "", 10)
			for _, s := range daySessions {
				sessionText := fmt.Sprintf("  %s - %s | %s | %s | %.1f min",
					s.StartTime.Format("15:04"),
					s.EndTime.Format("15:04"),
					s.Activity,
					strings.Join(s.Tags, ", "),
					s.Duration.Minutes())
				pdf.Cell(0, 6, sessionText)
				pdf.Ln(6)
			}
		}
		pdf.Ln(4)
	}
	
	// Monthly summary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("Monthly Total: %.1f minutes (%.1f hours)", monthlyTotal, monthlyTotal/60))
	
	return pdf.OutputFileAndClose(filename)
}

func jsonTags(tags []string) string {
	b, _ := json.Marshal(tags)
	return string(b)
}

func formatMinutes(d time.Duration) string {
	return fmt.Sprintf("%.1f", d.Minutes())
}
