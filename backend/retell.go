package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var RETELL_API_KEY = os.Getenv("RETELL_API_KEY")

// Агент для сценария 1 — перезвонить клиенту после формы (данные уже есть)
const RETELL_AGENT_ID_OUTBOUND = "agent_fbdf096c0cabbeebbdfc934e69"

// Агент для сценария 2 — входящий звонок, собирает данные с нуля
const RETELL_AGENT_ID_INBOUND = "agent_aa549d074d5da9f9e8317bc990"

// InitiateRetellCall — исходящий AI звонок клиенту (сценарий 1)
// Используется когда клиент заполнил форму — данные уже известны
func InitiateRetellCall(phone string, orderID int, clientName string, device string) error {
	return initiateCall(phone, orderID, RETELL_AGENT_ID_OUTBOUND, map[string]string{
		"customer_name":  clientName,
		"appliance_type": device,
		"order_id":       fmt.Sprintf("%d", orderID),
	})
}

// InitiateRetellQuoteCall — исходящий AI звонок для сценария 2 (кнопка Get a Quote)
// AI собирает все данные с нуля и вызывает confirm_order function
func InitiateRetellQuoteCall(phone string, orderID int) error {
	return initiateCall(phone, orderID, RETELL_AGENT_ID_INBOUND, map[string]string{
		"order_id": fmt.Sprintf("%d", orderID),
	})
}

// initiateCall — общая функция для запуска звонка через Retell
func initiateCall(phone string, orderID int, agentID string, dynamicVars map[string]string) error {
	apiURL := "https://api.retellai.com/v2/create-phone-call"

	payload := map[string]interface{}{
		"from_number":                  TWILIO_PHONE,
		"to_number":                    phone,
		"agent_id":                     agentID,
		"retell_llm_dynamic_variables": dynamicVars,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+RETELL_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("retell returned status %d", resp.StatusCode)
	}

	fmt.Printf("✅ Retell call initiated to %s for order #%d (agent: %s)\n", phone, orderID, agentID)
	return nil
}
