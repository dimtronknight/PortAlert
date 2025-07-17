package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PortfolioSnapshot represents a daily snapshot of portfolio values
type PortfolioSnapshot struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Date             time.Time          `bson:"date"`
	Trading212BGN    float64            `bson:"trading212_bgn"`
	CryptoBGN        float64            `bson:"crypto_bgn"`
	BullionBGN       float64            `bson:"bullion_bgn"`
	TotalBGN         float64            `bson:"total_bgn"`
	CreatedAt        time.Time          `bson:"created_at"`
	
	// Optional: Store individual asset breakdowns
	CryptoAssets     []AssetValue       `bson:"crypto_assets,omitempty"`
	BullionAssets    []AssetValue       `bson:"bullion_assets,omitempty"`
	
	// Optional: Store exchange rates used
	USDToBGNRate     float64            `bson:"usd_to_bgn_rate,omitempty"`
}

// AssetValue represents individual asset values
type AssetValue struct {
	Symbol    string  `bson:"symbol"`
	Amount    float64 `bson:"amount"`
	Price     float64 `bson:"price"`
	Value     float64 `bson:"value"`
	Currency  string  `bson:"currency"`
}

// PortfolioStats represents aggregated statistics
type PortfolioStats struct {
	TotalGrowth      float64   `bson:"total_growth"`
	GrowthPercentage float64   `bson:"growth_percentage"`
	BestDay          time.Time `bson:"best_day"`
	BestDayValue     float64   `bson:"best_day_value"`
	WorstDay         time.Time `bson:"worst_day"`
	WorstDayValue    float64   `bson:"worst_day_value"`
	AverageValue     float64   `bson:"average_value"`
	DaysTracked      int       `bson:"days_tracked"`
}