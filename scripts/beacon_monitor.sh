#!/bin/bash
# Beacon Auto-Submitter Monitor
# Run this to check progress: bash beacon_monitor.sh

echo "========================================"
echo "  BEACON AUTO-SUBMITTER MONITOR"
echo "========================================"
echo ""

# Check if process is running
if ps aux | grep beacon_auto_submit | grep -v grep > /dev/null; then
    echo "✓ Process: RUNNING"
    ps aux | grep beacon_auto | grep -v grep | awk '{print "  PID: "$2", CPU: "$3"%", Mem: "$4"%"}'
else
    echo "✗ Process: NOT RUNNING"
fi
echo ""

# Show state
echo "=== Current State ==="
if [ -f /tmp/beacon_state.json ]; then
    cat /tmp/beacon_state.json | jq '{
        agent_id,
        registered,
        heartbeat_count: (.heartbeats | length),
        created_at: (.created_at | todate)
    }'
else
    echo "No state file yet"
fi
echo ""

# Show last 20 lines of output
echo "=== Latest Output ==="
tail -20 /tmp/beacon_output.log 2>/dev/null || echo "No output yet"
echo ""

# Estimate completion
if [ -f /tmp/beacon_state.json ]; then
    registered=$(cat /tmp/beacon_state.json | jq -r '.registered')
    heartbeat_count=$(cat /tmp/beacon_state.json | jq -r '.heartbeats | length')
    
    echo "=== Progress ==="
    if [ "$registered" = "true" ]; then
        echo "Registration: ✓ COMPLETE"
        echo "Heartbeats: $heartbeat_count/3"
        
        if [ "$heartbeat_count" -ge 3 ]; then
            echo ""
            echo "🎉 SUBMISSION COMPLETE!"
            echo "Claim comment: /tmp/beacon_claim.md"
        else
            remaining=$((3 - heartbeat_count))
            eta=$((remaining * 5))
            echo "ETA: ~$eta minutes"
        fi
    else
        echo "Registration: ⏳ IN PROGRESS (retrying...)"
    fi
fi
