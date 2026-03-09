package main

import "database/sql"

type CallLog struct {
	ID            int    `json:"id"`
	OrderID       int    `json:"order_id"`
	ContractorID  int    `json:"contractor_id"`
	TwilioCallSID string `json:"twilio_call_sid"`
	Duration      int    `json:"duration"`
	Status        string `json:"status"`
}

type CallLogStorage struct {
	db *sql.DB
}

func NewCallLogStorage(db *sql.DB) *CallLogStorage {
	return &CallLogStorage{db: db}
}

// SaveCallLog сохраняет лог звонка
func (s *CallLogStorage) SaveCallLog(orderID int, contractorID int, sid string, duration int, status string) error {
	query := `
		INSERT INTO call_logs (order_id, contractor_id, twilio_call_sid, duration, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.Exec(query, orderID, contractorID, sid, duration, status)
	return err
}

// HasSuccessfulCall проверяет был ли звонок 30+ секунд по заказу
func (s *CallLogStorage) HasSuccessfulCall(orderID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM call_logs
		WHERE order_id = $1 AND duration >= 30
	`
	var count int
	err := s.db.QueryRow(query, orderID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasAnyCallAttempt проверяет были ли вообще звонки по заказу
func (s *CallLogStorage) HasAnyCallAttempt(orderID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM call_logs
		WHERE order_id = $1
	`
	var count int
	err := s.db.QueryRow(query, orderID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
