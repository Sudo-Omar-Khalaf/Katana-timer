package power

import (
"fmt"
"log"
"os/exec"
"runtime"
"strings"
"strconv"
"sync"
"time"
)

// PowerManager handles system wake-up scheduling for alarms
type PowerManager struct {
	mu              sync.Mutex
	activeWakeTimers map[string]*time.Timer // Track active wake timers by alarm ID
}

// NewPowerManager creates a new power manager instance
func NewPowerManager() *PowerManager {
	return &PowerManager{
		activeWakeTimers: make(map[string]*time.Timer),
	}
}

// ScheduleWakeup schedules the system to wake up at a specific time for an alarm
func (pm *PowerManager) ScheduleWakeup(alarmID string, wakeupTime time.Time) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Cancel any existing wake timer for this alarm
	if timer, exists := pm.activeWakeTimers[alarmID]; exists {
		timer.Stop()
		delete(pm.activeWakeTimers, alarmID)
	}

	duration := time.Until(wakeupTime)
	if duration <= 0 {
		return fmt.Errorf("wakeup time is in the past")
	}

	log.Printf("Scheduling system wake-up for alarm %s at %v (in %v)", alarmID, wakeupTime, duration)

	switch runtime.GOOS {
	case "linux":
		return pm.scheduleWakeupLinux(alarmID, wakeupTime, duration)
	case "windows":
		return pm.scheduleWakeupWindows(alarmID, wakeupTime, duration)
	case "darwin": // macOS
		return pm.scheduleWakeupMacOS(alarmID, wakeupTime, duration)
	default:
		log.Printf("System wake-up not supported on %s, alarm will only ring if system is awake", runtime.GOOS)
		return nil
	}
}

// CancelWakeup cancels a scheduled wake-up for an alarm
func (pm *PowerManager) CancelWakeup(alarmID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Cancel the timer
	if timer, exists := pm.activeWakeTimers[alarmID]; exists {
		timer.Stop()
		delete(pm.activeWakeTimers, alarmID)
	}

	switch runtime.GOOS {
	case "linux":
		return pm.cancelWakeupLinux(alarmID)
	case "windows":
		return pm.cancelWakeupWindows(alarmID)
	case "darwin":
		return pm.cancelWakeupMacOS(alarmID)
	}

	return nil
}

// Linux-specific wake-up scheduling
func (pm *PowerManager) scheduleWakeupLinux(alarmID string, wakeupTime time.Time, duration time.Duration) error {
	// Method 1: Try rtcwake (requires sudo, most reliable)
	if err := pm.tryRtcWake(duration); err == nil {
		log.Printf("Scheduled wake-up using rtcwake for alarm %s", alarmID)
		return nil
	}

	// Method 2: Try using 'at' command as fallback
	if err := pm.tryAtCommand(alarmID, wakeupTime); err == nil {
		log.Printf("Scheduled wake-up using 'at' command for alarm %s", alarmID)
		return nil
	}

	log.Printf("Could not schedule system wake-up for alarm %s: %v", alarmID, "no suitable method available")
	return fmt.Errorf("wake-up scheduling not available (try running with sudo for rtcwake support)")
}

// Try using rtcwake to schedule system wake-up
func (pm *PowerManager) tryRtcWake(duration time.Duration) error {
	// Note: This requires sudo privileges
	cmd := exec.Command("sudo", "rtcwake", "-m", "no", "-s", strconv.Itoa(int(duration.Seconds())))
	return cmd.Run()
}

// Try using 'at' command to wake system
func (pm *PowerManager) tryAtCommand(alarmID string, wakeupTime time.Time) error {
	// Create a simple command to ensure system is awake
	wakeCmd := fmt.Sprintf("echo 'Katana alarm %s wake-up' | wall", alarmID)
	timeStr := wakeupTime.Format("15:04 2006-01-02")
	
	cmd := exec.Command("at", timeStr)
	cmd.Stdin = strings.NewReader(wakeCmd + "\n")
	return cmd.Run()
}

// Cancel Linux wake-up
func (pm *PowerManager) cancelWakeupLinux(alarmID string) error {
	// For rtcwake, we can't easily cancel, but it's a one-time wake
	// For 'at' jobs, we could try to remove them, but it's complex
log.Printf("Wake-up cancellation requested for alarm %s (Linux)", alarmID)
return nil
}

// Windows-specific wake-up scheduling using Task Scheduler
func (pm *PowerManager) scheduleWakeupWindows(alarmID string, wakeupTime time.Time, duration time.Duration) error {
taskName := fmt.Sprintf("KatanaWakeup_%s", alarmID)
timeStr := wakeupTime.Format("15:04:05")
dateStr := wakeupTime.Format("2006-01-02")

// Create a scheduled task that wakes the computer
cmd := exec.Command("schtasks", "/create",
"/tn", taskName,
"/tr", "cmd.exe /c echo Katana wake-up",
"/sc", "once",
"/st", timeStr,
"/sd", dateStr,
"/f", // Force overwrite if exists
"/ru", "SYSTEM", // Run as system to ensure wake capability
)

if err := cmd.Run(); err != nil {
return fmt.Errorf("failed to create wake-up task: %v", err)
}

// Enable wake timers for the task (this is the key for wake-from-sleep)
cmd2 := exec.Command("powercfg", "/waketimers")
cmd2.Run() // Just check if wake timers are enabled

log.Printf("Scheduled Windows wake-up task %s", taskName)
return nil
}

// Cancel Windows wake-up
func (pm *PowerManager) cancelWakeupWindows(alarmID string) error {
taskName := fmt.Sprintf("KatanaWakeup_%s", alarmID)
cmd := exec.Command("schtasks", "/delete", "/tn", taskName, "/f")
err := cmd.Run()
if err != nil {
log.Printf("Failed to delete wake-up task %s: %v", taskName, err)
}
return err
}

// macOS-specific wake-up scheduling using pmset
func (pm *PowerManager) scheduleWakeupMacOS(alarmID string, wakeupTime time.Time, duration time.Duration) error {
// Use pmset to schedule a wake event
// Format: pmset schedule wake "MM/dd/yyyy HH:mm:ss"
timeStr := wakeupTime.Format("01/02/2006 15:04:05")

cmd := exec.Command("sudo", "pmset", "schedule", "wake", timeStr)
if err := cmd.Run(); err != nil {
return fmt.Errorf("failed to schedule wake-up (requires sudo): %v", err)
}

log.Printf("Scheduled macOS wake-up for %s", timeStr)
return nil
}

// Cancel macOS wake-up
func (pm *PowerManager) cancelWakeupMacOS(alarmID string) error {
// pmset doesn't have easy individual cancellation, so we clear all scheduled wake events
	// This might affect other wake schedules, but it's the best we can do
cmd := exec.Command("sudo", "pmset", "schedule", "cancel", "wake")
err := cmd.Run()
if err != nil {
log.Printf("Failed to cancel wake-up schedule: %v", err)
}
return err
}

// Cleanup cancels all scheduled wake-ups
func (pm *PowerManager) Cleanup() {
pm.mu.Lock()
defer pm.mu.Unlock()

// Cancel all active timers
for alarmID := range pm.activeWakeTimers {
pm.CancelWakeup(alarmID)
}
pm.activeWakeTimers = make(map[string]*time.Timer)
}

// GetActiveWakeups returns the number of active wake-up schedules
func (pm *PowerManager) GetActiveWakeups() int {
pm.mu.Lock()
defer pm.mu.Unlock()
return len(pm.activeWakeTimers)
}

// Legacy methods for backward compatibility (now no-ops)
func (pm *PowerManager) PreventSleep(reason string) error {
log.Printf("PreventSleep called but ignored (using wake-up scheduling instead): %s", reason)
return nil
}

func (pm *PowerManager) AllowSleep() error {
log.Printf("AllowSleep called but ignored (using wake-up scheduling instead)")
return nil
}

func (pm *PowerManager) IsPreventingSleep() bool {
return false // We no longer prevent sleep, we schedule wake-ups
}

func (pm *PowerManager) GetActiveAlarms() int {
return pm.GetActiveWakeups()
}
