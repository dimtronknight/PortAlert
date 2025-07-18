package main

import (
	"fmt"
	"investment-tracker/internal/notifications"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Check if required environment variables are set
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}
	if chatID == "" {
		log.Fatal("TELEGRAM_CHAT_ID environment variable is required")
	}

	fmt.Printf("Bot Token: %s...\n", botToken[:10])
	fmt.Printf("Chat ID: %s\n", chatID)

	// Create Telegram notifier
	notifier := notifications.NewTelegramNotifier()

	// Test connection
	fmt.Println("\n1. Testing bot connection...")
	err = notifier.TestConnection()
	if err != nil {
		log.Fatalf("Failed to test connection: %v", err)
	}
	fmt.Println("✅ Connection test successful!")

	// Get chat info
	fmt.Println("\n2. Getting chat information...")
	err = notifier.GetChatInfo()
	if err != nil {
		log.Printf("Failed to get chat info: %v", err)
	} else {
		fmt.Println("✅ Chat info retrieved successfully!")
	}

	// Send test investment update
	fmt.Println("\n3. Sending test investment update...")
	err = notifier.SendInvestmentUpdate(5000.00, 2000.00, 1500.00, 1500.00)
	if err != nil {
		log.Fatalf("Failed to send investment update: %v", err)
	}
	fmt.Println("✅ Investment update sent successfully!")

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("\nIf you're still not receiving notifications, check:")
	fmt.Println("1. Your phone's notification settings for Telegram")
	fmt.Println("2. Telegram's notification settings for the bot")
	fmt.Println("3. Whether the bot is muted in the chat")
	fmt.Println("4. Your device's Do Not Disturb settings")
}
