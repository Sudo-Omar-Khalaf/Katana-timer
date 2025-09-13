package ui

import (
	"katana/storage"
	"katana/tracker"
	"katana/export"
	"log"
	"sync"
	"time"
	"fmt"
	"strings"
	"io"
	"os"

	"github.com/gen2brain/beeep"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

func init() {
	// Discard all log output by default for a silent terminal
	log.SetOutput(io.Discard)
}

type MainUI struct {
	Container      fyne.CanvasObject
	activityEntry  *widget.Entry
	startStopBtn   *widget.Button
	isTracking     bool
	currentSession *tracker.Session
	storage        *storage.Storage
	mu             sync.Mutex
	viewerTabs     *container.AppTabs
	timerLabel     *widget.Label
	activityList   *widget.List
	sessionsToday  []*tracker.Session
	updateActivityListPlaceholder func()
}

func NewMainUI() *MainUI {
	st, err := storage.NewStorage()
	if err != nil {
		// Temporarily enable error output for critical errors
		log.SetOutput(os.Stderr)
		log.Println("Failed to initialize storage: ", err)
		log.SetOutput(io.Discard)
	}
	sessionsToday, _ := st.LoadSessionsForDay(time.Now())
	ui := &MainUI{
		activityEntry: widget.NewEntry(),
		isTracking:    false,
		storage:       st,
		timerLabel:     widget.NewLabel("00:00:00"),
		sessionsToday:  sessionsToday,
	}
	ui.startStopBtn = widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func() {
		ui.toggleTracking()
	})

	// --- UI ENHANCEMENT: Modern pastel background, accent colors, classic layout ---
	appBg := canvas.NewRectangle(color.RGBA{R: 245, G: 248, B: 255, A: 255})

	title := canvas.NewText("Katana Time Tracker", color.RGBA{R: 60, G: 180, B: 220, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 24

	tagFilter := widget.NewEntry()
	tagFilter.SetPlaceHolder("Filter by tag (optional)")

	exportCSV := widget.NewButton("Export CSV", func() {
		sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
		export.ExportToCSV(sessions, "katana_export.csv")
	})
	exportJSON := widget.NewButton("Export JSON", func() {
		sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
		export.ExportToJSON(sessions, "katana_export.json")
	})
	exportPDF := widget.NewButton("Export PDF", func() {
		sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
		export.ExportToPDF(sessions, "katana_export.pdf")
	})

	// Accent color for buttons (handled by Importance, no need for variable)
	for _, btn := range []*widget.Button{ui.startStopBtn, exportCSV, exportJSON, exportPDF} {
		btn.Importance = widget.HighImportance
	}

	analyticsLabel := widget.NewLabel("")
	updateAnalytics := func() {
		today := time.Now()
		totalToday := 0.0
		totalWeek := 0.0
		totalMonth := 0.0
		sessionsToday, _ := ui.storage.LoadSessionsForDay(today)
		for _, s := range sessionsToday {
			totalToday += s.Duration.Hours()
		}
		for i := 0; i < 7; i++ {
			day := today.AddDate(0, 0, -i)
			s, _ := ui.storage.LoadSessionsForDay(day)
			for _, sess := range s {
				totalWeek += sess.Duration.Hours()
			}
		}
		for i := 0; i < 30; i++ {
			day := today.AddDate(0, 0, -i)
			s, _ := ui.storage.LoadSessionsForDay(day)
			for _, sess := range s {
				totalMonth += sess.Duration.Hours()
			}
		}
		analyticsLabel.SetText(fmt.Sprintf("Today: %.1fh | Week: %.1fh | Month: %.1fh", totalToday, totalWeek, totalMonth))
	}
	updateAnalytics()

	// --- Activity List ---
	ui.activityList = widget.NewList(
		func() int { return len(ui.sessionsToday) },
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(color.White)
			bg.CornerRadius = 8
			label := widget.NewLabel("")
			return container.NewStack(bg, label)
		},
		func(i int, o fyne.CanvasObject) {
			if i < len(ui.sessionsToday) {
				s := ui.sessionsToday[i]
				tags := ""
				if len(s.Tags) > 0 {
					tags = " #" + strings.Join(s.Tags, " #")
				}
				label := o.(*fyne.Container).Objects[1].(*widget.Label)
				label.TextStyle = fyne.TextStyle{Monospace: true, Bold: s.EndTime.IsZero()}
				label.SetText(fmt.Sprintf("%s - %s | %s [%s]%s", s.StartTime.Format("15:04"), s.EndTime.Format("15:04"), s.Activity, s.Category, tags))
				bg := o.(*fyne.Container).Objects[0].(*canvas.Rectangle)
				if s.EndTime.IsZero() {
					bg.FillColor = color.RGBA{R: 220, G: 245, B: 255, A: 255}
				} else {
					bg.FillColor = color.White
				}
				canvas.Refresh(bg)
			}
		},
	)

	// Activity list background and border
	activityListBorder := canvas.NewRectangle(color.RGBA{R: 220, G: 230, B: 245, A: 255})
	activityListBorder.StrokeColor = color.RGBA{R: 180, G: 200, B: 220, A: 255}
	activityListBorder.StrokeWidth = 1
	activityListScroll := container.NewVScroll(ui.activityList)
	// Make today's activity list bigger and allow user to resize dynamically
	activityListScroll.SetMinSize(fyne.NewSize(0, 180))

	placeholderLabel := widget.NewLabelWithStyle("No activities yet. Start tracking to see your sessions!", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	placeholderLabel.Hide()
	activityListStack := container.NewStack(activityListBorder, activityListScroll, placeholderLabel)

	ui.updateActivityListPlaceholder = func() {
		if len(ui.sessionsToday) == 0 {
			activityListScroll.Hide()
			placeholderLabel.Show()
		} else {
			activityListScroll.Show()
			placeholderLabel.Hide()
		}
	}
	ui.updateActivityListPlaceholder()

	// --- Viewers ---
	dailyGrid := makeHourGrid(sessionsToday)
	weeklyGrid := makeWeekGrid(st)
	monthlyGrid := makeMonthGrid(st)
	ui.viewerTabs = container.NewAppTabs(
		container.NewTabItem("Daily", dailyGrid),
		container.NewTabItem("Weekly", weeklyGrid),
		container.NewTabItem("Monthly", monthlyGrid),
	)
	ui.viewerTabs.SetTabLocation(container.TabLocationTop)

	// Use a VSplit to allow user to resize activity list and viewers dynamically
	centerSplit := container.NewVSplit(activityListStack, ui.viewerTabs)
	centerSplit.Offset = 0.4 // More space for activity list by default

	// --- Compose main content (classic layout, VSplit for vertical resizing) ---
	controls := container.NewVBox(
		container.NewCenter(title),
		widget.NewLabelWithStyle("Activity:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		ui.activityEntry,
		tagFilter,
		ui.startStopBtn,
		ui.timerLabel,
		analyticsLabel,
	)

	mainContent := container.NewBorder(
		controls, // top
		nil,      // bottom
		nil,      // left
		nil,      // right
		centerSplit, // center
	)

	ui.Container = container.NewStack(
		appBg,
		mainContent,
	)

	// --- Resource/logic enhancements for lightweight app ---
	// 1. Only update analytics and activity list when needed (not on every UI refresh)
	// 2. Reduce timer update frequency to 500ms for less CPU usage
	// 3. Use a single goroutine for all periodic background tasks
	// 4. Avoid unnecessary DB queries

	// --- Live timer and background update loop ---
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		lastDay := time.Now().Day()
		for {
			<-ticker.C
			ui.mu.Lock()
			if ui.isTracking && ui.currentSession != nil {
				dur := time.Since(ui.currentSession.StartTime)
				ui.timerLabel.SetText(dur.Truncate(time.Second).String())
				if dur.Hours() >= 2 {
					beeep.Notify("Katana Time Tracker", "Session running over 2 hours!", "")
				}
			} else {
				ui.timerLabel.SetText("00:00:00")
			}
			// Only refresh analytics and activity list if the day has changed
			if time.Now().Day() != lastDay {
				lastDay = time.Now().Day()
				ui.sessionsToday, _ = ui.storage.LoadSessionsForDay(time.Now())
				ui.activityList.Refresh()
				updateAnalytics()
				ui.updateActivityListPlaceholder()
			}
			ui.mu.Unlock()
		}
	}()

	return ui
}

// Helper: pastel color palette for comfort
func colorForCategory(cat string) color.Color {
	if cat == "" {
		return color.RGBA{R: 180, G: 220, B: 180, A: 255}
	}
	hash := 0
	for _, c := range cat {
		hash += int(c)
	}
	return color.RGBA{
		R: uint8(180 + (hash*37)%60),
		G: uint8(180 + (hash*53)%60),
		B: uint8(180 + (hash*97)%60),
		A: 255,
	}
}

// Daily 24h grid with hour labels, responsive
func makeHourGrid(sessions []*tracker.Session) fyne.CanvasObject {
	boxes := make([]fyne.CanvasObject, 24)
	labels := make([]fyne.CanvasObject, 24)
	for i := 0; i < 24; i++ {
		rect := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255})
		labels[i] = widget.NewLabelWithStyle(fmt.Sprintf("%02d:00", i), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		boxes[i] = container.NewStack(rect, labels[i])
	}
	for _, s := range sessions {
		catColor := colorForCategory(s.Category)
		startHour := s.StartTime.Hour()
		endHour := s.EndTime.Hour()
		for h := startHour; h <= endHour; h++ {
			if h >= 0 && h < 24 {
				rect := boxes[h].(*fyne.Container).Objects[0].(*canvas.Rectangle)
				rect.FillColor = catColor
				canvas.Refresh(rect)
			}
		}
	}
	// Use a scroll container to allow shrinking and remove VBox padding
	return container.NewVScroll(
		container.NewGridWrap(fyne.NewSize(48, 32), boxes...),
	)
}

// Weekly grid with day labels and total time tracked, responsive
func makeWeekGrid(storage *storage.Storage) fyne.CanvasObject {
	days := 7
	boxes := make([]fyne.CanvasObject, days)
	today := time.Now()
	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, -i)
		dayLabel := widget.NewLabelWithStyle(date.Format("Mon"), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		sessions, _ := storage.LoadSessionsForDay(date)
		total := 0.0
		cat := ""
		for _, s := range sessions {
			total += s.Duration.Hours()
			if cat == "" && s.Category != "" {
				cat = s.Category
			}
		}
		timeLabel := widget.NewLabelWithStyle(
			func() string {
				if total > 0 {
					return fmt.Sprintf("%.1fh", total)
				}
				return ""
			}(),
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		)
		boxContent := container.NewVBox(
			container.NewCenter(dayLabel),
			container.NewCenter(timeLabel),
		)
		rect := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255})
		if total > 0 {
			rect.FillColor = colorForCategory(cat)
		}
		boxes[days-1-i] = container.NewStack(rect, boxContent)
	}
	// Use a scroll container to allow shrinking and remove VBox padding
	return container.NewVScroll(
		container.NewGridWrap(fyne.NewSize(64, 40), boxes...),
	)
}

// Monthly grid with day-of-month labels, dynamic days, responsive
func makeMonthGrid(storage *storage.Storage) fyne.CanvasObject {
	today := time.Now()
	firstOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	nextMonth := firstOfMonth.AddDate(0, 1, 0)
	days := int(nextMonth.Sub(firstOfMonth).Hours() / 24)
	boxes := make([]fyne.CanvasObject, days)
	labels := make([]fyne.CanvasObject, days)
	for i := 0; i < days; i++ {
		date := firstOfMonth.AddDate(0, 0, i)
		labels[i] = widget.NewLabelWithStyle(fmt.Sprintf("%d", date.Day()), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		rect := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255})
		boxes[i] = container.NewStack(rect, labels[i])
	}
	for i := 0; i < days; i++ {
		date := firstOfMonth.AddDate(0, 0, i)
		sessions, _ := storage.LoadSessionsForDay(date)
		cat := ""
		for _, s := range sessions {
			if cat == "" && s.Category != "" {
				cat = s.Category
			}
		}
		if len(sessions) > 0 {
			rect := boxes[i].(*fyne.Container).Objects[0].(*canvas.Rectangle)
			rect.FillColor = colorForCategory(cat)
			canvas.Refresh(rect)
		}
	}
	// Use a scroll container to allow shrinking and remove VBox padding
	return container.NewVScroll(
		container.NewGridWrap(fyne.NewSize(36, 36), boxes...),
	)
}

func (ui *MainUI) toggleTracking() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if !ui.isTracking {
		activity := ui.activityEntry.Text
		if activity == "" {
			return
		}
		ui.currentSession = tracker.NewSession(activity)
		ui.isTracking = true
		ui.startStopBtn.SetText("Stop")
		ui.startStopBtn.SetIcon(theme.MediaStopIcon())
		ui.activityEntry.Disable()
	} else {
		ui.currentSession.Stop()
		err := ui.storage.SaveSession(ui.currentSession)
		if err != nil {
			log.SetOutput(os.Stderr)
			log.Println("Failed to save session:", err)
			log.SetOutput(io.Discard)
		}
		ui.sessionsToday, _ = ui.storage.LoadSessionsForDay(time.Now())
		ui.activityList.Refresh()
		sessionsToday := ui.sessionsToday
		ui.viewerTabs.Items[0].Content = makeHourGrid(sessionsToday)
		ui.viewerTabs.Items[1].Content = makeWeekGrid(ui.storage)
		ui.viewerTabs.Items[2].Content = makeMonthGrid(ui.storage)
		ui.viewerTabs.Refresh()
		ui.isTracking = false
		ui.startStopBtn.SetText("Start")
		ui.startStopBtn.SetIcon(theme.MediaPlayIcon())
		ui.activityEntry.Enable()
		ui.activityEntry.SetText("")
		ui.updateActivityListPlaceholder()
	}
}
