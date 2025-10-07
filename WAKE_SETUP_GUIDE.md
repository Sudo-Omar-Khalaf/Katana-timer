# 🔔 Katana Alarm - System Wake-Up Setup Guide

## Overview

Katana can wake your computer from sleep/suspend mode to ring alarms on time. This requires a one-time permission setup to allow the application to program your computer's hardware Real-Time Clock (RTC).

## Quick Setup (Linux)

Run this command in the Katana directory:

```bash
./setup-wake-permissions.sh
```

**That's it!** The script will:
- Configure sudo permissions for `rtcwake` (one-time password prompt)
- Test the configuration
- Display setup confirmation

## What This Does

The setup script creates a file at `/etc/sudoers.d/katana-alarm-wake` that allows your user to run `rtcwake` without a password prompt. This is safe because:

1. **Limited scope**: Only the `rtcwake` command is affected
2. **No general sudo**: Doesn't grant any other elevated permissions
3. **Standard practice**: Same approach used by other alarm/timer applications
4. **Removable**: Can be uninstalled anytime

## How It Works

### When You Set an Alarm:

1. **User Action**: You create an alarm for a specific time
2. **RTC Programming**: Katana uses `rtcwake` to program your hardware RTC
3. **Sleep Normally**: Your computer can sleep/suspend as usual
4. **Hardware Wake-Up**: RTC wakes the computer at alarm time
5. **Alarm Rings**: Application is running and sounds the alarm

### Technical Details:

```
User Sets Alarm (e.g., 7:00 AM)
         ↓
Katana calls: sudo rtcwake -m no -t <unix_timestamp>
         ↓
Hardware RTC is programmed with wake time
         ↓
Computer enters sleep/suspend mode
         ↓
RTC triggers wake-up at 7:00 AM (hardware level)
         ↓
Computer wakes up, Katana alarm rings
```

## Platform Support

| Platform | Method | Setup Required | Notes |
|----------|--------|----------------|-------|
| **Linux** | `rtcwake` (RTC) | Yes (one-time) | Hardware-level wake |
| **Windows** | Task Scheduler | No | Built into Windows |
| **macOS** | `pmset` | Yes (one-time) | macOS power management |

## Verification

### Test the Setup:

1. **Run the setup script**:
   ```bash
   ./setup-wake-permissions.sh
   ```

2. **Verify permissions**:
   ```bash
   sudo -n rtcwake --version
   ```
   Should display version without password prompt

3. **Test wake-up** (optional):
   ```bash
   # Schedule wake in 2 minutes
   sudo rtcwake -m no -t $(date -d '+2 minutes' +%s)
   
   # Check if alarm is set
   cat /sys/class/rtc/rtc0/wakealarm
   ```

4. **Create test alarm in Katana**:
   - Set alarm for 2-3 minutes from now
   - Let screen turn off or manually suspend
   - Computer should wake and alarm should ring

## Troubleshooting

### Alarm Doesn't Wake Computer

**Check 1: Permissions**
```bash
sudo -n rtcwake --version
```
If it asks for password, run `./setup-wake-permissions.sh` again

**Check 2: RTC Support**
```bash
ls /sys/class/rtc/
```
Should show `rtc0` or similar

**Check 3: BIOS Settings**
- Reboot into BIOS/UEFI
- Look for "Wake on RTC" or "RTC Alarm"
- Ensure it's enabled (usually enabled by default)

**Check 4: System Logs**
```bash
journalctl -xe | grep -i rtc
```
Look for errors related to RTC wake

### Virtual Machines

RTC wake may not work in virtual machines because:
- VMs don't have direct hardware RTC access
- Hypervisor may not support RTC passthrough
- **Workaround**: Don't let VM sleep, or use host-level scheduling

### Laptops with Aggressive Power Settings

Some laptops have power settings that override RTC wake:
1. Open Power Settings
2. Check "Sleep" settings
3. Ensure "Allow wake timers" is enabled
4. Disable "Hibernate" if wake issues persist

### Wayland/X11 Desktop Environments

RTC wake works at hardware level, so desktop environment doesn't matter. Both Wayland and X11 are supported.

## Uninstall

To remove the wake-up permissions:

```bash
sudo rm /etc/sudoers.d/katana-alarm-wake
```

Your alarms will still work if the computer is awake, but won't wake the computer from sleep.

## Security Considerations

### Is This Safe?

**Yes.** The setup:
- Only grants access to `rtcwake`
- Doesn't provide general sudo access
- Uses standard Linux security mechanisms
- Is easily reversible

### What Can Go Wrong?

The **worst case** scenario:
- Wake-up scheduling fails (alarm won't ring if asleep)
- No system damage possible
- No security vulnerabilities introduced

### Why Not Use Alternative Methods?

Other approaches were considered:

1. **Sleep Prevention** (systemd-inhibit):
   - ❌ Prevents computer from sleeping (defeats purpose)
   - ❌ Drains laptop battery
   - ❌ Not what user wanted

2. **User-Space Timers**:
   - ❌ Stop working when system sleeps
   - ❌ No hardware-level wake capability

3. **At Command**:
   - ❌ Can't wake system from sleep
   - ❌ Only schedules command execution

**RTC Wake** is the only proper solution for wake-from-sleep functionality.

## FAQ

**Q: Do I need to run the setup for each alarm?**
A: No, one-time setup is sufficient for all future alarms.

**Q: Will this drain my battery?**
A: No, RTC uses minimal power (same as keeping time). System sleeps normally.

**Q: Can I use Katana without this setup?**
A: Yes, alarms work fine if your computer is awake. Wake-up just won't work.

**Q: Does this work on laptops?**
A: Yes, most modern laptops support RTC wake.

**Q: What if I don't have sudo access?**
A: You won't be able to use wake-from-sleep. Alarms still work when system is awake.

**Q: Can other users abuse this?**
A: No, only your user account can use rtcwake without password.

## Support

If you encounter issues:

1. Check the **Troubleshooting** section above
2. Review `/var/log/syslog` or `journalctl` for errors
3. Open an issue on GitHub with system details

---

**Happy Waking! 🌅⏰**
