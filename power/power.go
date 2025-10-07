package power

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// PowerManager handles system power management and wake-up scheduling
type PowerManager struct {
	mu               sync.Mutex
	activeWakeTimers map[string]*time.Timer
}

// NewPowerManager creates a new power manager instance
func NewPowerManager() *PowerManager {
	return &PowerManager{
		activeWakeTimers: make(map[string]*time.Timer),
	}
}

// ScheduleWakeup schedules the system to wake up at a specific time
// This sets a system-level wake timer that will wake the PC from sleep
func (pm *PowerManager) ScheduleWakeup(alarmID string, wakeTime time.Time) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Cancel any existing wake timer for this alarm
	pm.cancelWakeupLocked(alarmID)

	duration := time.Until(wakeTime)
	if duration <= 0 {
		return fmt.Errorf("wake time is in the past")
	}

	log.Printf("Scheduling system wakeup for alarm %s in %v at %v", alarmID, duration, wakeTime)

	switch runtime.GOOS {
	case "linux":
		return pm.scheduleLinuxWakeup(alarmID, wakeTime)
	case "windows":
		return pm.scheduleWindowsWakeup(alarmID, wakeTime)
	case "darwin":
		return pm.scheduleMacOSWakeup(alarmID, wakeTime)
	default:
		log.Printf("System wake-up scheduling not supported on %s, alarm will still work if system is awake", runtime.GOOS)
		return nil
	}
}

// CancelWakeup cancels a scheduled wake-up timer
func (pm *PowerManager) CancelWakeup(alarmID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.cancelWakeupLocked(alarmID)
}

// cancelWakeupLocked cancels a wake-up timer (must be called with mutex locked)
func (pm *PowerManager) cancelWakeupLocked(alarmID string) {
	if timer, exists := pm.activeWakeTimers[alarmID]; exists {
		timer.Stop()
		delete(pm.activeWakeTimers, alarmID)
		log.Printf("Cancelled wake timer for alarm %s", alarmID)
	}

	switch runtime.GOOS {
	case "linux":
		pm.cancelLinuxWakeup(alarmID)
	case "windows":
		pm.cancelWindowsWakeup(alarmID)
	case "darwin":
		pm.cancelMacOSWakeup(alarmID)
	}
}

// scheduleLinuxWakeup schedules wake-up on Linux using rtcwake
func (pm *PowerManager) scheduleLinuxWakeup(alarmID string, wakeTime time.Time) error {
	timestamp := wakeTime.Unix()

	cmd := exec.Command("sudo", "rtcwake", "-m", "no", "-t", fmt.Sprintf("%d", timestamp))
	if err := cmd.Run(); err != nil {
		log.Printf("rtcwake failed: %v", err)
		if err := pm.scheduleLinuxAtCommand(alarmID, wakeTime); err != nil {
			log.Printf("Warning: Could not schedule system wake-up. Alarm will only ring if system is awake.")
			return nil
		}
		return nil
	}

	log.Printf("Linux wake-up scheduled using rtcwake for %v", wakeTime)
	return nil
}

// scheduleLinuxAtCommand uses the 'at' command as fallback
func (pm *PowerManager) scheduleLinuxAtCommand(alarmID string, wakeTime time.Time) error {
	if _, err := exec.LookPath("at"); err != nil {
		return fmt.Errorf("'at' command not available")
	}

	atTime := wakeTime.Format("15:04 01/02/2006")
	wakeCommand := fmt.Sprintf("echo 'Katana alarm %s wake-up' | wall", alarmID)

	cmd := exec.Command("bash", "-c", fmt.Sprintf(`echo "%s" | at %s`, wakeCommand, atTime))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to schedule with 'at' command: %v", err)
	}

	log.Printf("Linux wake-up scheduled using 'at' command for %v", wakeTime)
	return nil
}

// cancelLinuxWakeup cancels Linux wake-up timers
func (pm *PowerManager) cancelLinuxWakeup(alarmID string) {
	exec.Command("sudo", "rtcwake", "-m", "disable").Run()
	cmd := exec.Command("bash", "-c", fmt.Sprintf("atq | grep 'Katana alarm %s' | cut -f1 | xargs -r atrm", alarmID))
	cmd.Run()
}

// scheduleWindowsWakeup schedules wake-up on Windows using Task Scheduler
func (pm *PowerManager) scheduleWindowsWakeup(alarmID string, wakeTime time.Time) error {
	taskName := fmt.Sprintf("KatanaAlarm_%s", alarmID)
	timeStr := wakeTime.Format("15:04")
	dateStr := wakeTime.Format("01/02/2006")

	args := []string{
		"/create", "/tn", taskName,
		"/tr", "echo Katana Alarm Wake-up",
		"/sc", "once",
		"/st", timeStr,
		"/sd", dateStr,
		"/ru", "SYSTEM",
		"/rl", "HIGHEST",
		"/f",
	}

	cmd := exec.Command("schtasks", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Windows wake timer: %v", err)
	}

	log.Printf("Windows wake-up scheduled for %v", wakeTime)
	return nil
}

// cancelWindowsWakeup cancels Windows wake-up timers
func (pm *PowerManager) cancelWindowsWakeup(alarmID string) {
	taskName := fmt.Sprintf("KatanaAlarm_%s", alarmID)
	cmd := exec.Command("schtasks", "/delete", "/tn", taskName, "/f")
	cmd.Run()
}

// scheduleMacOSWakeup schedules wake-up on macOS using pmset
func (pm *PowerManager) scheduleMacOSWakeup(alarmID string, wakeTime time.Time) error {
	timeStr := wakeTime.Format("01/02/06 15:04:05")

	cmd := exec.Command("sudo", "pmset", "schedule", "wake", timeStr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to schedule macOS wake-up: %v", err)
	}

	log.Printf("macOS wake-up scheduled for %v", wakeTime)
	return nil
}

// cancelMacOSWakeup cancels macOS wake-up timers
func (pm *PowerManager) cancelMacOSWakeup(alarmID string) {
	cmd := exec.Command("sudo", "pmset", "schedule", "cancel")
	cmd.Run()
}

// AllowSleep allows the system to sleep (cancels wake-up scheduling)
func (pm *PowerManager) AllowSleep(alarmID string) {
	pm.CancelWakeup(alarmID)
}

// GetActiveWakeTimers returns the list of active wake timer IDs
func (pm *PowerManager) GetActiveWakeTimers() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	timers := make([]string, 0, len(pm.activeWakeTimers))
	for id := range pm.activeWakeTimers {
		timers = append(timers, id)
	}
	return timers
}

// Cleanup cleans up any resources used by the power manager
func (pm *PowerManager) Cleanup() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for alarmID := range pm.activeWakeTimers {
		pm.cancelWakeupLocked(alarmID)
	}
}
