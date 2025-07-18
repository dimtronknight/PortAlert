# Telegram Bot Notifications Troubleshooting Guide

## Common Issues and Solutions

### 1. **Messages Arrive But No Notifications**

This is the most common issue. The bot is working, but you're not getting push notifications.

#### Check Device Settings:
- **iPhone**: Settings → Notifications → Telegram → Allow Notifications
- **Android**: Settings → Apps → Telegram → Notifications → Turn on

#### Check Telegram App Settings:
1. Open Telegram app
2. Go to Settings → Notifications and Sounds
3. Make sure notifications are enabled for:
   - Private Chats
   - Bots (if there's a separate setting)

#### Check Bot-Specific Settings:
1. Go to your chat with the bot
2. Tap the bot name at the top
3. Check if notifications are muted
4. Unmute if necessary

### 2. **Test Your Setup**

Run the test program to verify your bot configuration:

```bash
go run cmd/test-telegram/main.go
```

This will:
- Test bot connection
- Get chat information
- Send a test notification
- Provide troubleshooting tips

### 3. **Environment Variables**

Make sure these are set correctly:
```bash
export TELEGRAM_BOT_TOKEN="your_bot_token"
export TELEGRAM_CHAT_ID="your_chat_id"
```

### 4. **Getting Your Chat ID**

If you're unsure about your chat ID:

1. **Method 1**: Message your bot, then visit:
   ```
   https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates
   ```

2. **Method 2**: Add `@userinfobot` to get your user ID

3. **Method 3**: Use the test program - it will show chat info

### 5. **Bot Setup Issues**

#### Create a New Bot:
1. Message `@BotFather` on Telegram
2. Type `/newbot`
3. Follow the instructions
4. Save the bot token

#### Set Bot Commands (Optional):
Message `@BotFather`:
```
/setcommands
```
Then select your bot and add:
```
start - Start the bot
help - Get help
status - Check investment status
```

### 6. **Notification Timing Issues**

If notifications work sometimes but not others:

#### Check Rate Limits:
- Telegram has rate limits (30 messages/second to different chats)
- Your bot shouldn't hit this limit, but worth checking

#### Check Bot Status:
- Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getMe`
- Should return bot information

### 7. **Advanced Troubleshooting**

#### Enable Debug Mode:
Add this to your notification code:
```go
fmt.Printf("Sending message to chat %s\n", chatID)
fmt.Printf("Message: %s\n", message)
```

#### Check Raw Response:
The updated code now shows full API responses on errors.

#### Test with Different Message Types:
```go
// Test with different parse modes
telegramMsg := TelegramMessage{
    ChatID:    t.ChatID,
    Text:      "Test message",
    ParseMode: "HTML", // or "Markdown" or ""
}
```

### 8. **Phone/Device Specific Issues**

#### iPhone:
- Check Focus/Do Not Disturb settings
- Ensure Telegram has notification permissions
- Check notification scheduling settings

#### Android:
- Check battery optimization settings
- Ensure Telegram is not in "doze" mode
- Check notification channels

### 9. **Quick Fixes to Try**

1. **Restart Telegram app**
2. **Log out and log back in to Telegram**
3. **Clear Telegram app cache** (Android)
4. **Update Telegram app**
5. **Test with a different device**

### 10. **Verify Your Bot is Working**

Run this test:
```bash
# Test the bot directly
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage" \
-H "Content-Type: application/json" \
-d '{
  "chat_id": "<YOUR_CHAT_ID>",
  "text": "Test notification from curl"
}'
```

If this works but your Go code doesn't, there's an issue with your Go implementation.

### 11. **Still Not Working?**

1. **Check bot logs** for any errors
2. **Try creating a new bot** with BotFather
3. **Test with a group chat** instead of private chat
4. **Contact Telegram support** if the issue persists

---

## Testing Commands

```bash
# Test your setup
go run cmd/test-telegram/main.go

# Build and run main program
go build -o investment-tracker .
./investment-tracker

# Test with curl
curl -X POST "https://api.telegram.org/bot<TOKEN>/sendMessage" \
-H "Content-Type: application/json" \
-d '{"chat_id": "<CHAT_ID>", "text": "Test"}'
```

## Environment Setup

Create a `.env` file:
```
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
```

Or export environment variables:
```bash
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export TELEGRAM_CHAT_ID="your_chat_id_here"
```