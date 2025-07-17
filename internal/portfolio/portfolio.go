package portfolio

import (
	"encoding/json"
	"fmt"
	"investment-tracker/internal/bullion"
	"investment-tracker/internal/crypto"
	"io/ioutil"
)

// func sumCryptoHoldings()

type Holdings struct {
	Crypto  map[string]float64 `json:"crypto"`
	Bullion map[string]float64 `json:"bullion"`
}

func LoadHoldings(path string) (*Holdings, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var holdings Holdings

	err = json.Unmarshal([]byte(bytes), &holdings)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return &holdings, nil
}

func SumHoldingsTotal(h *Holdings) (float64, error) {

	var sum float64

	for symbol, amount := range h.Crypto {
		price, err := crypto.GetCryptoPrice(symbol)

		if err != nil {
			return 0, fmt.Errorf("could not sum crypto networth %s", err)
		}

		totalPricePerCrypto := amount * price

		sum += totalPricePerCrypto

	}

	fmt.Printf("Crypto Networth: %v", sum*1.7346)

	for metal, amount := range h.Bullion {

		priceResponse, err := bullion.GetBullionPrices(metal, "USD")
		if err != nil {
			return 0, fmt.Errorf("could not sum bullion networth %s", err)
		}

		sum += amount * priceResponse.Price
	}

	return sum, nil
}

func SumCryptoHoldings(h *Holdings) (float64, error) {
	var sum float64

	for symbol, amount := range h.Crypto {
		price, err := crypto.GetCryptoPrice(symbol)
		if err != nil {
			return 0, fmt.Errorf("could not sum crypto networth %s", err)
		}

		totalPricePerCrypto := amount * price
		sum += totalPricePerCrypto
	}

	return sum, nil
}

func SumBullionHoldings(h *Holdings) (float64, error) {
	var sum float64

	for metal, amount := range h.Bullion {
		priceResponse, err := bullion.GetBullionPrices(metal, "USD")
		if err != nil {
			return 0, fmt.Errorf("could not sum bullion networth %s", err)
		}

		sum += amount * priceResponse.Price
	}

	return sum, nil
}
