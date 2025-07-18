package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type TelegramNotifier struct {
	BotToken string
	ChatID   string
}

type TelegramMessage struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableNotification   bool   `json:"disable_notification"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
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
		ChatID:                t.ChatID,
		Text:                  message,
		ParseMode:             "Markdown",
		DisableNotification:   false,
		DisableWebPagePreview: true,
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
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API returned status: %d, body: %s", resp.StatusCode, string(body))
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

// TestConnection tests if the bot can send messages
func (t *TelegramNotifier) TestConnection() error {
	if t.BotToken == "" || t.ChatID == "" {
		return fmt.Errorf("Telegram bot token and chat ID are required")
	}

	// Test bot info
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", t.BotToken)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get bot info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get bot info, status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Send test message
	return t.SendMessage("🔔 *Test Notification*\n\nThis is a test message to check if notifications are working properly.")
}

// GetChatInfo gets information about the chat
func (t *TelegramNotifier) GetChatInfo() error {
	if t.BotToken == "" || t.ChatID == "" {
		return fmt.Errorf("Telegram bot token and chat ID are required")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat", t.BotToken)
	
	data := map[string]string{
		"chat_id": t.ChatID,
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal chat request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to get chat info: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get chat info, status: %d, body: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Chat info: %s\n", string(body))
	return nil
}
