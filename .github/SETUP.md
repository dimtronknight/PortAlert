# GitHub Actions Setup Guide

This guide will help you set up GitHub Actions for your Investment Tracker project.

## Prerequisites

1. Push your code to a GitHub repository
2. Have API keys for the services you want to use

## Step 1: Set up Repository Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions

Add the following secrets:

### Required Secrets

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `TRADING212_API_KEY` | Your Trading212 API key | `your_trading212_api_key` |
| `TRADING212_IS_LIVE` | Set to "true" for live account, "false" for demo | `false` |
| `CMC_API_KEY` | CoinMarketCap API key for crypto prices | `your_cmc_api_key` |
| `BULLION_API_KEY` | Bullion API key for gold/silver prices | `your_bullion_api_key` |
| `BULLION_API_URL` | Bullion API base URL | `https://api.metals.live/v1/spot/` |

### Optional Secrets (for notifications)

| Secret Name | Description |
|-------------|-------------|
| `TELEGRAM_BOT_TOKEN` | Telegram bot token for notifications |
| `TELEGRAM_CHAT_ID` | Telegram chat ID to send notifications to |
| `SMTP_HOST` | SMTP server host for email notifications |
| `SMTP_PORT` | SMTP server port |
| `SMTP_USERNAME` | SMTP username |
| `SMTP_PASSWORD` | SMTP password |
| `EMAIL_FROM` | Email address to send from |
| `EMAIL_TO` | Email address to send to |

## Step 2: Understanding the Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)
- **Triggers**: Push to main/develop, PRs to main/develop
- **Actions**: 
  - Runs tests
  - Builds the application
  - Checks code formatting
  - Creates build artifacts for multiple platforms

### 2. Deploy Workflow (`.github/workflows/deploy.yml`)
- **Triggers**: 
  - Push to main branch
  - Daily at 9 AM UTC (scheduled)
  - Manual trigger via GitHub Actions UI
- **Actions**:
  - Builds and runs the investment tracker
  - Uses all your API keys to fetch current portfolio value
  - Uploads logs as artifacts

### 3. Release Workflow (`.github/workflows/release.yml`)
- **Triggers**: When you push a git tag starting with 'v' (e.g., v1.0.0)
- **Actions**:
  - Creates binaries for multiple platforms
  - Creates a GitHub release with downloadable files
  - Generates changelog automatically

## Step 3: Set up API Keys

### Trading212 API Key
1. Log into Trading212 web platform
2. Go to Invest account → Practice mode
3. Settings → API (Beta)
4. Generate API key
5. Add to GitHub secrets as `TRADING212_API_KEY`

### CoinMarketCap API Key
1. Visit https://coinmarketcap.com/api/
2. Sign up for free account
3. Generate API key
4. Add to GitHub secrets as `CMC_API_KEY`

### Bullion API Key
1. Visit your bullion price provider (e.g., metals.live)
2. Sign up and get API key
3. Add to GitHub secrets as `BULLION_API_KEY`

## Step 4: Test the Setup

1. **Push code to GitHub**: This will trigger the CI workflow
2. **Manual deployment**: Go to Actions → Deploy → Run workflow
3. **Create a release**: 
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Step 5: Monitor Your Investments

### Automated Daily Reports
- The deploy workflow runs daily at 9 AM UTC
- It calculates your total portfolio value
- Logs are saved as artifacts in GitHub Actions

### Manual Runs
- You can trigger the deployment manually anytime
- Go to Actions → Deploy → Run workflow
- Choose environment (production/staging)

## Environment Management

The workflows support different environments:
- **Production**: Uses live API keys and real data
- **Staging**: Can use demo/test API keys

To use different environments:
1. Go to Settings → Environments
2. Create "production" and "staging" environments
3. Add environment-specific secrets

## Troubleshooting

### Common Issues

1. **API Key Issues**: Make sure all required secrets are set
2. **Build Failures**: Check Go version compatibility
3. **Permission Issues**: Ensure repository has proper permissions

### Viewing Logs
- Go to Actions tab in your repository
- Click on any workflow run
- View logs for each step
- Download artifacts for detailed logs

## Security Best Practices

1. Never commit API keys to your repository
2. Use environment-specific secrets for different deployments
3. Regularly rotate your API keys
4. Review workflow permissions

## Customization

### Changing Schedule
Edit `.github/workflows/deploy.yml` and modify the cron expression:
```yaml
schedule:
  - cron: '0 9 * * *'  # Daily at 9 AM UTC
```

### Adding More Platforms
Add more GOOS/GOARCH combinations in the build steps:
```yaml
GOOS=linux GOARCH=arm64 go build -o bin/investment-tracker-linux-arm64 .
```

### Notifications
The workflows can be extended to send notifications on success/failure by adding steps that use your notification secrets.