package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN  = os.Getenv("TWILIO_AUTH_TOKEN")
	TWILIO_PHONE       = os.Getenv("TWILIO_PHONE")
)

// SendSMS отправляет SMS через Twilio REST API
func SendSMS(to string, message string) error {
	// URL Twilio API
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", TWILIO_ACCOUNT_SID)

	// Параметры
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", TWILIO_PHONE)
	data.Set("Body", message)

	// Создаём HTTP запрос
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	// Basic Auth (Account SID + Auth Token)
	auth := base64.StdEncoding.EncodeToString([]byte(TWILIO_ACCOUNT_SID + ":" + TWILIO_AUTH_TOKEN))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправляем
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус
	if resp.StatusCode != 201 {
		return fmt.Errorf("twilio returned status %d", resp.StatusCode)
	}

	fmt.Printf("✅ SMS sent to %s\n", to)
	return nil
}

// GenerateToken создаёт случайный токен
func GenerateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// BroadcastJobToContractors рассылает SMS всем подрядчикам
func BroadcastJobToContractors(order Order, contractors []Contractor, baseURL string) string {
	// Генерируем уникальный токен
	token := GenerateToken()

	// Формируем ссылку
	jobURL := fmt.Sprintf("%s/accept/%s", baseURL, token)

	// Формируем текст SMS
	message := fmt.Sprintf(
		"New Job - %s repair - %s\nAccept: %s",
		order.Device,
		order.ZipCode,
		jobURL,
	)

	// Отправляем всем подрядчикам
	successCount := 0
	for _, contractor := range contractors {
		err := SendSMS(contractor.Phone, message)
		if err != nil {
			fmt.Printf("❌ Failed to send to %s: %v\n", contractor.Phone, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("📤 Broadcast sent to %d/%d contractors for order #%d\n", successCount, len(contractors), order.ID)

	return token
}

// InitiateCall звонит подрядчику, потом соединяет с клиентом
func InitiateCall(contractorPhone string, clientPhone string, orderID int) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", TWILIO_ACCOUNT_SID)

	// URL нашего TwiML endpoint — Twilio будет его вызывать
	twimlURL := fmt.Sprintf("https://rare-zebras-write.loca.lt/api/twiml?client_phone=%s&order_id=%d",
		url.QueryEscape(clientPhone), orderID)

	data := url.Values{}
	data.Set("To", contractorPhone)
	data.Set("From", TWILIO_PHONE)
	data.Set("Url", twimlURL)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(TWILIO_ACCOUNT_SID + ":" + TWILIO_AUTH_TOKEN))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("twilio returned status %d", resp.StatusCode)
	}

	fmt.Printf("✅ Call initiated to contractor %s for order #%d\n", contractorPhone, orderID)
	return nil
}
