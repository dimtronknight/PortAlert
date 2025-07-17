package stocks

import (
	"fmt"
	"log"
)

func ExampleUsage() {
	client, err := NewClientFromConfig()
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	accountInfo, err := client.GetAccountInfo()
	if err != nil {
		log.Fatal("Failed to get account info:", err)
	}
	fmt.Printf("Account ID: %d, Currency: %s\n", accountInfo.ID, accountInfo.CurrencyCode)

	cash, err := client.GetAccountCash()
	if err != nil {
		log.Fatal("Failed to get account cash:", err)
	}
	fmt.Printf("Available Cash: %.2f %s\n", cash.Free, cash.CurrencyCode)

	portfolio, err := client.GetPortfolio()
	if err != nil {
		log.Fatal("Failed to get portfolio:", err)
	}

	fmt.Printf("\nPortfolio Positions:\n")
	for _, position := range portfolio {
		currentValue := position.CurrentPrice * position.Quantity
		pnl := currentValue - (position.AveragePrice * position.Quantity)
		fmt.Printf("Ticker: %s, Quantity: %.2f, Current Price: %.2f, Current Value: %.2f, P&L: %.2f\n",
			position.Ticker, position.Quantity, position.CurrentPrice, currentValue, pnl)
	}

	orders, err := client.GetOrders()
	if err != nil {
		log.Fatal("Failed to get orders:", err)
	}

	fmt.Printf("\nActive Orders:\n")
	for _, order := range orders {
		fmt.Printf("Order ID: %d, Ticker: %s, Type: %s, Status: %s, Quantity: %.2f\n",
			order.ID, order.Ticker, order.Type, order.Status, order.Quantity)
	}

	dividends, err := client.GetDividends("", 10)
	if err != nil {
		log.Fatal("Failed to get dividends:", err)
	}

	fmt.Printf("\nRecent Dividends:\n")
	for _, dividend := range dividends {
		fmt.Printf("Ticker: %s, Amount: %.2f, Date: %s\n",
			dividend.Ticker, dividend.CashAmount, dividend.PaidOn.Format("2006-01-02"))
	}
}

func CalculatePortfolioValue(client *Client) (float64, error) {
	portfolio, err := client.GetPortfolio()
	if err != nil {
		return 0, err
	}

	totalValue := 0.0
	for _, position := range portfolio {
		totalValue += position.CurrentPrice * position.Quantity
	}

	return totalValue, nil
}

func GetPortfolioSummary(client *Client) (*PortfolioSummary, error) {
	portfolio, err := client.GetPortfolio()
	if err != nil {
		return nil, err
	}

	cash, err := client.GetAccountCash()
	if err != nil {
		return nil, err
	}

	summary := &PortfolioSummary{
		TotalCash:       cash.Total,
		FreeCash:        cash.Free,
		InvestedAmount:  cash.Invested,
		TotalPositions:  len(portfolio),
		Positions:       portfolio,
		CurrencyCode:    cash.CurrencyCode,
	}

	totalValue := 0.0
	totalPnL := 0.0
	for _, position := range portfolio {
		currentValue := position.CurrentPrice * position.Quantity
		cost := position.AveragePrice * position.Quantity
		totalValue += currentValue
		totalPnL += currentValue - cost
	}

	summary.TotalValue = totalValue
	summary.TotalPnL = totalPnL
	summary.TotalPortfolioValue = totalValue + cash.Free

	return summary, nil
}

type PortfolioSummary struct {
	TotalCash            float64    `json:"totalCash"`
	FreeCash             float64    `json:"freeCash"`
	InvestedAmount       float64    `json:"investedAmount"`
	TotalValue           float64    `json:"totalValue"`
	TotalPnL             float64    `json:"totalPnL"`
	TotalPortfolioValue  float64    `json:"totalPortfolioValue"`
	TotalPositions       int        `json:"totalPositions"`
	CurrencyCode         string     `json:"currencyCode"`
	Positions            []Position `json:"positions"`
}