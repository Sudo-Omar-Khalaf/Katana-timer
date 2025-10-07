package ui

import (
	"katana/storage"
	"katana/tracker"
	"katana/export"
	"katana/sound"
	"katana/power"
	"log"
	"sync"
	"time"
	"fmt"
	"strings"
	"io"
	"os"
	"strconv"

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
	// Time Tracker tab components
	activityEntry  *widget.Entry
	tagEntry       *widget.Entry // New: for entering tags
	startStopBtn   *TerminalButton
	isTracking     bool
	currentSession *tracker.Session
	storage        *storage.Storage
	soundPlayer    *sound.Player // Sound player for alarm sounds
	powerManager   *power.PowerManager // Power manager for sleep prevention
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
	notificationSent bool // Track if 2-hour notification has been sent
	
	// Main application tabs
	mainTabContainer *CustomMainTabContainer
	timeTrackerTab   *container.TabItem
	stopwatchTab     *container.TabItem
	countdownTab     *container.TabItem
	alarmTab         *container.TabItem
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
	
	// Update label text from button's Label field
	if r.btn.isTimer {
		if r.btn.isStop {
			r.label.Text = "Stop"
		} else {
			r.label.Text = "Start"
		}
	} else {
		// For regular buttons, use the Label field
		r.label.Text = r.btn.Label
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
	if err != nil {
		return nil, err
	}
	
	// Initialize sound player
	soundPlayer, err := sound.NewPlayer()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sound player: %w", err)
	}
	
	// Load alarm sounds from assets directory
	err = soundPlayer.LoadSoundsFromDirectory("assets/sounds")
	if err != nil {
		return nil, fmt.Errorf("failed to load alarm sounds: %w", err)
	}
	
	// Initialize power manager
	powerManager := power.NewPowerManager()
	
	sessionsToday, _ := st.LoadSessionsForDay(time.Now())
	ui := &MainUI{
		isTracking:       false,
		storage:          st,
		soundPlayer:      soundPlayer,
		powerManager:     powerManager,
		timerLabel:       widget.NewLabel("00:00:00"),
		sessionsToday:    sessionsToday,
		allSessionsToday: sessionsToday, // Store unfiltered sessions
		originalTabLabels: []string{"Daily", "Weekly", "Monthly"},
	}

	// Create the main application title
	terminalGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	appTitle := canvas.NewText("Katana Multi-Timer", terminalGreen)
	appTitle.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	appTitle.TextSize = 24

	// Create main tabs
	ui.timeTrackerTab = ui.createTimeTrackerTab()
	ui.stopwatchTab = ui.createStopwatchTab()
	ui.countdownTab = ui.createCountdownTab()
	ui.alarmTab = ui.createAlarmTab()

	// Create the main tab container with custom styling
	ui.mainTabContainer = NewCustomMainTabContainer(
		ui.timeTrackerTab,
		ui.stopwatchTab,
		ui.countdownTab,
		ui.alarmTab,
	)

	// Main application layout
	ui.Container = container.NewVBox(
		widget.NewSeparator(), // Additional top padding
		widget.NewSeparator(), // Additional top padding
		container.NewCenter(appTitle),
		ui.mainTabContainer,
	)

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
			// Show error dialog for empty activity
			dialog.NewError(
				fmt.Errorf("activity name cannot be empty"),
				fyne.CurrentApp().Driver().AllWindows()[0],
			).Show()
			return
		}
		// Validate activity length
		if len(strings.TrimSpace(activity)) > 100 {
			dialog.NewError(
				fmt.Errorf("activity name too long (max 100 characters)"),
				fyne.CurrentApp().Driver().AllWindows()[0],
			).Show()
			return
		}
		tags := []string{}
		for _, t := range strings.FieldsFunc(tagsText, func(r rune) bool { return r == ',' || r == ' ' }) {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				// Validate tag length
				if len(trimmed) > 20 {
					dialog.NewError(
						fmt.Errorf("tag '%s' too long (max 20 characters)", trimmed),
						fyne.CurrentApp().Driver().AllWindows()[0],
					).Show()
					return
				}
				tags = append(tags, trimmed)
			}
		}
		// Limit number of tags
		if len(tags) > 5 {
			dialog.NewError(
				fmt.Errorf("too many tags (max 5 allowed)"),
				fyne.CurrentApp().Driver().AllWindows()[0],
			).Show()
			return
		}
		sess := tracker.NewSession(activity)
		sess.Tags = tags
		// Validate session before starting
		if err := sess.Validate(); err != nil {
			dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
			return
		}
		ui.currentSession = sess
		ui.isTracking = true
		ui.notificationSent = false // Reset notification flag for new session
		ui.startStopBtn.SetStopState(true)
		ui.activityEntry.Disable()
		ui.tagEntry.Disable()
	} else {
		ui.currentSession.Stop()
		// Validate session before saving
		if err := ui.currentSession.Validate(); err != nil {
			log.SetOutput(os.Stderr)
			log.Printf("Session validation failed: %v", err)
			log.SetOutput(io.Discard)
			dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
			return
		}
		err := ui.storage.SaveSession(ui.currentSession)
		if err != nil {
			log.SetOutput(os.Stderr)
			log.Printf("Failed to save session: %v", err)
			log.SetOutput(io.Discard)
			dialog.NewError(
				fmt.Errorf("failed to save session: %v", err),
				fyne.CurrentApp().Driver().AllWindows()[0],
			).Show()
			return
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

// Cleanup properly shuts down the UI and releases resources
func (ui *MainUI) Cleanup() {
	if ui.storage != nil {
		ui.storage.Close()
	}
	if ui.soundPlayer != nil {
		ui.soundPlayer.Close()
	}
	if ui.powerManager != nil {
		ui.powerManager.Cleanup()
	}
}

// --- Custom theme for green tabs (always green, underline when selected) ---
type terminalTheme struct{}
func (t *terminalTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameSelection, theme.ColorNamePrimary, theme.ColorNameFocus:
		return color.RGBA{0, 255, 0, 255} // Bright green for selected/active elements
	case theme.ColorNameButton:
		return color.RGBA{0, 255, 0, 255} // Green for buttons
	case theme.ColorNameDisabled, theme.ColorNameDisabledButton:
		return color.RGBA{0, 128, 0, 128} // Dimmed green for disabled elements
	case theme.ColorNameForeground:
		return color.RGBA{0, 255, 0, 255} // Green text for foreground
	case theme.ColorNameBackground:
		return color.RGBA{0, 0, 0, 255} // Black background
	case theme.ColorNameHover:
		return color.RGBA{0, 255, 128, 255} // Lighter green for hover
	case theme.ColorNameHeaderBackground:
		return color.RGBA{0, 0, 0, 255} // Black background for tab headers
	case theme.ColorNameMenuBackground:
		return color.RGBA{0, 0, 0, 255} // Black background for menus
	case theme.ColorNameOverlayBackground:
		return color.RGBA{0, 0, 0, 200} // Semi-transparent black for overlays
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}
func (t *terminalTheme) Font(s fyne.TextStyle) fyne.Resource { 
	// Force monospace font for all text
	s.Monospace = true
	return theme.DefaultTheme().Font(s) 
}
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

// Helper function to format time for stopwatch display
func formatStopwatchTime(d time.Duration) string {
	totalMillis := d.Milliseconds()
	hours := totalMillis / (1000 * 60 * 60)
	minutes := (totalMillis % (1000 * 60 * 60)) / (1000 * 60)
	seconds := (totalMillis % (1000 * 60)) / 1000
	millis := totalMillis % 1000
	
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
}

// Helper function to safely parse integer from string
func parseIntSafe(s string) int {
	if s == "" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// createTimeTrackerTab creates the time tracking functionality tab
func (ui *MainUI) createTimeTrackerTab() *container.TabItem {
	terminalGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	title := canvas.NewText("Time Tracker", terminalGreen)
	title.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	title.TextSize = 20

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

	exportMonthlyCSV := NewTerminalButton("Export Month CSV", func() {
		dialog.NewFileSave(
			func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil { return }
				export.ExportMonthlyToCSV(ui.storage, uc.URI().Path())
				uc.Close()
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		).Show()
	})

	exportMonthlyPDF := NewTerminalButton("Export Month PDF", func() {
		dialog.NewFileSave(
			func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil { return }
				export.ExportMonthlyToPDF(ui.storage, uc.URI().Path())
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
	sessionsToday, _ := ui.storage.LoadSessionsForDay(time.Now())
	dailyGrid := container.NewCenter(makeHourGrid(sessionsToday, terminalGreen))
	weeklyGrid := container.NewCenter(makeWeekGrid(ui.storage, terminalGreen))
	monthlyGrid := container.NewCenter(makeMonthGrid(ui.storage, terminalGreen))
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

	// Overlay timerLabel with a green canvas.Text for color
	timerText := canvas.NewText("00:00:00", terminalGreen)
	timerText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	timerText.Alignment = fyne.TextAlignCenter
	ui.timerLabel.SetText("00:00:00")

	controls := container.NewVBox(
		widget.NewSeparator(), // Additional top padding
		widget.NewSeparator(), // Additional top padding
		widget.NewSeparator(), // Additional top padding
		widget.NewSeparator(), // Additional top padding
		container.NewCenter(title),
		canvas.NewText("Activity:", terminalGreen),
		activityEntry,
		canvas.NewText("Tags:", terminalGreen),
		tagEntry,
		tagFilterEntry,
		startStopBtn,
		container.NewGridWithColumns(2, exportCSV, exportPDF),
		container.NewGridWithColumns(2, exportMonthlyCSV, exportMonthlyPDF),
		container.NewCenter(timerText),
		analyticsText,
	)

	mainContent := container.NewVSplit(
		controls, // top: controls (activity entry, buttons, etc.)
		centerSplit, // bottom: activity list and viewers
	)
	mainContent.Offset = 0.22 // Adjust as needed for initial split

	// Start the background timer update goroutine for this tab
	ui.startTimeTrackerUpdates(timerText, updateAnalytics)

	return container.NewTabItem("Time Tracker", mainContent)
}

// CapturedTime represents a captured stopwatch time
type CapturedTime struct {
	Time       time.Duration
	Difference time.Duration
}

// createStopwatchTab creates the stopwatch functionality tab
func (ui *MainUI) createStopwatchTab() *container.TabItem {
	terminalGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Stopwatch state
	var stopwatchRunning bool
	var stopwatchStartTime time.Time
	var stopwatchElapsed time.Duration
	var capturedTimes []*CapturedTime

	// Stopwatch display - bigger size
	stopwatchDisplay := canvas.NewText("00:00:00.000", terminalGreen)
	stopwatchDisplay.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	stopwatchDisplay.TextSize = 48 // Bigger counter
	stopwatchDisplay.Alignment = fyne.TextAlignCenter

	// Captured times list
	captureList := widget.NewList(
		func() int { return len(capturedTimes) },
		func() fyne.CanvasObject {
			label := canvas.NewText("", color.RGBA{R: 180, G: 180, B: 180, A: 255})
			label.TextStyle = fyne.TextStyle{Monospace: true}
			label.Alignment = fyne.TextAlignLeading
			return label
		},
		func(i int, o fyne.CanvasObject) {
			if i < len(capturedTimes) {
				capture := capturedTimes[len(capturedTimes)-1-i] // Show most recent first
				label := o.(*canvas.Text)
				
				captureNum := len(capturedTimes) - i
				timeStr := formatStopwatchTime(capture.Time)
				
				if capture.Difference > 0 {
					diffStr := formatStopwatchTime(capture.Difference)
					label.Text = fmt.Sprintf("%d. %s (+%s)", captureNum, timeStr, diffStr)
				} else {
					label.Text = fmt.Sprintf("%d. %s", captureNum, timeStr)
				}
				canvas.Refresh(label)
			}
		},
	)

	captureScroll := container.NewVScroll(captureList)
	captureScroll.SetMinSize(fyne.NewSize(350, 200))

	// Captured times border
	captureBorder := canvas.NewRectangle(terminalGreen)
	captureBorder.StrokeColor = terminalGreen
	captureBorder.StrokeWidth = 1
	captureBorder.FillColor = color.Transparent

	captureContainer := container.NewStack(captureBorder, container.NewPadded(captureScroll))

	// Update function
	updateStopwatchDisplay := func() {
		var totalTime time.Duration
		if stopwatchRunning {
			totalTime = stopwatchElapsed + time.Since(stopwatchStartTime)
		} else {
			totalTime = stopwatchElapsed
		}
		stopwatchDisplay.Text = formatStopwatchTime(totalTime)
		canvas.Refresh(stopwatchDisplay)
	}

	// Start/Stop button
	startStopBtn := NewTerminalButton("Start", nil)
	
	startStopBtn.OnTap = func() {
		if !stopwatchRunning {
			// Start stopwatch
			stopwatchRunning = true
			stopwatchStartTime = time.Now()
			startStopBtn.SetLabel("Stop")
		} else {
			// Stop stopwatch
			stopwatchRunning = false
			stopwatchElapsed += time.Since(stopwatchStartTime)
			startStopBtn.SetLabel("Start")
		}
	}

	// Reset button
	resetBtn := NewTerminalButton("Reset", func() {
		stopwatchRunning = false
		stopwatchElapsed = 0
		capturedTimes = []*CapturedTime{}
		startStopBtn.SetLabel("Start")
		updateStopwatchDisplay()
		captureList.Refresh()
	})

	// Capture button
	captureBtn := NewTerminalButton("Capture", func() {
		var currentTime time.Duration
		if stopwatchRunning {
			currentTime = stopwatchElapsed + time.Since(stopwatchStartTime)
		} else {
			currentTime = stopwatchElapsed
		}
		
		if currentTime > 0 {
			var difference time.Duration
			if len(capturedTimes) > 0 {
				difference = currentTime - capturedTimes[len(capturedTimes)-1].Time
			}
			
			capture := &CapturedTime{
				Time:       currentTime,
				Difference: difference,
			}
			
			capturedTimes = append(capturedTimes, capture)
			captureList.Refresh()
		}
	})

	// Button container
	buttonContainer := container.NewGridWithColumns(3, startStopBtn, resetBtn, captureBtn)

	// Start background timer
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond) // High precision for smooth display
		defer ticker.Stop()
		for {
			<-ticker.C
			if stopwatchRunning {
				updateStopwatchDisplay()
			}
		}
	}()

	// Layout - clean and simple with better spacing and center alignment
	// Center the entire content vertically and horizontally
	centeredContent := container.NewCenter(
		container.NewVBox(
			widget.NewSeparator(),
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			container.NewCenter(stopwatchDisplay),
			widget.NewSeparator(),
			container.NewCenter(buttonContainer),
			widget.NewSeparator(),
			container.NewCenter(canvas.NewText("Captured Times", terminalGreen)),
			captureContainer,
		),
	)

	return container.NewTabItem("Stopwatch", centeredContent)
}

// createCountdownTab creates the countdown timer functionality tab
func (ui *MainUI) createCountdownTab() *container.TabItem {
	terminalGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Countdown state
	var countdownRunning bool
	var countdownEndTime time.Time
	var countdownDuration time.Duration

	// Countdown display
	countdownDisplay := canvas.NewText("00:00:00", terminalGreen)
	countdownDisplay.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	countdownDisplay.TextSize = 32
	countdownDisplay.Alignment = fyne.TextAlignCenter

	// Time input fields
	hoursEntry := widget.NewEntry()
	hoursEntry.SetText("0")
	hoursEntry.TextStyle = fyne.TextStyle{Monospace: true}
	hoursEntry.Validator = func(s string) error {
		if s == "" { return nil }
		if val, err := fmt.Sscanf(s, "%d", new(int)); err != nil || val != 1 {
			return fmt.Errorf("invalid number")
		}
		return nil
	}

	minutesEntry := widget.NewEntry()
	minutesEntry.SetText("5")
	minutesEntry.TextStyle = fyne.TextStyle{Monospace: true}
	minutesEntry.Validator = hoursEntry.Validator

	secondsEntry := widget.NewEntry()
	secondsEntry.SetText("0")
	secondsEntry.TextStyle = fyne.TextStyle{Monospace: true}
	secondsEntry.Validator = hoursEntry.Validator

	// Time input container
	timeInputContainer := container.NewGridWithColumns(5,
		hoursEntry,
		canvas.NewText(":", terminalGreen),
		minutesEntry,
		canvas.NewText(":", terminalGreen),
		secondsEntry,
	)

	// Update function
	updateCountdownDisplay := func() {
		if countdownRunning {
			remaining := time.Until(countdownEndTime)
			if remaining <= 0 {
				// Timer finished
				countdownRunning = false
				countdownDisplay.Text = "00:00:00"
				countdownDisplay.Color = color.RGBA{255, 0, 0, 255} // Red when finished
				
				// Show notification
				beeep.Notify("Katana Timer", "Countdown finished!", "")
				
				// Flash effect
				go func() {
					for i := 0; i < 6; i++ {
						if i%2 == 0 {
							countdownDisplay.Color = color.RGBA{255, 0, 0, 255} // Red
						} else {
							countdownDisplay.Color = terminalGreen // Green
						}
						canvas.Refresh(countdownDisplay)
						time.Sleep(500 * time.Millisecond)
					}
					countdownDisplay.Color = terminalGreen // Back to green
					canvas.Refresh(countdownDisplay)
				}()
			} else {
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60
				seconds := int(remaining.Seconds()) % 60
				countdownDisplay.Text = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
				
				// Change color based on remaining time
				if remaining <= 10*time.Second {
					countdownDisplay.Color = color.RGBA{255, 0, 0, 255} // Red for last 10 seconds
				} else if remaining <= 60*time.Second {
					countdownDisplay.Color = color.RGBA{255, 255, 0, 255} // Yellow for last minute
				} else {
					countdownDisplay.Color = terminalGreen // Green otherwise
				}
			}
		} else {
			countdownDisplay.Color = terminalGreen
		}
		canvas.Refresh(countdownDisplay)
	}

	// Start/Stop button
	startStopBtn := NewTerminalButton("Start", nil)
	
	startStopBtn.OnTap = func() {
		if !countdownRunning {
			// Parse time input
			hours := parseIntSafe(hoursEntry.Text)
			minutes := parseIntSafe(minutesEntry.Text)
			seconds := parseIntSafe(secondsEntry.Text)
			
			totalSeconds := hours*3600 + minutes*60 + seconds
			if totalSeconds <= 0 {
				// Show error for invalid time
				beeep.Alert("Invalid Time", "Please set a valid countdown time", "")
				return
			}
			
			// Start countdown
			countdownDuration = time.Duration(totalSeconds) * time.Second
			countdownEndTime = time.Now().Add(countdownDuration)
			countdownRunning = true
			startStopBtn.SetLabel("Stop")
			
			// Disable input fields
			hoursEntry.Disable()
			minutesEntry.Disable()
			secondsEntry.Disable()
		} else {
			// Stop countdown
			countdownRunning = false
			startStopBtn.SetLabel("Start")
			countdownDisplay.Color = terminalGreen
			
			// Re-enable input fields
			hoursEntry.Enable()
			minutesEntry.Enable()
			secondsEntry.Enable()
		}
	}

	// Reset button
	resetBtn := NewTerminalButton("Reset", func() {
		countdownRunning = false
		startStopBtn.SetLabel("Start")
		countdownDisplay.Text = "00:00:00"
		countdownDisplay.Color = terminalGreen
		
		// Re-enable input fields
		hoursEntry.Enable()
		minutesEntry.Enable()
		secondsEntry.Enable()
		
		updateCountdownDisplay()
	})

	// Button container
	buttonContainer := container.NewGridWithColumns(2, startStopBtn, resetBtn)

	// Progress bar
	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 1

	// Start background timer
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			<-ticker.C
			if countdownRunning {
				updateCountdownDisplay()
				
				// Update progress bar
				elapsed := time.Since(countdownEndTime.Add(-countdownDuration))
				progress := float64(elapsed) / float64(countdownDuration)
				if progress > 1 {
					progress = 1
				}
				progressBar.SetValue(progress)
			}
		}
	}()

	// Layout - centered content
	centeredContent := container.NewCenter(
		container.NewVBox(
			widget.NewSeparator(),
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			container.NewCenter(canvas.NewText("Set Time (HH:MM:SS):", terminalGreen)),
			container.NewCenter(timeInputContainer),
			widget.NewSeparator(),
			container.NewCenter(countdownDisplay),
			progressBar,
			container.NewCenter(buttonContainer),
		),
	)

	return container.NewTabItem("Countdown", centeredContent)
}

// Alarm represents a single alarm
type Alarm struct {
	ID          string
	Name        string
	Time        string // Format: "15:04"
	Enabled     bool
	Recurring   bool
	DaysOfWeek  []bool // [Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday]
	LastTriggered time.Time
	SoundName   string // Selected sound name
}

// createAlarmTab creates the alarm functionality tab
func (ui *MainUI) createAlarmTab() *container.TabItem {
	terminalGreen := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Alarm storage
	var alarms []*Alarm

	// Current time display
	timeDisplay := canvas.NewText("", terminalGreen)
	timeDisplay.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	timeDisplay.TextSize = 24
	timeDisplay.Alignment = fyne.TextAlignCenter

	// Update current time
	updateTimeDisplay := func() {
		now := time.Now()
		timeDisplay.Text = now.Format("15:04:05 Monday")
		canvas.Refresh(timeDisplay)
	}

	// Declare alarm list first
	var alarmList *widget.List
	
	// Alarm list implementation
	alarmList = widget.NewList(
		func() int { return len(alarms) },
		func() fyne.CanvasObject {
			// Enhanced alarm item with better visual separation
			bg := canvas.NewRectangle(color.RGBA{R: 20, G: 20, B: 20, A: 255})
			bg.StrokeColor = terminalGreen
			bg.StrokeWidth = 1
			
			nameLabel := canvas.NewText("", terminalGreen)
			nameLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
			nameLabel.TextSize = 14
			
			timeLabel := canvas.NewText("", color.RGBA{R: 200, G: 200, B: 200, A: 255})
			timeLabel.TextStyle = fyne.TextStyle{Monospace: true}
			timeLabel.TextSize = 12
			
			statusLabel := canvas.NewText("", terminalGreen)
			statusLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
			statusLabel.TextSize = 11
			
			// Smaller, better styled buttons
			toggleBtn := NewTerminalButton("Toggle", nil)
			deleteBtn := NewTerminalButton("Delete", nil)
			
			// Better layout with padding and spacing
			content := container.NewBorder(
				nil, nil, nil, nil,
				container.NewPadded(
					container.NewVBox(
						nameLabel,
						timeLabel,
						statusLabel,
						widget.NewSeparator(),
						container.NewGridWithColumns(2, toggleBtn, deleteBtn),
					),
				),
			)
			
			return container.NewStack(bg, content)
		},
		func(i int, o fyne.CanvasObject) {
			if i < len(alarms) {
				alarm := alarms[i]
				
				// Access the enhanced container structure
				stack := o.(*fyne.Container)
				borderContainer := stack.Objects[1].(*fyne.Container)
				paddedContainer := borderContainer.Objects[0].(*fyne.Container)
				vboxContainer := paddedContainer.Objects[0].(*fyne.Container)
				
				nameLabel := vboxContainer.Objects[0].(*canvas.Text)
				timeLabel := vboxContainer.Objects[1].(*canvas.Text)
				statusLabel := vboxContainer.Objects[2].(*canvas.Text)
				buttonContainer := vboxContainer.Objects[4].(*fyne.Container)
				toggleBtn := buttonContainer.Objects[0].(*TerminalButton)
				deleteBtn := buttonContainer.Objects[1].(*TerminalButton)
				
				// Enhanced alarm display with icons and better formatting
				nameLabel.Text = fmt.Sprintf("🔔 %s", alarm.Name)
				timeLabel.Text = fmt.Sprintf("⏰ %s", alarm.Time)
				if alarm.SoundName != "" {
					timeLabel.Text += fmt.Sprintf(" | 🔊 %s", alarm.SoundName)
				}
				
				if alarm.Enabled {
					statusLabel.Text = "Status: ACTIVE"
					statusLabel.Color = terminalGreen
					toggleBtn.SetLabel("Disable")
					// Highlight background for active alarms
					stack.Objects[0].(*canvas.Rectangle).FillColor = color.RGBA{R: 0, G: 40, B: 0, A: 255}
				} else {
					statusLabel.Text = "Status: DISABLED"
					statusLabel.Color = color.RGBA{R: 128, G: 128, B: 128, A: 255}
					toggleBtn.SetLabel("Enable")
					// Normal background for disabled alarms
					stack.Objects[0].(*canvas.Rectangle).FillColor = color.RGBA{R: 20, G: 20, B: 20, A: 255}
				}
				
				if alarm.Recurring {
					days := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
					recurringDays := []string{}
					for j, enabled := range alarm.DaysOfWeek {
						if enabled {
							recurringDays = append(recurringDays, days[j])
						}
					}
					if len(recurringDays) > 0 {
						timeLabel.Text += fmt.Sprintf(" (Recurring: %s)", strings.Join(recurringDays, ","))
					}
				}
				
				// Set button callbacks
				toggleBtn.OnTap = func() {
					alarm.Enabled = !alarm.Enabled
					
					// Manage system wake-up based on alarm state
					if alarm.Enabled {
						// Parse alarm time and schedule system wake-up
						if alarmTime, err := time.Parse("15:04", alarm.Time); err == nil {
							now := time.Now()
							alarmDateTime := time.Date(now.Year(), now.Month(), now.Day(), alarmTime.Hour(), alarmTime.Minute(), 0, 0, now.Location())
							
							// If alarm is for today but the time has passed, schedule for tomorrow
							if alarmDateTime.Before(now) {
								alarmDateTime = alarmDateTime.Add(24 * time.Hour)
							}
							
							// Schedule system wake-up for this alarm
							if err := ui.powerManager.ScheduleWakeup(alarm.ID, alarmDateTime); err != nil {
								log.Printf("Warning: Could not schedule system wake-up: %v", err)
							}
						}
					} else {
						// Cancel wake-up when alarm is disabled
						ui.powerManager.CancelWakeup(alarm.ID)
					}
					
					alarmList.Refresh()
				}
				
				deleteBtn.OnTap = func() {
					// Cancel wake-up when alarm is deleted
					if alarm.Enabled {
						ui.powerManager.CancelWakeup(alarm.ID)
					}
					
					// Remove alarm from slice
					for j, a := range alarms {
						if a.ID == alarm.ID {
							alarms = append(alarms[:j], alarms[j+1:]...)
							break
						}
					}
					alarmList.Refresh()
				}
				
				canvas.Refresh(nameLabel)
				canvas.Refresh(timeLabel)
				canvas.Refresh(statusLabel)
			}
		},
	)

	alarmScroll := container.NewVScroll(alarmList)
	alarmScroll.SetMinSize(fyne.NewSize(400, 200))

	// Alarm list border
	alarmBorder := canvas.NewRectangle(terminalGreen)
	alarmBorder.StrokeColor = terminalGreen
	alarmBorder.StrokeWidth = 1
	alarmBorder.FillColor = color.Transparent

	alarmContainer := container.NewStack(alarmBorder, container.NewMax(alarmScroll))

	// Add new alarm form - with proper container sizing
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name...")
	nameEntry.TextStyle = fyne.TextStyle{Monospace: true}
	
	// Create a properly sized name entry for grid layout
	nameEntry.Resize(fyne.NewSize(250, 36)) // Adjusted for grid layout
	
	// Create a simple container for the entry
	nameEntryContainer := container.NewBorder(nil, nil, nil, nil, nameEntry)

	hourEntry := widget.NewEntry()
	hourEntry.SetText("08")
	hourEntry.TextStyle = fyne.TextStyle{Monospace: true}

	minuteEntry := widget.NewEntry()
	minuteEntry.SetText("00")
	minuteEntry.TextStyle = fyne.TextStyle{Monospace: true}

	timeInputContainer := container.NewGridWithColumns(3,
		container.NewBorder(nil, nil, nil, nil, hourEntry),
		container.NewCenter(canvas.NewText(":", terminalGreen)),
		container.NewBorder(nil, nil, nil, nil, minuteEntry),
	)

	// Sound selection - with proper container and sizing
	availableSounds := ui.soundPlayer.GetAvailableSounds()
	
	// Create sound selection dropdown
	soundSelect := widget.NewSelect(availableSounds, nil)
	
	if len(availableSounds) > 0 {
		soundSelect.SetSelected(availableSounds[0]) // Default to first sound
	}

	// Test sound button
	testSoundBtn := NewTerminalButton("Play", func() {
		if soundSelect.Selected != "" {
			// Play for 3 seconds as a test
			ui.soundPlayer.PlaySound(soundSelect.Selected, 3*time.Second)
		}
	})
	
	// Create a properly sized sound container with spacing
	soundSelectContainer := container.NewBorder(nil, nil, nil, nil, soundSelect)
	soundContainer := container.NewGridWithColumns(2, 
		soundSelectContainer,
		testSoundBtn,
	)

	addAlarmBtn := NewTerminalButton("Add Alarm", func() {
		name := strings.TrimSpace(nameEntry.Text)
		if name == "" {
			name = "Alarm"
		}
		
		hour := parseIntSafe(hourEntry.Text)
		minute := parseIntSafe(minuteEntry.Text)
		
		if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
			beeep.Alert("Invalid Time", "Please enter valid time (00-23:00-59)", "")
			return
		}
		
		// Get selected sound
		selectedSound := soundSelect.Selected
		if selectedSound == "" && len(availableSounds) > 0 {
			selectedSound = availableSounds[0] // Default to first sound
		}
		
		// Create new alarm (simplified - no recurring)
		alarm := &Alarm{
			ID:        fmt.Sprintf("alarm_%d", time.Now().UnixNano()),
			Name:      name,
			Time:      fmt.Sprintf("%02d:%02d", hour, minute),
			Enabled:   true,
			Recurring: false,
			DaysOfWeek: make([]bool, 7),
			SoundName: selectedSound,
		}
		
		alarms = append(alarms, alarm)
		alarmList.Refresh()
		
		// Parse the alarm time and schedule system wake-up
		if alarmTime, err := time.Parse("15:04", alarm.Time); err == nil {
			now := time.Now()
			alarmDateTime := time.Date(now.Year(), now.Month(), now.Day(), alarmTime.Hour(), alarmTime.Minute(), 0, 0, now.Location())
			
			// If alarm is for today but the time has passed, schedule for tomorrow
			if alarmDateTime.Before(now) {
				alarmDateTime = alarmDateTime.Add(24 * time.Hour)
			}
			
			// Schedule system wake-up for this alarm
			if err := ui.powerManager.ScheduleWakeup(alarm.ID, alarmDateTime); err != nil {
				log.Printf("Warning: Could not schedule system wake-up: %v", err)
			}
		}
		
		// Clear form
		nameEntry.SetText("")
		hourEntry.SetText("08")
		minuteEntry.SetText("00")
		if len(availableSounds) > 0 {
			soundSelect.SetSelected(availableSounds[0])
		}
	})

	// Alarm checking goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		for {
			<-ticker.C
			updateTimeDisplay()
			
			now := time.Now()
			currentTime := now.Format("15:04")
			currentWeekday := int(now.Weekday())
			
			for _, alarm := range alarms {
				if !alarm.Enabled {
					continue
				}
				
				// Check if it's time for this alarm
				if alarm.Time == currentTime {
					shouldTrigger := false
					
					if alarm.Recurring {
						// Check if today is enabled for this recurring alarm
						if len(alarm.DaysOfWeek) > currentWeekday && alarm.DaysOfWeek[currentWeekday] {
							// Only trigger once per day
							today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
							if alarm.LastTriggered.Before(today) {
								shouldTrigger = true
								alarm.LastTriggered = now
							}
						}
					} else {
						// One-time alarm - trigger once then disable
						if alarm.LastTriggered.IsZero() || now.Sub(alarm.LastTriggered) > 24*time.Hour {
							shouldTrigger = true
							alarm.LastTriggered = now
							alarm.Enabled = false // Disable one-time alarms after triggering
						}
					}
					
					if shouldTrigger {
						// Trigger alarm notification
						beeep.Notify("Katana Alarm", fmt.Sprintf("Alarm: %s", alarm.Name), "")
						
						// Play alarm sound for 5 minutes
						if alarm.SoundName != "" {
							go func(soundName string) {
								ui.soundPlayer.PlaySound(soundName, 5*time.Minute)
							}(alarm.SoundName)
						}
						
						// Cancel wake-up after one-time alarm triggers (since it gets disabled)
						if !alarm.Recurring && alarm.Enabled == false {
							ui.powerManager.CancelWakeup(alarm.ID)
						}
						
						alarmList.Refresh()
					}
				}
			}
		}
	}()

	// Layout - clean and simple
	formContainer := container.NewVBox(
		container.NewGridWithColumns(2, 
			canvas.NewText("Alarm Name:", terminalGreen),
			nameEntryContainer,
		),
		widget.NewSeparator(),
		container.NewVBox(
			canvas.NewText("Time (HH:MM):", terminalGreen),
			timeInputContainer,
		),
		widget.NewSeparator(),
		container.NewVBox(
			canvas.NewText("Alarm Sound:", terminalGreen),
			soundContainer,
		),
		widget.NewSeparator(),
		addAlarmBtn,
	)

	// Stop alarm button for currently playing alarms
	stopAlarmBtn := NewTerminalButton("Stop Alarm", func() {
		ui.soundPlayer.StopSound()
	})

	// Better centered layout with proper spacing
	content := container.NewBorder(
		// Top: time display and stop button
		container.NewVBox(
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			widget.NewSeparator(), // Additional top padding
			container.NewCenter(timeDisplay),
			container.NewCenter(stopAlarmBtn),
			widget.NewSeparator(),
		),
		// Bottom: active alarms
		container.NewVBox(
			widget.NewSeparator(),
			container.NewCenter(canvas.NewText("Active Alarms:", terminalGreen)),
			alarmContainer,
		),
		nil, nil,
		// Center: alarm form
		container.NewCenter(formContainer),
	)

	return container.NewTabItem("Alarm", content)
}

// startTimeTrackerUpdates starts the background goroutine for time tracker updates
func (ui *MainUI) startTimeTrackerUpdates(timerText *canvas.Text, updateAnalytics func()) {
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
				// Send notification only once when crossing 2 hours
				if dur.Hours() >= 2 && !ui.notificationSent {
					beeep.Notify("Katana Time Tracker", "Session running over 2 hours!", "")
					// Play classic alarm sound for 2 seconds
					go func() {
						ui.soundPlayer.PlaySound("Classic Alarm 995", 2*time.Second)
					}()
					ui.notificationSent = true
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
}

// --- CustomMainTabContainer: uses TerminalTabButton styling for main app tabs ---
type CustomMainTabContainer struct {
	widget.BaseWidget
	tabs     []*container.TabItem
	selected int
	tabBar   *TerminalTabBar
	content  *fyne.Container
}

func NewCustomMainTabContainer(tabs ...*container.TabItem) *CustomMainTabContainer {
	c := &CustomMainTabContainer{
		tabs:     tabs,
		selected: 0,
	}
	c.ExtendBaseWidget(c)
	c.createContent()
	return c
}

func (c *CustomMainTabContainer) createContent() {
	// Create tab labels
	labels := make([]string, len(c.tabs))
	for i, tab := range c.tabs {
		labels[i] = tab.Text
	}
	
	// Create terminal tab bar
	c.tabBar = NewTerminalTabBar(labels, c.selected, func(idx int) {
		c.SelectTab(idx)
	})
	
	// Create content container
	if len(c.tabs) > 0 {
		c.content = container.NewMax(c.tabs[c.selected].Content)
	} else {
		c.content = container.NewMax()
	}
}

func (c *CustomMainTabContainer) SelectTab(index int) {
	if index >= 0 && index < len(c.tabs) {
		c.selected = index
		c.content.Objects = []fyne.CanvasObject{c.tabs[index].Content}
		c.content.Refresh()
		
		// Update tab bar selection
		for i, btn := range c.tabBar.buttons {
			btn.Selected = (i == index)
			btn.Refresh()
		}
	}
}

func (c *CustomMainTabContainer) CreateRenderer() fyne.WidgetRenderer {
	container := container.NewVBox(
		c.tabBar,
		c.content,
	)
	return widget.NewSimpleRenderer(container)
}
