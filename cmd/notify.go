package main

import (
	"fmt"
	"investment-tracker/internal/conversion"
	"investment-tracker/internal/database"
	"investment-tracker/internal/notifications"
	"investment-tracker/internal/portfolio"
	"investment-tracker/internal/stocks"
	"log"
	"os"
	"strconv"
	
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/notify.go test           - Send test notification")
		fmt.Println("  go run cmd/notify.go now            - Send notification now")
		fmt.Println("  go run cmd/notify.go schedule 8 30  - Schedule daily at 8:30 AM")
		os.Exit(1)
	}

	command := os.Args[1]
	
	notificationService := notifications.NewNotificationService()
	scheduler := notifications.NewScheduler(notificationService)

	switch command {
	case "test":
		fmt.Println("Sending test notification...")
		err := notificationService.TestNotifications()
		if err != nil {
			log.Fatal("Test notification failed:", err)
		}
		fmt.Println("Test notification sent successfully!")

	case "now":
		fmt.Println("Sending current portfolio notification...")
		err := scheduler.SendNow(getPortfolioValues)
		if err != nil {
			log.Fatal("Failed to send notification:", err)
		}
		
		// Save to database
		err = savePortfolioToDatabase()
		if err != nil {
			log.Printf("Warning: Failed to save to database: %v", err)
		}
		
		fmt.Println("Notification sent successfully!")

	case "schedule":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run cmd/notify.go schedule <hour> <minute>")
			fmt.Println("Example: go run cmd/notify.go schedule 8 30")
			os.Exit(1)
		}

		hour, err := strconv.Atoi(os.Args[2])
		if err != nil || hour < 0 || hour > 23 {
			log.Fatal("Invalid hour (must be 0-23):", os.Args[2])
		}

		minute, err := strconv.Atoi(os.Args[3])
		if err != nil || minute < 0 || minute > 59 {
			log.Fatal("Invalid minute (must be 0-59):", os.Args[3])
		}

		fmt.Printf("Starting daily notification scheduler for %02d:%02d...\n", hour, minute)
		
		// Set up database callback for scheduled notifications
		scheduler.SetOnNotificationSent(func(total, stocks, crypto, bullion float64) {
			err := savePortfolioToDatabase()
			if err != nil {
				log.Printf("Error saving scheduled notification to database: %v", err)
			}
		})
		
		scheduler.Start(hour, minute, getPortfolioValues)

		// Keep the program running
		select {}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func getPortfolioValues() (total, stocks, crypto, bullion float64, err error) {
	// Load existing holdings
	h, err := portfolio.LoadHoldings("/Users/ivanrusev/Documents/LP/Go/Investment-tracker/config/holdings.json")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to load holdings: %w", err)
	}

	// Get crypto value (in USD)
	cryptoUSD, err := portfolio.SumCryptoHoldings(h)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to sum crypto holdings: %w", err)
	}

	// Get bullion value (in USD)
	bullionUSD, err := portfolio.SumBullionHoldings(h)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to sum bullion holdings: %w", err)
	}

	// Get Trading212 stocks value (in BGN)
	stocksValueBGN, err := getTrading212Value()
	if err != nil {
		log.Printf("Warning: Could not get Trading212 stocks value: %v", err)
		stocksValueBGN = 0
	}

	// Convert crypto and bullion from USD to BGN
	cryptoBGN := conversion.USDToBGN(cryptoUSD, 1.7346)
	bullionBGN := conversion.USDToBGN(bullionUSD, 1.7346)
	
	// Total = crypto (BGN) + bullion (BGN) + stocks (BGN)
	totalBGN := cryptoBGN + bullionBGN + stocksValueBGN

	return totalBGN, stocksValueBGN, cryptoBGN, bullionBGN, nil
}

func getTrading212Value() (float64, error) {
	client, err := stocks.NewClientFromConfig()
	if err != nil {
		return 0, fmt.Errorf("failed to create Trading212 client: %w", err)
	}

	// Get account cash - this contains all the value we need
	cash, err := client.GetAccountCash()
	if err != nil {
		return 0, fmt.Errorf("failed to get account cash: %w", err)
	}

	// We only need the total cash - it includes both free cash and invested amount
	return cash.Total, nil
}

func savePortfolioToDatabase() error {
	// Get portfolio values
	total, stocks, crypto, bullion, err := getPortfolioValues()
	if err != nil {
		return fmt.Errorf("failed to get portfolio values: %w", err)
	}

	// Connect to database
	db, err := database.NewMongoDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Save snapshot
	portfolioService := database.NewPortfolioService(db)
	err = portfolioService.SaveDailySnapshot(stocks, crypto, bullion, total)
	if err != nil {
		return fmt.Errorf("failed to save portfolio snapshot: %w", err)
	}

	log.Printf("Successfully saved portfolio snapshot to database")
	return nil
}