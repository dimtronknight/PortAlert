package notifications

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type NotificationService struct {
	telegram *TelegramNotifier
	twilio   *TwilioNotifier
	email    *EmailNotifier
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		telegram: NewTelegramNotifier(),
		twilio:   NewTwilioNotifier(),
		email:    NewEmailNotifier(),
	}
}

func (ns *NotificationService) SendDailyUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN float64) error {
	// Get notification methods from environment variable
	methods := os.Getenv("NOTIFICATION_METHODS")
	if methods == "" {
		methods = "telegram" // Default to telegram
	}

	methodList := strings.Split(strings.ToLower(methods), ",")
	var errors []string

	for _, method := range methodList {
		method = strings.TrimSpace(method)
		
		switch method {
		case "telegram":
			if err := ns.telegram.SendInvestmentUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN); err != nil {
				errors = append(errors, fmt.Sprintf("Telegram: %v", err))
			} else {
				log.Printf("Successfully sent Telegram notification")
			}
		case "sms":
			if err := ns.twilio.SendInvestmentUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN); err != nil {
				errors = append(errors, fmt.Sprintf("SMS: %v", err))
			} else {
				log.Printf("Successfully sent SMS notification")
			}
		case "email":
			if err := ns.email.SendInvestmentUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN); err != nil {
				errors = append(errors, fmt.Sprintf("Email: %v", err))
			} else {
				log.Printf("Successfully sent email notification")
			}
		default:
			errors = append(errors, fmt.Sprintf("Unknown notification method: %s", method))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (ns *NotificationService) TestNotifications() error {
	// Test with dummy data
	return ns.SendDailyUpdate(5000.0, 2000.0, 2000.0, 1000.0)
}