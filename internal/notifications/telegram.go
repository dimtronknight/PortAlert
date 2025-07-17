package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TelegramNotifier struct {
	BotToken string
	ChatID   string
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewTelegramNotifier() *TelegramNotifier {
	return &TelegramNotifier{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}
}

func (t *TelegramNotifier) SendMessage(message string) error {
	if t.BotToken == "" || t.ChatID == "" {
		return fmt.Errorf("Telegram bot token and chat ID are required")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	telegramMsg := TelegramMessage{
		ChatID:    t.ChatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API returned status: %d", resp.StatusCode)
	}

	return nil
}

func (t *TelegramNotifier) SendInvestmentUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN float64) error {
	message := fmt.Sprintf(
		"*Daily Investment Update*\n\n"+
			"- Trading212: 	`%.2f BGN`\n"+
			"- Crypto: 		`%.2f BGN`\n"+
			"- Bullion: 	`%.2f BGN`\n\n"+
			" *Total: %.2f BGN*\n\n",
		stocksBGN, cryptoBGN, bullionBGN, totalBGN,
	)

	return t.SendMessage(message)
}
