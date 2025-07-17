# Investment Tracker 💰

A comprehensive investment portfolio tracker that monitors your Trading212 stocks, cryptocurrency, and bullion holdings. Features real-time portfolio tracking, daily notifications, and historical data storage with MongoDB.

## Features ✨

- **Real-time Portfolio Tracking**: Monitor Trading212, crypto, and bullion values
- **Daily Notifications**: Telegram, SMS, and email notifications
- **Historical Data**: MongoDB integration for portfolio history and analytics
- **Multi-Currency Support**: USD to BGN conversion
- **Automated Scheduling**: Set up daily portfolio updates

## Setup 🚀

### Prerequisites
- Go 1.22.5 or higher
- MongoDB (local or Atlas)
- Trading212 API key
- Notification service credentials (Telegram/SMS/Email)

### Installation

1. **Clone and setup:**
```bash
git clone <repository-url>
cd Investment-tracker
go mod tidy
```

2. **Configure environment variables:**
```bash
cp .env.example .env
# Edit .env with your credentials
```

3. **Install MongoDB:**
```bash
# macOS
brew install mongodb-community
brew services start mongodb-community

# Or use MongoDB Atlas (cloud)
```

## Configuration 🔧

### Environment Variables (.env)

```bash
# Trading212 API Configuration
TRADING212_API_KEY=your_trading212_api_key_here
TRADING212_IS_LIVE=false  # true for live, false for demo

# Notification Configuration
NOTIFICATION_METHODS=telegram  # telegram, sms, email (comma-separated)

# Telegram Configuration
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_CHAT_ID=your_telegram_chat_id

# SMS Configuration (Twilio)
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIO_FROM_NUMBER=+1234567890
TWILIO_TO_NUMBER=+1234567890

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
FROM_EMAIL=your_email@gmail.com
TO_EMAIL=your_email@gmail.com

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=investment_tracker
```

### Trading212 API Setup

1. Go to Trading212 Settings → API
2. Create API key with permissions:
   - `account` - Account information
   - `portfolio` - Portfolio data
   - `orders:read` - Order history
   - `history` - Historical data

### Telegram Bot Setup

1. Message @BotFather on Telegram
2. Send `/newbot` and follow instructions
3. Get your bot token
4. Send a message to your bot
5. Get your chat ID from: `https://api.telegram.org/bot<TOKEN>/getUpdates`

## Commands 🎯

### Portfolio & Notifications

#### Main Investment Tracker
```bash
# Run main portfolio calculation
go run main.go
```

#### Notification Commands
```bash
# Send test notification (dummy data)
go run cmd/notify.go test

# Send notification with current portfolio values
go run cmd/notify.go now

# Schedule daily notifications at specific time
go run cmd/notify.go schedule 8 30    # 8:30 AM daily
go run cmd/notify.go schedule 18 00   # 6:00 PM daily
```

### Database Management

#### Save Data
```bash
# Save current portfolio snapshot to database
go run cmd/database.go save
```

#### View Data
```bash
# View today's portfolio data
go run cmd/database.go today

# View recent portfolio history
go run cmd/database.go recent 7     # Last 7 days
go run cmd/database.go recent 30    # Last 30 days
go run cmd/database.go recent 90    # Last 90 days

# View portfolio statistics
go run cmd/database.go stats
```

#### Data Maintenance
```bash
# Clean up old data (keep last N days)
go run cmd/database.go cleanup 365  # Keep last 365 days
go run cmd/database.go cleanup 90   # Keep last 90 days
```

## Usage Examples 📋

### Daily Routine
```bash
# Check current portfolio and get notification
go run cmd/notify.go now

# View today's saved data
go run cmd/database.go today

# Check portfolio statistics
go run cmd/database.go stats
```

### Weekly Analysis
```bash
# View last week's performance
go run cmd/database.go recent 7

# View last month's performance
go run cmd/database.go recent 30
```

### Setup Automation
```bash
# Set up daily 8:30 AM notifications
go run cmd/notify.go schedule 8 30

# Or use cron for automation
# Add to crontab: 30 8 * * * cd /path/to/project && go run cmd/notify.go now
```

## Sample Outputs 📊

### Telegram Notification
```
💰 Daily Investment Update

🏦 Trading212: 3,905.04 BGN
₿ Crypto: 873.87 BGN
🥇 Bullion: 8,855.50 BGN

💎 Total: 13,634.41 BGN

📈 Have a great day!
```

### Database Statistics
```
📊 Portfolio Statistics
═════════════════════════
Days tracked: 45
Average value: 12,450.00 BGN
Total growth: 1,184.41 BGN (9.5%)
Best day: 2024-03-15 (14,200.00 BGN)
Worst day: 2024-02-28 (11,890.00 BGN)
```

### Recent History
```
📈 Portfolio History (Last 7 days)
═══════════════════════════════════════════════════════════════
Date         Trading212   Crypto       Bullion      Total       
───────────────────────────────────────────────────────────────
2024-03-10   3800.00      950.00       8500.00      13250.00     
2024-03-11   3850.00      920.00       8600.00      13370.00     
2024-03-12   3900.00      890.00       8700.00      13490.00     
2024-03-13   3950.00      870.00       8750.00      13570.00     
2024-03-14   3905.04      873.87       8855.50      13634.41     
```

## File Structure 📁

```
Investment-tracker/
├── cmd/
│   ├── notify.go          # Notification commands
│   └── database.go        # Database management
├── config/
│   └── holdings.json      # Portfolio configuration
├── internal/
│   ├── bullion/          # Bullion price fetching
│   ├── conversion/       # Currency conversion
│   ├── crypto/           # Cryptocurrency data
│   ├── database/         # MongoDB integration
│   ├── notifications/    # Notification services
│   ├── portfolio/        # Portfolio calculations
│   └── stocks/           # Trading212 integration
├── main.go               # Main portfolio tracker
├── .env.example          # Environment template
└── README.md            # This file
```

## API Integrations 🔌

- **Trading212 API**: Real-time stock portfolio data
- **Cryptocurrency APIs**: Current crypto prices
- **Bullion APIs**: Precious metals pricing
- **Telegram Bot API**: Push notifications
- **Twilio API**: SMS notifications
- **SMTP**: Email notifications

## Database Schema 📚

### PortfolioSnapshot
- `date`: Daily snapshot date
- `trading212_bgn`: Trading212 value in BGN
- `crypto_bgn`: Crypto value in BGN
- `bullion_bgn`: Bullion value in BGN
- `total_bgn`: Total portfolio value
- `created_at`: Timestamp
- `usd_to_bgn_rate`: Exchange rate used

## Troubleshooting 🔧

### Common Issues

1. **MongoDB Connection Error**
   - Ensure MongoDB is running: `brew services start mongodb-community`
   - Check MONGODB_URI in .env

2. **Trading212 API Error**
   - Verify API key permissions
   - Check if using correct environment (live vs demo)

3. **Notification Failures**
   - Verify Telegram bot token and chat ID
   - Check network connectivity

4. **Missing Holdings**
   - Ensure `config/holdings.json` exists and is properly formatted
   - Check crypto/bullion API keys

### Debug Commands
```bash
# Test individual components
go run cmd/notify.go test
go run cmd/database.go today
go run main.go
```

## Contributing 🤝

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License 📄

This project is licensed under the MIT License.

---

**Happy Investing! 📈💰**