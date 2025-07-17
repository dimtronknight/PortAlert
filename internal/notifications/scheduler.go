package notifications

import (
	"fmt"
	"log"
	"time"
)

type Scheduler struct {
	service          *NotificationService
	done             chan bool
	onNotificationSent func(total, stocks, crypto, bullion float64)
}

func NewScheduler(service *NotificationService) *Scheduler {
	return &Scheduler{
		service: service,
		done:    make(chan bool),
	}
}

func (s *Scheduler) SetOnNotificationSent(callback func(total, stocks, crypto, bullion float64)) {
	s.onNotificationSent = callback
}

func (s *Scheduler) Start(hour, minute int, getPortfolioValues func() (total, stocks, crypto, bullion float64, err error)) {
	log.Printf("Starting daily notification scheduler for %02d:%02d", hour, minute)

	// Calculate time until next notification
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	
	// If the time has already passed today, schedule for tomorrow
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	log.Printf("Next notification scheduled for: %v", nextRun)

	// Wait until the first scheduled time
	timer := time.NewTimer(nextRun.Sub(now))
	
	go func() {
		for {
			select {
			case <-timer.C:
				// Send notification
				s.sendScheduledNotification(getPortfolioValues)
				
				// Schedule next notification (24 hours from now)
				timer.Reset(24 * time.Hour)
				
			case <-s.done:
				timer.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) sendScheduledNotification(getPortfolioValues func() (total, stocks, crypto, bullion float64, err error)) {
	log.Println("Sending scheduled investment notification...")
	
	total, stocks, crypto, bullion, err := getPortfolioValues()
	if err != nil {
		log.Printf("Error getting portfolio values: %v", err)
		return
	}

	err = s.service.SendDailyUpdate(total, stocks, crypto, bullion)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	} else {
		log.Println("Successfully sent scheduled notification")
		
		// Call the callback to save to database
		if s.onNotificationSent != nil {
			s.onNotificationSent(total, stocks, crypto, bullion)
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.done)
}

// Helper function to run once immediately (for testing)
func (s *Scheduler) SendNow(getPortfolioValues func() (total, stocks, crypto, bullion float64, err error)) error {
	log.Println("Sending immediate investment notification...")
	
	total, stocks, crypto, bullion, err := getPortfolioValues()
	if err != nil {
		return fmt.Errorf("error getting portfolio values: %w", err)
	}

	return s.service.SendDailyUpdate(total, stocks, crypto, bullion)
}