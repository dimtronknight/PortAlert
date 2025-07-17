package bullion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func GetBullionPrices(bullion string, currency string) (*BullionPriceResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("BULLION_API_KEY")
	baseURL := os.Getenv("BULLION_API_URL")
	if apiKey == "" {
		return nil, fmt.Errorf("BULLION_API_KEY environment variable is not set")
	}

	apiURL := baseURL + bullion + "/" + currency

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("x-access-token", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed:: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("goldapi.io API error: HTTP %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var bullionResponse BullionPriceResponse
	err = json.Unmarshal(body, &bullionResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response %w", err)
	}

	return &bullionResponse, nil
}
