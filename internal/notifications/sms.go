package notifications

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TwilioNotifier struct {
	AccountSID string
	AuthToken  string
	FromNumber string
	ToNumber   string
}

func NewTwilioNotifier() *TwilioNotifier {
	return &TwilioNotifier{
		AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		FromNumber: os.Getenv("TWILIO_FROM_NUMBER"),
		ToNumber:   os.Getenv("TWILIO_TO_NUMBER"),
	}
}

func (t *TwilioNotifier) SendSMS(message string) error {
	if t.AccountSID == "" || t.AuthToken == "" || t.FromNumber == "" || t.ToNumber == "" {
		return fmt.Errorf("Twilio credentials are required")
	}

	data := url.Values{}
	data.Set("To", t.ToNumber)
	data.Set("From", t.FromNumber)
	data.Set("Body", message)

	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSID)
	
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.AccountSID, t.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Twilio API returned status: %d", resp.StatusCode)
	}

	return nil
}

func (t *TwilioNotifier) SendInvestmentUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN float64) error {
	message := fmt.Sprintf(
		"💰 Daily Investment Update\n\n"+
			"Trading212: %.2f BGN\n"+
			"Crypto: %.2f BGN\n"+
			"Bullion: %.2f BGN\n\n"+
			"Total: %.2f BGN\n\n"+
			"Have a great day! 📈",
		stocksBGN, cryptoBGN, bullionBGN, totalBGN,
	)

	return t.SendSMS(message)
}