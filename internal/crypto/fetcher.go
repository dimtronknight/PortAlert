package crypto

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var cmcBaseURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"

func GetCryptoPrice(symbol string) (float64, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("CMC_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("CMC_API_KEY environment variable is not set")
	}

	req, err := http.NewRequest("GET", cmcBaseURL, nil)
	if err != nil {
		return 0, err
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed:: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("CoinMarketCap API error: HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format (data)")
	}

	coin, ok := data[symbol].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("coin %s not found", symbol)
	}

	quote, ok := coin["quote"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("missing quote")
	}

	usd, ok := quote["USD"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("missing USD price")
	}

	price, ok := usd["price"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid price format")
	}

	return price, nil

}
