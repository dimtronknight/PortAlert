package stocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) GetAccountInfo() (*AccountInfo, error) {
	data, err := c.makeRequest(http.MethodGet, "/equity/account/info", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	var accountInfo AccountInfo
	if err := json.Unmarshal(data, &accountInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account info: %w", err)
	}

	return &accountInfo, nil
}

func (c *Client) GetAccountCash() (*AccountCash, error) {
	data, err := c.makeRequest(http.MethodGet, "/equity/account/cash", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account cash: %w", err)
	}

	var accountCash AccountCash
	if err := json.Unmarshal(data, &accountCash); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account cash: %w", err)
	}

	return &accountCash, nil
}

func (c *Client) GetPortfolio() ([]Position, error) {
	data, err := c.makeRequest(http.MethodGet, "/equity/portfolio", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	var positions []Position
	if err := json.Unmarshal(data, &positions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal portfolio: %w", err)
	}

	return positions, nil
}

func (c *Client) GetPosition(ticker string) (*Position, error) {
	endpoint := fmt.Sprintf("/equity/portfolio/%s", ticker)
	data, err := c.makeRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get position for %s: %w", ticker, err)
	}

	var position Position
	if err := json.Unmarshal(data, &position); err != nil {
		return nil, fmt.Errorf("failed to unmarshal position: %w", err)
	}

	return &position, nil
}

func (c *Client) GetOrders() ([]Order, error) {
	data, err := c.makeRequest(http.MethodGet, "/equity/orders", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var orders []Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders: %w", err)
	}

	return orders, nil
}

func (c *Client) GetOrder(orderID int64) (*Order, error) {
	endpoint := fmt.Sprintf("/equity/orders/%d", orderID)
	data, err := c.makeRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}

func (c *Client) PlaceMarketOrder(ticker string, quantity float64) (*Order, error) {
	orderRequest := MarketOrderRequest{
		Ticker:   ticker,
		Quantity: quantity,
	}

	data, err := c.makeRequest(http.MethodPost, "/equity/orders/market", orderRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to place market order: %w", err)
	}

	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &order, nil
}

func (c *Client) PlaceLimitOrder(ticker string, quantity, limitPrice float64, timeInForce string) (*Order, error) {
	orderRequest := LimitOrderRequest{
		Ticker:      ticker,
		Quantity:    quantity,
		LimitPrice:  limitPrice,
		TimeInForce: timeInForce,
	}

	data, err := c.makeRequest(http.MethodPost, "/equity/orders/limit", orderRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to place limit order: %w", err)
	}

	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &order, nil
}

func (c *Client) CancelOrder(orderID int64) error {
	endpoint := fmt.Sprintf("/equity/orders/%d", orderID)
	_, err := c.makeRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel order %d: %w", orderID, err)
	}

	return nil
}

func (c *Client) GetExchanges() ([]Exchange, error) {
	data, err := c.makeRequest(http.MethodGet, "/equity/metadata/exchanges", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchanges: %w", err)
	}

	var exchanges []Exchange
	if err := json.Unmarshal(data, &exchanges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exchanges: %w", err)
	}

	return exchanges, nil
}

func (c *Client) GetInstruments() ([]Instrument, error) {
	data, err := c.makeRequest(http.MethodGet, "/equity/metadata/instruments", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get instruments: %w", err)
	}

	var instruments []Instrument
	if err := json.Unmarshal(data, &instruments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instruments: %w", err)
	}

	return instruments, nil
}

func (c *Client) GetHistoricalOrders(cursor string, limit int) ([]Order, error) {
	endpoint := "/equity/history/orders"
	if cursor != "" || limit > 0 {
		endpoint += "?"
		if cursor != "" {
			endpoint += "cursor=" + cursor
		}
		if limit > 0 {
			if cursor != "" {
				endpoint += "&"
			}
			endpoint += "limit=" + strconv.Itoa(limit)
		}
	}

	data, err := c.makeRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical orders: %w", err)
	}

	var orders []Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal historical orders: %w", err)
	}

	return orders, nil
}

func (c *Client) GetDividends(cursor string, limit int) ([]Dividend, error) {
	endpoint := "/history/dividends"
	if cursor != "" || limit > 0 {
		endpoint += "?"
		if cursor != "" {
			endpoint += "cursor=" + cursor
		}
		if limit > 0 {
			if cursor != "" {
				endpoint += "&"
			}
			endpoint += "limit=" + strconv.Itoa(limit)
		}
	}

	data, err := c.makeRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get dividends: %w", err)
	}

	var dividends []Dividend
	if err := json.Unmarshal(data, &dividends); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dividends: %w", err)
	}

	return dividends, nil
}

func (c *Client) GetTransactions(cursor string, limit int) ([]Transaction, error) {
	endpoint := "/history/transactions"
	if cursor != "" || limit > 0 {
		endpoint += "?"
		if cursor != "" {
			endpoint += "cursor=" + cursor
		}
		if limit > 0 {
			if cursor != "" {
				endpoint += "&"
			}
			endpoint += "limit=" + strconv.Itoa(limit)
		}
	}

	data, err := c.makeRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var transactions []Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions: %w", err)
	}

	return transactions, nil
}