# GitHub Repository Secrets Setup Guide

This guide shows you exactly how to set up repository secrets for your investment tracker.

## 🔑 Step 1: Access Repository Secrets

1. Go to your GitHub repository
2. Click **Settings** (top menu)
3. Click **Secrets and variables** (left sidebar)
4. Click **Actions**
5. Click **New repository secret**

## 📝 Step 2: Required Secrets

Add these secrets one by one:

### Trading212 API
```
Name: TRADING212_API_KEY
Value: your_trading212_api_key_here
```

```
Name: TRADING212_IS_LIVE
Value: false
```
> Set to `true` for live account, `false` for demo account

### CoinMarketCap API
```
Name: CMC_API_KEY
Value: your_coinmarketcap_api_key_here
```

### Bullion API
```
Name: BULLION_API_KEY
Value: your_bullion_api_key_here
```

```
Name: BULLION_API_URL
Value: https://api.metals.live/v1/spot/
```

### Telegram Bot (for notifications)
```
Name: TELEGRAM_BOT_TOKEN
Value: your_telegram_bot_token_here
```

```
Name: TELEGRAM_CHAT_ID
Value: your_telegram_chat_id_here
```

## 🔧 Step 3: Get Your API Keys

### Trading212 API Key
1. Log into Trading212 web platform
2. Switch to **Practice Mode** (required for API access)
3. Go to **Settings** → **API (Beta)**
4. Generate API key
5. Copy the key and add it to GitHub secrets

### CoinMarketCap API Key
1. Visit https://coinmarketcap.com/api/
2. Sign up for free account
3. Go to **API Keys** section
4. Generate new API key
5. Copy the key and add it to GitHub secrets

### Bullion API Key
1. Visit https://metals.live/ (or your preferred bullion API provider)
2. Sign up for account
3. Generate API key
4. Copy the key and add it to GitHub secrets

### Telegram Bot Token
1. Message `@BotFather` on Telegram
2. Type `/newbot`
3. Follow instructions to create your bot
4. Copy the bot token and add it to GitHub secrets

### Telegram Chat ID
1. Message your bot on Telegram
2. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
3. Look for your chat ID in the response
4. Copy the chat ID and add it to GitHub secrets

## 🧪 Step 4: Test Your Setup

1. **Push your code** to trigger the workflow:
   ```bash
   git add .
   git commit -m "Fix environment variable handling"
   git push origin main
   ```

2. **Check the workflow** in the Actions tab
3. **Verify the output** shows your API keys (first 10 characters)
4. **Check your Telegram** for the notification

## 🔍 Step 5: Verify Secrets Are Set

After adding all secrets, you should see them listed in your repository secrets:

- ✅ TRADING212_API_KEY
- ✅ TRADING212_IS_LIVE  
- ✅ CMC_API_KEY
- ✅ BULLION_API_KEY
- ✅ BULLION_API_URL
- ✅ TELEGRAM_BOT_TOKEN
- ✅ TELEGRAM_CHAT_ID

## 🚨 Troubleshooting

### Workflow still failing?
- Check that all required secrets are added
- Verify API keys are correct
- Check the Actions tab for detailed error messages

### Not receiving Telegram notifications?
- Verify bot token is correct
- Check chat ID is correct
- Make sure you've messaged the bot first
- Check Telegram notification settings

### API errors?
- Verify API keys have proper permissions
- Check API rate limits
- Ensure Trading212 is in Practice Mode

## 🎯 Expected Workflow Output

When everything is set up correctly, you should see:

```
🚀 Starting daily portfolio tracking...
Environment variables check:
- CMC_API_KEY: 1234567890...
- TRADING212_API_KEY: abcdef1234...
- TELEGRAM_BOT_TOKEN: 1234567890...
Sending current portfolio notification...
Info: No .env file found, using environment variables: open .env: no such file or directory
2025/07/18 10:15:16 Sending immediate investment notification...
2025/07/18 10:15:16 Successfully sent Telegram notification
Notification sent successfully!
Investment tracker completed successfully
```

## 📱 What You'll Receive

A Telegram message like:
```
*Daily Investment Update*

- Trading212: 	2000.00 BGN
- Crypto: 		1500.00 BGN
- Bullion: 	1000.00 BGN

*Total: 4500.00 BGN*
```

## 🔄 Automatic Schedule

The workflow runs:
- **Daily at 9 AM UTC** (automatic)
- **On every push to main** (automatic)  
- **Manually triggered** (from Actions tab)

You can change the schedule by editing `.github/workflows/deploy.yml`