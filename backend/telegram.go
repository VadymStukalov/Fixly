package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	TELEGRAM_BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	TELEGRAM_CHAT_ID   = os.Getenv("TELEGRAM_CHAT_ID")
)

// SendTelegramMessage отправляет сообщение админу в Telegram
func SendTelegramMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_BOT_TOKEN)

	body, _ := json.Marshal(map[string]string{
		"chat_id":    TELEGRAM_CHAT_ID,
		"text":       message,
		"parse_mode": "HTML",
	})

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram returned status %d", resp.StatusCode)
	}

	fmt.Println("📨 Telegram notification sent to admin")
	return nil
}

// NotifyAdminClientUnreachable — уведомление когда клиент не берёт трубку 3 раза
func NotifyAdminClientUnreachable(order Order, contractorID int) {
	message := fmt.Sprintf(
		"⚠️ <b>Client Unreachable</b>\n\n"+
			"📋 Order #%d\n"+
			"👤 Client: %s\n"+
			"📱 Phone: %s\n"+
			"🔧 Device: %s\n"+
			"📍 ZIP: %s\n\n"+
			"🔨 Contractor #%d called 3 times — no answer.\n"+
			"Please contact the client manually.",
		order.ID,
		order.ClientName,
		order.Phone,
		order.Device,
		order.ZipCode,
		contractorID,
	)

	err := SendTelegramMessage(message)
	if err != nil {
		fmt.Printf("❌ Failed to send Telegram notification: %v\n", err)
	}
}
