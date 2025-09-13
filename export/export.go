package export

import (
	"encoding/csv"
	"encoding/json"
	"katana/tracker"
	"os"
	"time"
	"fmt"

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

// ExportToJSON exports sessions to a JSON file
func ExportToJSON(sessions []*tracker.Session, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(sessions)
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

func jsonTags(tags []string) string {
	b, _ := json.Marshal(tags)
	return string(b)
}

func formatMinutes(d time.Duration) string {
	return fmt.Sprintf("%.1f", d.Minutes())
}
