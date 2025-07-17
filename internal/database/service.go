package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PortfolioService struct {
	db         *MongoDB
	collection *mongo.Collection
}

func NewPortfolioService(db *MongoDB) *PortfolioService {
	return &PortfolioService{
		db:         db,
		collection: db.GetCollection("portfolio_snapshots"),
	}
}

// SaveDailySnapshot saves a daily portfolio snapshot
func (ps *PortfolioService) SaveDailySnapshot(trading212, crypto, bullion, total float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use today's date (without time) as the key
	today := time.Now().Truncate(24 * time.Hour)

	snapshot := PortfolioSnapshot{
		Date:          today,
		Trading212BGN: trading212,
		CryptoBGN:     crypto,
		BullionBGN:    bullion,
		TotalBGN:      total,
		CreatedAt:     time.Now(),
		USDToBGNRate:  1.7346, // You might want to make this dynamic
	}

	// Use upsert to replace if already exists for today
	filter := bson.M{"date": today}
	update := bson.M{"$set": snapshot}
	opts := options.Update().SetUpsert(true)

	result, err := ps.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save portfolio snapshot: %w", err)
	}

	if result.UpsertedCount > 0 {
		log.Printf("Created new portfolio snapshot for %s", today.Format("2006-01-02"))
	} else {
		log.Printf("Updated existing portfolio snapshot for %s", today.Format("2006-01-02"))
	}

	return nil
}

// GetLatestSnapshot returns the most recent portfolio snapshot
func (ps *PortfolioService) GetLatestSnapshot() (*PortfolioSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var snapshot PortfolioSnapshot
	opts := options.FindOne().SetSort(bson.D{{"date", -1}})

	err := ps.collection.FindOne(ctx, bson.M{}, opts).Decode(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no portfolio snapshots found")
		}
		return nil, fmt.Errorf("failed to get latest snapshot: %w", err)
	}

	return &snapshot, nil
}

// GetSnapshotsByDateRange returns snapshots between two dates
func (ps *PortfolioService) GetSnapshotsByDateRange(start, end time.Time) ([]PortfolioSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"date": bson.M{
			"$gte": start.Truncate(24 * time.Hour),
			"$lte": end.Truncate(24 * time.Hour),
		},
	}

	opts := options.Find().SetSort(bson.D{{"date", 1}})
	cursor, err := ps.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshots: %w", err)
	}
	defer cursor.Close(ctx)

	var snapshots []PortfolioSnapshot
	err = cursor.All(ctx, &snapshots)
	if err != nil {
		return nil, fmt.Errorf("failed to decode snapshots: %w", err)
	}

	return snapshots, nil
}

// GetLastNDays returns the last N days of portfolio data
func (ps *PortfolioService) GetLastNDays(days int) ([]PortfolioSnapshot, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -days)
	
	return ps.GetSnapshotsByDateRange(start, end)
}

// GetPortfolioStats calculates portfolio statistics
func (ps *PortfolioService) GetPortfolioStats() (*PortfolioStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all snapshots
	cursor, err := ps.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots for stats: %w", err)
	}
	defer cursor.Close(ctx)

	var snapshots []PortfolioSnapshot
	err = cursor.All(ctx, &snapshots)
	if err != nil {
		return nil, fmt.Errorf("failed to decode snapshots: %w", err)
	}

	if len(snapshots) == 0 {
		return nil, fmt.Errorf("no data available for statistics")
	}

	// Calculate stats
	stats := &PortfolioStats{
		DaysTracked: len(snapshots),
	}

	var totalValue float64
	stats.BestDayValue = snapshots[0].TotalBGN
	stats.WorstDayValue = snapshots[0].TotalBGN
	stats.BestDay = snapshots[0].Date
	stats.WorstDay = snapshots[0].Date

	for _, snapshot := range snapshots {
		totalValue += snapshot.TotalBGN

		if snapshot.TotalBGN > stats.BestDayValue {
			stats.BestDayValue = snapshot.TotalBGN
			stats.BestDay = snapshot.Date
		}

		if snapshot.TotalBGN < stats.WorstDayValue {
			stats.WorstDayValue = snapshot.TotalBGN
			stats.WorstDay = snapshot.Date
		}
	}

	stats.AverageValue = totalValue / float64(len(snapshots))

	// Calculate growth (first vs last)
	if len(snapshots) > 1 {
		firstValue := snapshots[0].TotalBGN
		lastValue := snapshots[len(snapshots)-1].TotalBGN
		stats.TotalGrowth = lastValue - firstValue
		stats.GrowthPercentage = (stats.TotalGrowth / firstValue) * 100
	}

	return stats, nil
}

// GetTodaySnapshot returns today's snapshot if it exists
func (ps *PortfolioService) GetTodaySnapshot() (*PortfolioSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	today := time.Now().Truncate(24 * time.Hour)
	filter := bson.M{"date": today}

	var snapshot PortfolioSnapshot
	err := ps.collection.FindOne(ctx, filter).Decode(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no snapshot found for today")
		}
		return nil, fmt.Errorf("failed to get today's snapshot: %w", err)
	}

	return &snapshot, nil
}

// DeleteOldSnapshots removes snapshots older than specified days
func (ps *PortfolioService) DeleteOldSnapshots(daysToKeep int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep).Truncate(24 * time.Hour)
	filter := bson.M{"date": bson.M{"$lt": cutoffDate}}

	result, err := ps.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete old snapshots: %w", err)
	}

	if result.DeletedCount > 0 {
		log.Printf("Deleted %d old portfolio snapshots", result.DeletedCount)
	}

	return nil
}