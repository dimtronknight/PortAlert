package stocks

import "time"

type AccountInfo struct {
	CurrencyCode string `json:"currencyCode"`
	ID           int64  `json:"id"`
}

type AccountCash struct {
	Free           float64 `json:"free"`
	Total          float64 `json:"total"`
	PieCash        float64 `json:"pieCash"`
	Result         float64 `json:"result"`
	Invested       float64 `json:"invested"`
	CurrencyCode   string  `json:"currencyCode"`
	BlockedForPies float64 `json:"blockedForPies"`
}

type Position struct {
	Ticker               string  `json:"ticker"`
	Quantity             float64 `json:"quantity"`
	AveragePrice         float64 `json:"averagePrice"`
	AveragePriceConverted float64 `json:"averagePriceConverted"`
	CurrentPrice         float64 `json:"currentPrice"`
	PieQuantity          float64 `json:"pieQuantity"`
	FxPnl                float64 `json:"fxPnl"`
	InitialFillDate      string  `json:"initialFillDate"`
	Frontend             string  `json:"frontend"`
	MaxBuy               float64 `json:"maxBuy"`
	MaxSell              float64 `json:"maxSell"`
}

type Order struct {
	CreationTime   time.Time `json:"creationTime"`
	FilledQuantity float64   `json:"filledQuantity"`
	ID             int64     `json:"id"`
	LimitPrice     float64   `json:"limitPrice"`
	Quantity       float64   `json:"quantity"`
	Status         string    `json:"status"`
	Strategy       string    `json:"strategy"`
	Ticker         string    `json:"ticker"`
	Type           string    `json:"type"`
	Value          float64   `json:"value"`
}

type MarketOrderRequest struct {
	Ticker   string  `json:"ticker"`
	Quantity float64 `json:"quantity"`
}

type LimitOrderRequest struct {
	Ticker     string  `json:"ticker"`
	Quantity   float64 `json:"quantity"`
	LimitPrice float64 `json:"limitPrice"`
	TimeInForce string  `json:"timeInForce"`
}

type Exchange struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	WorkingDays  string `json:"workingDays"`
	TimeZone     string `json:"timeZone"`
	OpeningTime  string `json:"openingTime"`
	ClosingTime  string `json:"closingTime"`
	LunchStart   string `json:"lunchStart"`
	LunchEnd     string `json:"lunchEnd"`
	OpeningDays  string `json:"openingDays"`
	ClosingDays  string `json:"closingDays"`
}

type Instrument struct {
	Ticker              string  `json:"ticker"`
	Name                string  `json:"name"`
	Type                string  `json:"type"`
	CurrencyCode        string  `json:"currencyCode"`
	ISIN                string  `json:"isin"`
	MaxOpenQuantity     float64 `json:"maxOpenQuantity"`
	MinTradeQuantity    float64 `json:"minTradeQuantity"`
	QuantityPrecision   int     `json:"quantityPrecision"`
	AddedOn             string  `json:"addedOn"`
	WorkingScheduleID   int     `json:"workingScheduleId"`
}

type Dividend struct {
	Ticker         string    `json:"ticker"`
	Type           string    `json:"type"`
	CashAmount     float64   `json:"cashAmount"`
	Quantity       float64   `json:"quantity"`
	Price          float64   `json:"price"`
	GrossAmountPerShare float64 `json:"grossAmountPerShare"`
	PaidOn         time.Time `json:"paidOn"`
	Reference      string    `json:"reference"`
}

type Transaction struct {
	ActionID        string    `json:"actionId"`
	DateTime        time.Time `json:"dateTime"`
	Reference       string    `json:"reference"`
	Type            string    `json:"type"`
	CurrencyCode    string    `json:"currencyCode"`
	Amount          float64   `json:"amount"`
	Ticker          string    `json:"ticker,omitempty"`
	Price           float64   `json:"price,omitempty"`
	Quantity        float64   `json:"quantity,omitempty"`
	Notes           string    `json:"notes,omitempty"`
}