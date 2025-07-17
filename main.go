package main

import (
	"fmt"
	"investment-tracker/internal/conversion"
	"investment-tracker/internal/portfolio"
	"investment-tracker/internal/stocks"
	"log"
)

func main() {
	var h *portfolio.Holdings

	h, err := portfolio.LoadHoldings("/Users/ivanrusev/Documents/LP/Go/Investment-tracker/config/holdings.json")
	if err != nil {
		log.Fatal(err)
	}

	sum, err := portfolio.SumHoldingsTotal(h)
	if err != nil {
		log.Fatal(err)
	}

	// Get Trading212 stocks value (in BGN)
	stocksValueBGN, err := getTrading212Value()
	if err != nil {
		log.Printf("Warning: Could not get Trading212 stocks value: %v", err)
		stocksValueBGN = 0
	}

	// Convert crypto/bullion from USD to BGN
	sumBGN := conversion.USDToBGN(sum, 1.7346)
	
	// Add Trading212 BGN value directly (no conversion needed)
	totalBGN := sumBGN + stocksValueBGN

	fmt.Printf("\n Trading212 Stocks Value: %.2f BGN", stocksValueBGN)
	fmt.Printf("\n Total Investment Worth: %.2f BGN", totalBGN)

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
	
	fmt.Printf("\n=== Trading212 Account Debug ===\n")
	fmt.Printf("Free Cash (available to invest): %.2f %s\n", cash.Free, cash.CurrencyCode)
	fmt.Printf("Invested Amount (in stocks): %.2f %s\n", cash.Invested, cash.CurrencyCode)
	fmt.Printf("Total Cash (free + invested): %.2f %s\n", cash.Total, cash.CurrencyCode)

	// We only need the total cash - it includes both free cash and invested amount
	finalValue := cash.Total
	
	fmt.Printf("Using Total Cash as Trading212 value: %.2f %s\n", finalValue, cash.CurrencyCode)
	fmt.Printf("=== End Debug ===\n")

	if cash.CurrencyCode != "BGN" {
		log.Printf("Warning: Trading212 account currency is %s, expected BGN", cash.CurrencyCode)
	}

	return finalValue, nil
}
