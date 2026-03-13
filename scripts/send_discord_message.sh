#!/bin/bash
# Quick script to send Discord message to @geekster007

set -e

echo "📧 Discord Message Sender for FinMind #144"
echo "══════════════════════════════════════════"
echo ""

# Check if discord.py is installed
if ! python3 -c "import discord" 2>/dev/null; then
    echo "📦 Installing discord.py..."
    pip install discord.py
fi

# Check for token
if [ -z "$DISCORD_BOT_TOKEN" ]; then
    echo "❌ DISCORD_BOT_TOKEN not set!"
    echo ""
    echo "To get your bot token:"
    echo "1. Go to https://discord.com/developers/applications"
    echo "2. Create a new application"
    echo "3. Go to 'Bot' tab → Add Bot"
    echo "4. Copy the token"
    echo ""
    echo "Then run: export DISCORD_BOT_TOKEN='your-token'"
    exit 1
fi

# Check for recipient ID
if [ -z "$RECIPIENT_USER_ID" ]; then
    echo "❌ RECIPIENT_USER_ID not set!"
    echo ""
    echo "To get @geekster007's User ID:"
    echo "1. Enable Developer Mode in Discord"
    echo "2. Right-click @geekster007 → Copy ID"
    echo ""
    echo "Then run: export RECIPIENT_USER_ID='copied-id'"
    exit 1
fi

echo "✅ Running Discord message sender..."
echo ""

# Run the Python script
python3 /home/administrator/projects/bountyOS/scripts/send_discord_dm.py

echo ""
echo "══════════════════════════════════════════"
echo "✅ Done! Check Discord for confirmation."
