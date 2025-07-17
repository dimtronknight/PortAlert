package bullion

type BullionPriceResponse struct {
	Timestamp      int64   `json:"timestamp"`
	Metal          string  `json:"metal"`
	Currency       string  `json:"currency"`
	Exchange       string  `json:"exchange"`
	Symbol         string  `json:"symbol"`
	PrevClosePrice float64 `json:"prev_close_price"`
	OpenPrice      float64 `json:"open_price"`
	LowPrice       float64 `json:"low_price"`
	HighPrice      float64 `json:"high_price"`
	OpenTime       int64   `json:"open_time"`
	Price          float64 `json:"price"`
	Change         float64 `json:"ch"`
	ChangePercent  float64 `json:"chp"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	PriceGram24K   float64 `json:"price_gram_24k"`
	PriceGram22K   float64 `json:"price_gram_22k"`
	PriceGram21K   float64 `json:"price_gram_21k"`
	PriceGram20K   float64 `json:"price_gram_20k"`
	PriceGram18K   float64 `json:"price_gram_18k"`
	PriceGram16K   float64 `json:"price_gram_16k"`
	PriceGram14K   float64 `json:"price_gram_14k"`
	PriceGram10K   float64 `json:"price_gram_10k"`
}
