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
	"time"
)

var (
	TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN  = os.Getenv("TWILIO_AUTH_TOKEN")
	TWILIO_PHONE       = os.Getenv("TWILIO_PHONE")
)

// SendSMS отправляет SMS через Twilio REST API
func SendSMS(to string, message string) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", TWILIO_ACCOUNT_SID)

	data := url.Values{}
	data.Set("To", to)
	data.Set("From", TWILIO_PHONE)
	data.Set("Body", message)

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

	fmt.Printf("✅ SMS sent to %s\n", to)
	return nil
}

// GenerateToken создаёт случайный токен
func GenerateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// BroadcastJobToContractors рассылает SMS всем подрядчикам с задержкой 200мс
func BroadcastJobToContractors(order Order, contractors []Contractor, baseURL string) string {
	token := GenerateToken()
	jobURL := fmt.Sprintf("%s/accept/%s", baseURL, token)
	message := fmt.Sprintf(
		"New Job - %s repair - %s\nAccept: %s",
		order.Device,
		order.ZipCode,
		jobURL,
	)

	successCount := 0
	for _, contractor := range contractors {
		err := SendSMS(contractor.Phone, message)
		if err != nil {
			fmt.Printf("❌ Failed to send to %s: %v\n", contractor.Phone, err)
		} else {
			successCount++
		}
		// Задержка между SMS чтобы не попасть под спам-фильтры операторов
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("📤 Broadcast sent to %d/%d contractors for order #%d\n", successCount, len(contractors), order.ID)

	return token
}

// InitiateCall звонит подрядчику, потом соединяет с клиентом
// Twilio после завершения звонка пришлёт webhook на /api/call-status
func InitiateCall(contractorPhone string, clientPhone string, orderID int, contractorID int) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", TWILIO_ACCOUNT_SID)

	twimlURL := fmt.Sprintf(
		"https://fixly-production.up.railway.app/api/twiml?client_phone=%s&order_id=%d",
		url.QueryEscape(clientPhone), orderID,
	)

	// StatusCallback — Twilio пришлёт сюда длительность звонка после завершения
	statusCallbackURL := fmt.Sprintf(
		"https://fixly-production.up.railway.app/api/call-status?order_id=%d&contractor_id=%d",
		orderID, contractorID,
	)

	data := url.Values{}
	data.Set("To", contractorPhone)
	data.Set("From", TWILIO_PHONE)
	data.Set("Url", twimlURL)
	data.Set("StatusCallback", statusCallbackURL)
	data.Set("StatusCallbackMethod", "POST")
	// Twilio шлёт callback на эти события; "completed" содержит итоговую длительность
	data.Set("StatusCallbackEvent", "completed")

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
