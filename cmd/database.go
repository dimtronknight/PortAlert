package main

import (
	"fmt"
	"investment-tracker/internal/conversion"
	"investment-tracker/internal/database"
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
		fmt.Println("  go run cmd/database.go save          - Save current portfolio to database")
		fmt.Println("  go run cmd/database.go stats         - Show portfolio statistics")
		fmt.Println("  go run cmd/database.go recent <days> - Show recent portfolio history")
		fmt.Println("  go run cmd/database.go today         - Show today's portfolio data")
		fmt.Println("  go run cmd/database.go cleanup <days> - Delete data older than N days")
		os.Exit(1)
	}

	// Connect to database
	db, err := database.NewMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	portfolioService := database.NewPortfolioService(db)
	command := os.Args[1]

	switch command {
	case "save":
		err := saveCurrentPortfolio(portfolioService)
		if err != nil {
			log.Fatal("Failed to save portfolio:", err)
		}
		fmt.Println("Portfolio saved successfully!")

	case "stats":
		showPortfolioStats(portfolioService)

	case "recent":
		days := 7 // Default to 7 days
		if len(os.Args) > 2 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		showRecentHistory(portfolioService, days)

	case "today":
		showTodayData(portfolioService)

	case "cleanup":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run cmd/database.go cleanup <days_to_keep>")
			os.Exit(1)
		}
		days, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid number of days:", os.Args[2])
		}
		cleanupOldData(portfolioService, days)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func saveCurrentPortfolio(service *database.PortfolioService) error {
	// Get current portfolio values (reusing the function from notify.go)
	total, stocks, crypto, bullion, err := getPortfolioValues()
	if err != nil {
		return fmt.Errorf("failed to get portfolio values: %w", err)
	}

	return service.SaveDailySnapshot(stocks, crypto, bullion, total)
}

func showPortfolioStats(service *database.PortfolioService) {
	stats, err := service.GetPortfolioStats()
	if err != nil {
		log.Fatal("Failed to get portfolio stats:", err)
	}

	fmt.Printf("📊 Portfolio Statistics\n")
	fmt.Printf("═════════════════════════\n")
	fmt.Printf("Days tracked: %d\n", stats.DaysTracked)
	fmt.Printf("Average value: %.2f BGN\n", stats.AverageValue)
	fmt.Printf("Total growth: %.2f BGN (%.2f%%)\n", stats.TotalGrowth, stats.GrowthPercentage)
	fmt.Printf("Best day: %s (%.2f BGN)\n", stats.BestDay.Format("2006-01-02"), stats.BestDayValue)
	fmt.Printf("Worst day: %s (%.2f BGN)\n", stats.WorstDay.Format("2006-01-02"), stats.WorstDayValue)
}

func showRecentHistory(service *database.PortfolioService, days int) {
	snapshots, err := service.GetLastNDays(days)
	if err != nil {
		log.Fatal("Failed to get recent history:", err)
	}

	fmt.Printf("📈 Portfolio History (Last %d days)\n", days)
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("%-12s %-12s %-12s %-12s %-12s\n", "Date", "Trading212", "Crypto", "Bullion", "Total")
	fmt.Printf("───────────────────────────────────────────────────────────────\n")

	for _, snapshot := range snapshots {
		fmt.Printf("%-12s %-12.2f %-12.2f %-12.2f %-12.2f\n",
			snapshot.Date.Format("2006-01-02"),
			snapshot.Trading212BGN,
			snapshot.CryptoBGN,
			snapshot.BullionBGN,
			snapshot.TotalBGN)
	}
}

func showTodayData(service *database.PortfolioService) {
	snapshot, err := service.GetTodaySnapshot()
	if err != nil {
		log.Fatal("Failed to get today's data:", err)
	}

	fmt.Printf("📊 Today's Portfolio (%s)\n", snapshot.Date.Format("2006-01-02"))
	fmt.Printf("═══════════════════════════════════\n")
	fmt.Printf("🏦 Trading212: %.2f BGN\n", snapshot.Trading212BGN)
	fmt.Printf("₿ Crypto: %.2f BGN\n", snapshot.CryptoBGN)
	fmt.Printf("🥇 Bullion: %.2f BGN\n", snapshot.BullionBGN)
	fmt.Printf("💎 Total: %.2f BGN\n", snapshot.TotalBGN)
	fmt.Printf("───────────────────────────────────\n")
	fmt.Printf("Recorded at: %s\n", snapshot.CreatedAt.Format("2006-01-02 15:04:05"))
}

func cleanupOldData(service *database.PortfolioService, daysToKeep int) {
	fmt.Printf("Cleaning up data older than %d days...\n", daysToKeep)
	err := service.DeleteOldSnapshots(daysToKeep)
	if err != nil {
		log.Fatal("Failed to cleanup old data:", err)
	}
	fmt.Println("Cleanup completed!")
}

// Import the portfolio functions from notify.go
// Note: In a real application, these should be in a shared package
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