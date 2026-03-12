package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var RETELL_API_KEY = os.Getenv("RETELL_API_KEY")

const RETELL_AGENT_ID = "agent_fbdf096c0cabbeebbdfc934e69"

// InitiateRetellCall — запускает исходящий AI звонок клиенту через Retell
func InitiateRetellCall(phone string, orderID int, clientName string, device string) error {
	apiURL := "https://api.retellai.com/v2/create-phone-call"

	payload := map[string]interface{}{
		"from_number": TWILIO_PHONE,
		"to_number":   phone,
		"agent_id":    RETELL_AGENT_ID,
		"retell_llm_dynamic_variables": map[string]string{
			"customer_name":  clientName,
			"appliance_type": device,
			"order_id":       fmt.Sprintf("%d", orderID),
		},
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

	fmt.Printf("✅ Retell call initiated to %s for order #%d\n", phone, orderID)
	return nil
}
