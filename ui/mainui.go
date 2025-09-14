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
	"fyne.io/fyne/v2/dialog"
	"image/color"
)

func init() {
	// Discard all log output by default for a silent terminal
	log.SetOutput(io.Discard)
}

type MainUI struct {
	Container      fyne.CanvasObject
	activityEntry  *widget.Entry
	tagEntry       *widget.Entry // New: for entering tags
	startStopBtn   *TerminalButton
	isTracking     bool
	currentSession *tracker.Session
	storage        *storage.Storage
	mu             sync.Mutex
	timerLabel     *widget.Label
	activityList   *widget.List
	sessionsToday  []*tracker.Session
	allSessionsToday []*tracker.Session
	updateActivityListPlaceholder func()
	originalTabLabels []string
	viewerContents []fyne.CanvasObject
	contentContainer *fyne.Container
	tabBar *TerminalTabBar
}

// --- TerminalButton: custom widget for terminal-style button with hover border and shadow/motion ---
type TerminalButton struct {
	widget.BaseWidget
	Label   string
	OnTap   func()
	isStop  bool // Track if button is in Stop state (only for timer button)
	isTimer bool // True if this is the Start/Stop timer button
	hovered bool
	active  bool
}

func NewTerminalButton(label string, onTap func()) *TerminalButton {
	btn := &TerminalButton{Label: label, OnTap: onTap, isTimer: false}
	btn.ExtendBaseWidget(btn)
	return btn
}

func NewTimerTerminalButton(label string, onTap func()) *TerminalButton {
	btn := &TerminalButton{Label: label, OnTap: onTap, isTimer: true}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *TerminalButton) SetLabel(label string) {
	b.Label = label
	b.Refresh()
}

func (b *TerminalButton) SetStopState(stop bool) {
	if b.isTimer {
		b.isStop = stop
		b.Refresh()
	}
}

func (b *TerminalButton) CreateRenderer() fyne.WidgetRenderer {
	terminalGreen := color.RGBA{0, 255, 0, 255}
	bg := canvas.NewRectangle(color.RGBA{0, 0, 0, 128}) // 50% transparent black
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 2
	border.StrokeColor = terminalGreen
	label := canvas.NewText(b.Label, terminalGreen)
	label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	label.Alignment = fyne.TextAlignCenter
	shadow := canvas.NewRectangle(color.RGBA{0, 255, 0, 40})
	shadow.Hide()
	objects := []fyne.CanvasObject{shadow, bg, label, border}
	return &terminalButtonRenderer{btn: b, bg: bg, border: border, label: label, shadow: shadow, objects: objects}
}

type terminalButtonRenderer struct {
	btn    *TerminalButton
	bg     *canvas.Rectangle
	border *canvas.Rectangle
	label  *canvas.Text
	shadow *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *terminalButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.border.Resize(size)
	r.label.Move(fyne.NewPos(0, (size.Height-r.label.MinSize().Height)/2))
	r.label.Resize(fyne.NewSize(size.Width, r.label.MinSize().Height))
	r.shadow.Resize(size)
}

func (r *terminalButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(120, 36)
}

func (r *terminalButtonRenderer) Refresh() {
	if r.btn.hovered || r.btn.active {
		r.border.StrokeColor = color.RGBA{0, 255, 0, 255} // bright green
		r.border.Show()
		r.shadow.Show()
	} else {
		r.border.StrokeColor = color.RGBA{0, 128, 0, 255} // dimmer green
		r.border.Show()
		r.shadow.Hide()
	}
	if r.btn.active {
		r.bg.FillColor = color.RGBA{0, 64, 0, 180}
	} else {
		r.bg.FillColor = color.RGBA{0, 0, 0, 128}
	}
	if r.btn.isTimer {
		if r.btn.isStop {
			r.label.Text = "Stop"
		} else {
			r.label.Text = "Start"
		}
	}
	canvas.Refresh(r.bg)
	canvas.Refresh(r.label)
	canvas.Refresh(r.border)
	canvas.Refresh(r.shadow)
}

func (r *terminalButtonRenderer) BackgroundColor() color.Color { return color.Transparent }
func (r *terminalButtonRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *terminalButtonRenderer) Destroy() {}

func (b *TerminalButton) MouseIn(*fyne.PointEvent) {
	b.hovered = true
	b.Refresh()
}
func (b *TerminalButton) MouseOut() {
	b.hovered = false
	b.active = false
	b.Refresh()
}
func (b *TerminalButton) MouseMoved(*fyne.PointEvent) {}
func (b *TerminalButton) Tapped(*fyne.PointEvent) {
	b.active = true
	b.Refresh()
	if b.OnTap != nil {
		b.OnTap()
	}
	b.active = false
	b.Refresh()
}

// --- TerminalTabButton: custom tab button for terminal-style tabs ---
type TerminalTabButton struct {
	widget.BaseWidget
	Label    string
	Selected bool
	OnTap    func()
	hovered  bool
}

func NewTerminalTabButton(label string, selected bool, onTap func()) *TerminalTabButton {
	btn := &TerminalTabButton{Label: label, Selected: selected, OnTap: onTap}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *TerminalTabButton) CreateRenderer() fyne.WidgetRenderer {
	terminalGreen := color.RGBA{0, 255, 0, 255}
	label := canvas.NewText(b.Label, terminalGreen)
	label.TextStyle = fyne.TextStyle{Monospace: true, Bold: b.Selected}
	label.Alignment = fyne.TextAlignCenter // Center the text horizontally
	underline := canvas.NewLine(terminalGreen)
	underline.StrokeWidth = 2
	if b.Selected {
		underline.Show()
	} else {
		underline.Hide()
	}
	if b.hovered {
		label.Color = color.RGBA{0, 255, 128, 255}
	}
	return &terminalTabButtonRenderer{btn: b, label: label, underline: underline, objects: []fyne.CanvasObject{label, underline}}
}

type terminalTabButtonRenderer struct {
	btn       *TerminalTabButton
	label     *canvas.Text
	underline *canvas.Line
	objects   []fyne.CanvasObject
}

func (r *terminalTabButtonRenderer) Layout(size fyne.Size) {
	// Center label horizontally and vertically
	labelSize := r.label.MinSize()
	labelX := (size.Width - labelSize.Width) / 2
	labelY := (size.Height - labelSize.Height) / 2
	r.label.Move(fyne.NewPos(labelX, labelY))
	r.label.Resize(labelSize)

	// Underline: match label width, center under text
	underlineY := labelY + labelSize.Height + 2
	underlineStart := labelX
	underlineEnd := labelX + labelSize.Width
	r.underline.Position1 = fyne.NewPos(underlineStart, underlineY)
	r.underline.Position2 = fyne.NewPos(underlineEnd, underlineY)
}

func (r *terminalTabButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(80, 28)
}

func (r *terminalTabButtonRenderer) Refresh() {
	terminalGreen := color.RGBA{0, 255, 0, 255}
	r.label.Text = r.btn.Label
	r.label.Color = terminalGreen
	r.label.TextStyle = fyne.TextStyle{Monospace: true, Bold: r.btn.Selected}
	if r.btn.Selected {
		r.underline.Show()
	} else {
		r.underline.Hide()
	}
	if r.btn.hovered {
		r.label.Color = color.RGBA{0, 255, 128, 255}
	}
	canvas.Refresh(r.label)
	canvas.Refresh(r.underline)
}

func (r *terminalTabButtonRenderer) BackgroundColor() color.Color { return color.Transparent }
func (r *terminalTabButtonRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *terminalTabButtonRenderer) Destroy() {}

func (b *TerminalTabButton) MouseIn(*fyne.PointEvent) {
	b.hovered = true
	b.Refresh()
}
func (b *TerminalTabButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}
func (b *TerminalTabButton) MouseMoved(*fyne.PointEvent) {}
func (b *TerminalTabButton) Tapped(*fyne.PointEvent) {
	if b.OnTap != nil {
		b.OnTap()
	}
}

// --- TerminalTabBar: horizontal row of TerminalTabButton ---
type TerminalTabBar struct {
	widget.BaseWidget
	Labels   []string
	Selected int
	OnSelect func(idx int)
	buttons  []*TerminalTabButton
}

func NewTerminalTabBar(labels []string, selected int, onSelect func(idx int)) *TerminalTabBar {
	bar := &TerminalTabBar{Labels: labels, Selected: selected, OnSelect: onSelect}
	bar.ExtendBaseWidget(bar)
	bar.buildButtons()
	return bar
}

func (b *TerminalTabBar) buildButtons() {
	b.buttons = make([]*TerminalTabButton, len(b.Labels))
	for i, label := range b.Labels {
		i := i
		b.buttons[i] = NewTerminalTabButton(label, i == b.Selected, func() {
			if b.OnSelect != nil {
				b.OnSelect(i)
			}
		})
	}
}

func (b *TerminalTabBar) CreateRenderer() fyne.WidgetRenderer {
	objs := make([]fyne.CanvasObject, len(b.buttons))
	for i, btn := range b.buttons {
		objs[i] = btn
	}
	return &terminalTabBarRenderer{bar: b, objects: objs}
}

type terminalTabBarRenderer struct {
	bar     *TerminalTabBar
	objects []fyne.CanvasObject
}

func (r *terminalTabBarRenderer) Layout(size fyne.Size) {
	btnW := size.Width / float32(len(r.bar.buttons))
	for i, btn := range r.bar.buttons {
		btn.Resize(fyne.NewSize(btnW, size.Height))
		btn.Move(fyne.NewPos(float32(i)*btnW, 0))
	}
}

func (r *terminalTabBarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(240, 32)
}

func (r *terminalTabBarRenderer) Refresh() {
	for _, btn := range r.bar.buttons {
		btn.Refresh()
	}
}
func (r *terminalTabBarRenderer) BackgroundColor() color.Color { return color.Transparent }
func (r *terminalTabBarRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *terminalTabBarRenderer) Destroy() {}

// --- TerminalEntry: custom entry widget for pure black background and green text/placeholder ---
type TerminalEntryWidget struct {
	widget.Entry
	placeholder string
	baseMinSize fyne.Size // cache default entry min size
	focused     bool      // track focus state
}

func NewTerminalEntry(placeholder string) *TerminalEntryWidget {
	te := &TerminalEntryWidget{placeholder: placeholder}
	te.ExtendBaseWidget(te)
	te.TextStyle = fyne.TextStyle{Monospace: true}
	te.Wrapping = fyne.TextTruncate
	te.SetPlaceHolder(placeholder)
	te.baseMinSize = widget.NewEntry().MinSize()
	return te
}

func (te *TerminalEntryWidget) FocusGained() {
	te.focused = true
	te.Refresh()
}
func (te *TerminalEntryWidget) FocusLost() {
	te.focused = false
	te.Refresh()
}

func (te *TerminalEntryWidget) CreateRenderer() fyne.WidgetRenderer {
	terminalGreen := color.RGBA{0, 255, 0, 255}
	bg := canvas.NewRectangle(color.Black)
	border := canvas.NewRectangle(color.Black)
	border.StrokeColor = terminalGreen
	border.StrokeWidth = 2
	border.FillColor = color.Black
	te.Entry.TextStyle = fyne.TextStyle{Monospace: true}
	te.Entry.Wrapping = fyne.TextTruncate
	te.Entry.PlaceHolder = te.placeholder
	// Remove placeholderText overlay, rely on Entry's built-in placeholder
	return widget.NewSimpleRenderer(container.NewMax(bg, border, &te.Entry))
}

// NewMainUI returns the main UI object and error for error handling
func NewMainUI() (*MainUI, error) {
	fyne.CurrentApp().Settings().SetTheme(&terminalTheme{})
	st, err := storage.NewStorage()
	if (err != nil) {
		return nil, err
	}
	sessionsToday, _ := st.LoadSessionsForDay(time.Now())
	ui := &MainUI{
		isTracking:    false,
		storage:       st,
		timerLabel:     widget.NewLabel("00:00:00"),
		sessionsToday:  sessionsToday,
		allSessionsToday: sessionsToday, // Store unfiltered sessions
		originalTabLabels: []string{"Daily", "Weekly", "Monthly"},
	}

	// --- UI ENHANCEMENT: Terminal-style solid black background, 50% transparent buttons, border on hover ---
	terminalGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	title := canvas.NewText("Katana Time Tracker", terminalGreen)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 24

	activityEntry := widget.NewEntry()
	activityEntry.SetPlaceHolder("Enter activity name")
	activityEntry.TextStyle = fyne.TextStyle{Monospace: true}
	tagEntry := widget.NewEntry()
	tagEntry.SetPlaceHolder("Tags (optional)")
	tagEntry.TextStyle = fyne.TextStyle{Monospace: true}
	tagFilterEntry := widget.NewEntry()
	tagFilterEntry.SetPlaceHolder("Filter by tag (optional)")
	tagFilterEntry.TextStyle = fyne.TextStyle{Monospace: true}
	tagFilterEntry.OnChanged = func(text string) {
		ui.FilterSessionsByTag(text)
	}
	ui.activityEntry = activityEntry
	ui.tagEntry = tagEntry

	exportCSV := NewTerminalButton("Export CSV", func() {
		dialog.NewFileSave(
			func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil { return }
				sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
				export.ExportToCSV(sessions, uc.URI().Path())
				uc.Close()
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		).Show()
	})
	exportJSON := NewTerminalButton("Export JSON", func() {
		dialog.NewFileSave(
			func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil { return }
				sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
				export.ExportToJSON(sessions, uc.URI().Path())
				uc.Close()
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		).Show()
	})
	exportPDF := NewTerminalButton("Export PDF", func() {
		dialog.NewFileSave(
			func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil { return }
				sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
				export.ExportToPDF(sessions, uc.URI().Path())
				uc.Close()
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		).Show()
	})

	startStopBtn := NewTimerTerminalButton("Start", func() {
		ui.toggleTracking()
	})
	ui.startStopBtn = startStopBtn
	ui.startStopBtn.SetStopState(ui.isTracking)

	analyticsText := canvas.NewText("", terminalGreen)
	analyticsText.TextStyle = fyne.TextStyle{Monospace: true}
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
		 analyticsText.Text = fmt.Sprintf("Today: %.1fh | Week: %.1fh | Month: %.1fh", totalToday, totalWeek, totalMonth)
		canvas.Refresh(analyticsText)
	}
	updateAnalytics()

	// --- Activity List ---
	activityListBorder := canvas.NewRectangle(terminalGreen)
	activityListBorder.StrokeColor = terminalGreen
	activityListBorder.StrokeWidth = 1
	ui.activityList = widget.NewList(
		func() int { return len(ui.sessionsToday) },
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(color.Black) // black background for activity row
			label := canvas.NewText("", color.RGBA{R: 180, G: 180, B: 180, A: 255}) // grey text for activity
			label.TextStyle = fyne.TextStyle{Monospace: true}
			return container.NewStack(bg, label)
		},
		func(i int, o fyne.CanvasObject) {
			if i < len(ui.sessionsToday) {
				s := ui.sessionsToday[i]
				tags := "[]"
				if len(s.Tags) > 0 {
					tags = "[" + strings.Join(s.Tags, ", ") + "]"
				}
				label := o.(*fyne.Container).Objects[1].(*canvas.Text)
				label.Text = fmt.Sprintf("%s - %s | %s %s", s.StartTime.Format("15:04"), s.EndTime.Format("15:04"), s.Activity, tags)
				canvas.Refresh(label)
			}
		},
	)

	activityListScroll := container.NewVScroll(ui.activityList)
	activityListScroll.SetMinSize(fyne.NewSize(0, 180))

	placeholderLabel := canvas.NewText("No activities yet. Start tracking to see your sessions!", terminalGreen)
	placeholderLabel.TextStyle = fyne.TextStyle{Italic: true, Monospace: true}
	placeholderLabel.Hide()
	activityListStack := container.NewStack(activityListBorder, container.NewMax(activityListScroll), placeholderLabel)

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
	dailyGrid := container.NewCenter(makeHourGrid(sessionsToday, terminalGreen))
	weeklyGrid := container.NewCenter(makeWeekGrid(st, terminalGreen))
	monthlyGrid := container.NewCenter(makeMonthGrid(st, terminalGreen))
	viewerContents := []fyne.CanvasObject{dailyGrid, weeklyGrid, monthlyGrid}
	ui.viewerContents = viewerContents
	selectedTab := 0
	ui.contentContainer = container.NewMax(ui.viewerContents[selectedTab])
	ui.tabBar = NewTerminalTabBar(ui.originalTabLabels, selectedTab, func(idx int) {
		ui.contentContainer.Objects = []fyne.CanvasObject{ui.viewerContents[idx]}
		ui.contentContainer.Refresh()
		for i, btn := range ui.tabBar.buttons {
			btn.Selected = (i == idx)
			btn.Refresh()
		}
	})
	centeredTabBar := container.NewCenter(ui.tabBar)
	// Use a VSplit to allow user to resize activity list and viewers dynamically
	centerSplit := container.NewVSplit(activityListStack, container.NewVBox(centeredTabBar, container.NewMax(ui.contentContainer)))
	centerSplit.Offset = 0.4 // More space for activity list by default

	// Timer label: green and monospace
	ui.timerLabel = widget.NewLabel("00:00:00")
	ui.timerLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	ui.timerLabel.Alignment = fyne.TextAlignCenter
	// We'll update its color using a canvas.Text overlay below

	// Overlay timerLabel with a green canvas.Text for color
	timerText := canvas.NewText("00:00:00", terminalGreen)
	timerText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	timerText.Alignment = fyne.TextAlignCenter
	ui.timerLabel.SetText("00:00:00")

	controls := container.NewVBox(
		container.NewCenter(title),
		canvas.NewText("Activity:", terminalGreen),
		activityEntry,
		canvas.NewText("Tags:", terminalGreen),
		tagEntry,
		tagFilterEntry,
		startStopBtn,
		exportCSV,
		exportJSON,
		exportPDF,
		container.NewCenter(timerText),
		analyticsText,
	)

	mainContent := container.NewVSplit(
		controls, // top: controls (activity entry, buttons, etc.)
		centerSplit, // bottom: activity list and viewers
	)
	mainContent.Offset = 0.22 // Adjust as needed for initial split
	ui.Container = container.NewMax(mainContent)

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
				h := int(dur.Hours())
				m := int(dur.Minutes()) % 60
				s := int(dur.Seconds()) % 60
				timerText.Text = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
				canvas.Refresh(timerText)
				if dur.Hours() >= 2 {
					beeep.Notify("Katana Time Tracker", "Session running over 2 hours!", "")
				}
			} else {
				timerText.Text = "00:00:00"
				canvas.Refresh(timerText)
			}
			// Only refresh analytics and activity list if the day has changed
			if time.Now().Day() != lastDay {
				lastDay = time.Now().Day()
				sessions, _ := ui.storage.LoadSessionsForDay(time.Now())
				ui.sessionsToday = sessions
				ui.allSessionsToday = sessions // Update unfiltered list
				ui.activityList.Refresh()
				updateAnalytics()
				ui.updateActivityListPlaceholder()
			}
			ui.mu.Unlock()
		}
	}()

	return ui, nil
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
func makeHourGrid(sessions []*tracker.Session, terminalGreen color.Color) fyne.CanvasObject {
	boxes := make([]fyne.CanvasObject, 24)
	// Build a map of hours with activity
	activeHours := make(map[int]bool)
	for _, s := range sessions {
		startHour := s.StartTime.Hour()
		endHour := s.EndTime.Hour()
		for h := startHour; h <= endHour; h++ {
			activeHours[h] = true
		}
	}
	for i := 0; i < 24; i++ {
		var rectColor, textColor color.Color
		if activeHours[i] {
			rectColor = terminalGreen
			textColor = color.Black
		} else {
			rectColor = color.Black
			textColor = terminalGreen
		}
		rect := canvas.NewRectangle(rectColor)
		rect.StrokeColor = terminalGreen
		rect.StrokeWidth = 1
		label := canvas.NewText(fmt.Sprintf("%02d:00", i), textColor)
		label.TextStyle = fyne.TextStyle{Monospace: true}
		label.Alignment = fyne.TextAlignCenter
		boxes[i] = container.NewMax(rect, container.NewCenter(label))
	}
	return container.NewGridWithColumns(8, boxes...)
}

// Weekly grid with day labels and total time tracked, responsive
func makeWeekGrid(storage *storage.Storage, terminalGreen color.Color) fyne.CanvasObject {
	days := 7
	boxes := make([]fyne.CanvasObject, days)
	today := time.Now()
	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, -i)
		sessions, _ := storage.LoadSessionsForDay(date)
		total := 0.0
		for _, s := range sessions {
			total += s.Duration.Hours()
		}
		var rectColor, textColor color.Color
		if total > 0 {
			rectColor = terminalGreen
			textColor = color.Black
		} else {
			rectColor = color.Black
			textColor = terminalGreen
		}
		rect := canvas.NewRectangle(rectColor)
		rect.StrokeColor = terminalGreen
		rect.StrokeWidth = 1
		dayLabel := canvas.NewText(date.Format("Mon"), textColor)
		dayLabel.TextStyle = fyne.TextStyle{Monospace: true}
		dayLabel.Alignment = fyne.TextAlignCenter
		timeLabel := canvas.NewText(func() string {
			if total > 0 {
				return fmt.Sprintf("%.1fh", total)
			}
			return ""
		}(), textColor)
		timeLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		timeLabel.Alignment = fyne.TextAlignCenter
		boxes[days-1-i] = container.NewMax(rect, container.NewVBox(
			container.NewCenter(dayLabel),
			container.NewCenter(timeLabel),
		))
	}
	return container.NewGridWithColumns(7, boxes...)
}

// Monthly grid with day-of-month labels, dynamic days, responsive
func makeMonthGrid(storage *storage.Storage, terminalGreen color.Color) fyne.CanvasObject {
	today := time.Now()
	firstOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	nextMonth := firstOfMonth.AddDate(0, 1, 0)
	days := int(nextMonth.Sub(firstOfMonth).Hours() / 24)
	boxes := make([]fyne.CanvasObject, days)
	for i := 0; i < days; i++ {
		date := firstOfMonth.AddDate(0, 0, i)
		sessions, _ := storage.LoadSessionsForDay(date)
		var rectColor, textColor color.Color
		if len(sessions) > 0 {
			rectColor = terminalGreen
			textColor = color.Black
		} else {
			rectColor = color.Black
			textColor = terminalGreen
		}
		rect := canvas.NewRectangle(rectColor)
		rect.StrokeColor = terminalGreen
		rect.StrokeWidth = 1
		label := canvas.NewText(fmt.Sprintf("%d", date.Day()), textColor)
		label.TextStyle = fyne.TextStyle{Monospace: true}
		label.Alignment = fyne.TextAlignCenter
		boxes[i] = container.NewMax(rect, container.NewCenter(label))
	}
	return container.NewGridWithColumns(8, boxes...)
}

func (ui *MainUI) toggleTracking() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if (!ui.isTracking) {
		activity := ui.activityEntry.Text
		tagsText := ui.tagEntry.Text
		if activity == "" {
			return
		}
		tags := []string{}
		for _, t := range strings.FieldsFunc(tagsText, func(r rune) bool { return r == ',' || r == ' ' }) {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
		sess := tracker.NewSession(activity)
		sess.Tags = tags
		ui.currentSession = sess
		ui.isTracking = true
		ui.startStopBtn.SetStopState(true)
		ui.activityEntry.Disable()
		ui.tagEntry.Disable()
	} else {
		ui.currentSession.Stop()
		err := ui.storage.SaveSession(ui.currentSession)
		if (err != nil) {
			log.SetOutput(os.Stderr)
			log.Println("Failed to save session:", err)
			log.SetOutput(io.Discard)
		}
		ui.sessionsToday, _ = ui.storage.LoadSessionsForDay(time.Now())
		ui.allSessionsToday = ui.sessionsToday // Update unfiltered list
		ui.activityList.Refresh()
		// --- Update tab content after session ends ---
		terminalGreen := color.RGBA{0, 255, 0, 255}
		ui.viewerContents[0] = container.NewCenter(makeHourGrid(ui.sessionsToday, terminalGreen))
		ui.viewerContents[1] = container.NewCenter(makeWeekGrid(ui.storage, terminalGreen))
		ui.viewerContents[2] = container.NewCenter(makeMonthGrid(ui.storage, terminalGreen))
		selectedTab := 0
		for i, btn := range ui.tabBar.buttons {
			if btn.Selected {
				selectedTab = i
				break
			}
		}
		ui.contentContainer.Objects = []fyne.CanvasObject{ui.viewerContents[selectedTab]}
		ui.contentContainer.Refresh()
		for i, btn := range ui.tabBar.buttons {
			btn.Selected = (i == selectedTab)
			btn.Refresh()
		}
		ui.isTracking = false
		ui.startStopBtn.SetStopState(false)
		ui.activityEntry.Enable()
		ui.tagEntry.Enable()
		ui.activityEntry.SetText("")
		ui.tagEntry.SetText("")
		ui.updateActivityListPlaceholder()
	}
}

// Add tag filtering to activity list
// Call this in tagFilterEntry.OnChanged
func (ui *MainUI) FilterSessionsByTag(tag string) {
	if tag == "" {
		ui.sessionsToday = ui.allSessionsToday
	} else {
		var filtered []*tracker.Session
		for _, s := range ui.allSessionsToday {
			for _, t := range s.Tags {
				if strings.Contains(strings.ToLower(t), strings.ToLower(tag)) {
					filtered = append(filtered, s)
					break
				}
			}
		}
		ui.sessionsToday = filtered
	}
	ui.activityList.Refresh()
	ui.updateActivityListPlaceholder()
}

// --- Custom theme for green tabs (always green, underline when selected) ---
type terminalTheme struct{}
func (t *terminalTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameSelection || n == theme.ColorNamePrimary || n == theme.ColorNameFocus || n == theme.ColorNameButton {
		return color.RGBA{0, 255, 0, 255}
	}
	if n == theme.ColorNameDisabled || n == theme.ColorNameDisabledButton {
		return color.RGBA{0, 128, 0, 128}
	}
	return theme.DefaultTheme().Color(n, v)
}
func (t *terminalTheme) Font(s fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(s) }
func (t *terminalTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (t *terminalTheme) Size(n fyne.ThemeSizeName) float32 { return theme.DefaultTheme().Size(n) }

// Helper to strip color/underline tags for tab text
func stripColorUnderline(s string) string {
	s = strings.ReplaceAll(s, "[color=#00ff00]", "")
	s = strings.ReplaceAll(s, "[/color]", "")
	s = strings.ReplaceAll(s, "[u]", "")
	s = strings.ReplaceAll(s, "[/u]", "")
	return s
}
