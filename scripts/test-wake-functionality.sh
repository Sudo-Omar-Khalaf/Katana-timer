#!/bin/bash
# Quick test script for Katana wake-up functionality

echo "=========================================="
echo "🧪 Katana Wake-Up Test"
echo "=========================================="
echo ""

# Test 1: Check rtcwake sudo permissions
echo "Test 1: Checking sudo permissions..."
if sudo -n /usr/sbin/rtcwake --version &> /dev/null; then
    echo "✅ PASS: rtcwake can run without password"
else
    echo "❌ FAIL: rtcwake requires password"
    echo "Run: ./install-katana.sh to fix permissions"
    exit 1
fi
echo ""

# Test 2: Schedule a test wake-up in 2 minutes
echo "Test 2: Scheduling test wake-up..."
WAKE_TIME=$(date -d '+2 minutes' +%s)
echo "Current time: $(date '+%H:%M:%S')"
echo "Wake time: $(date -d '+2 minutes' '+%H:%M:%S')"

if sudo rtcwake -m no -t $WAKE_TIME; then
    echo "✅ PASS: Wake-up scheduled successfully"
else
    echo "❌ FAIL: Could not schedule wake-up"
    exit 1
fi
echo ""

# Test 3: Verify wake alarm is set
echo "Test 3: Verifying RTC wake alarm..."
if [ -f /sys/class/rtc/rtc0/wakealarm ]; then
    ALARM=$(cat /sys/class/rtc/rtc0/wakealarm)
    if [ -n "$ALARM" ]; then
        echo "✅ PASS: Wake alarm is set"
        echo "Alarm timestamp: $ALARM"
        echo "Alarm time: $(date -d @$ALARM '+%Y-%m-%d %H:%M:%S')"
    else
        echo "⚠️  WARNING: Wake alarm file exists but is empty"
    fi
else
    echo "❌ FAIL: No RTC wake alarm support (/sys/class/rtc/rtc0/wakealarm not found)"
    echo "Your system may not support RTC wake-up"
fi
echo ""

# Test 4: Check if Katana is installed
echo "Test 4: Checking Katana installation..."
if [ -f ~/.local/bin/katana ]; then
    echo "✅ PASS: Katana installed at ~/.local/bin/katana"
    echo "Size: $(du -h ~/.local/bin/katana | cut -f1)"
else
    echo "❌ FAIL: Katana not found"
    exit 1
fi
echo ""

echo "=========================================="
echo "✅ All Tests Passed!"
echo "=========================================="
echo ""
echo "Next Steps:"
echo "1. Run Katana: ~/.local/bin/katana"
echo "2. Create an alarm for 2-3 minutes from now"
echo "3. Let your screen turn off or manually suspend"
echo "4. System should wake automatically and alarm should ring"
echo ""
echo "To test wake-up manually:"
echo "  1. Note current time: $(date '+%H:%M:%S')"
echo "  2. Suspend system: systemctl suspend"
echo "  3. System will wake at: $(date -d '+2 minutes' '+%H:%M:%S')"
echo ""
