package main

import (
	"database/sql"
	"fmt"
	"time"
)

// BidStorage работает со ставками
type BidStorage struct {
	db *sql.DB
}

// NewBidStorage создаёт хранилище
func NewBidStorage(db *sql.DB) *BidStorage {
	return &BidStorage{db: db}
}

// Create создаёт новую ставку
func (s *BidStorage) Create(bid Bid) (Bid, error) {
	query := `
        INSERT INTO bids (order_id, contractor_id, proposed_time)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err := s.db.QueryRow(
		query,
		bid.OrderID,
		bid.ContractorID,
		bid.ProposedTime,
	).Scan(&bid.ID)

	if err != nil {
		return Bid{}, err
	}

	return bid, nil
}

// GetByOrderID получает все ставки для заказа
func (s *BidStorage) GetByOrderID(orderID int) ([]Bid, error) {
	query := `
        SELECT id, order_id, contractor_id, proposed_time
        FROM bids
        WHERE order_id = $1
        ORDER BY created_at ASC
    `

	rows, err := s.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []Bid
	for rows.Next() {
		var bid Bid
		err := rows.Scan(
			&bid.ID,
			&bid.OrderID,
			&bid.ContractorID,
			&bid.ProposedTime,
		)
		if err != nil {
			continue
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

// GetByContractorID получает все ставки подрядчика
func (s *BidStorage) GetByContractorID(contractorID int) ([]Bid, error) {
	query := `
        SELECT id, order_id, contractor_id, proposed_time
        FROM bids
        WHERE contractor_id = $1
        ORDER BY created_at DESC
    `

	rows, err := s.db.Query(query, contractorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []Bid
	for rows.Next() {
		var bid Bid
		err := rows.Scan(
			&bid.ID,
			&bid.OrderID,
			&bid.ContractorID,
			&bid.ProposedTime,
		)
		if err != nil {
			continue
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

// HasBid проверяет сделал ли подрядчик уже ставку на этот заказ
func (s *BidStorage) HasBid(orderID int, contractorID int) (bool, error) {
	query := `
        SELECT COUNT(*) FROM bids
        WHERE order_id = $1 AND contractor_id = $2
    `

	var count int
	err := s.db.QueryRow(query, orderID, contractorID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SelectContractor выбирает подрядчика для заказа
func (s *BidStorage) SelectContractor(orderID int) (*Bid, error) {
	// Получаем все ставки для заказа
	bids, err := s.GetByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	if len(bids) == 0 {
		return nil, fmt.Errorf("no bids for order %d", orderID)
	}

	// Ищем ставку "Today"
	var todayBids []Bid
	for _, bid := range bids {
		if bid.ProposedTime == "Today" {
			todayBids = append(todayBids, bid)
		}
	}

	// Если есть "Today" — выбираем первого
	if len(todayBids) > 0 {
		return &todayBids[0], nil
	}

	// Иначе выбираем "Tomorrow"
	var tomorrowBids []Bid
	for _, bid := range bids {
		if bid.ProposedTime == "Tomorrow" {
			tomorrowBids = append(tomorrowBids, bid)
		}
	}

	if len(tomorrowBids) > 0 {
		return &tomorrowBids[0], nil
	}

	// Иначе выбираем любую первую
	return &bids[0], nil
}

// ScheduleSelection запускает таймер для выбора подрядчика
func (s *BidStorage) ScheduleSelection(orderID int, storage *OrderStorage, delay time.Duration) {
	go func() {
		fmt.Printf("⏰ Timer started for order #%d (waiting %v)\n", orderID, delay)

		// Ждём указанное время
		time.Sleep(delay)

		fmt.Printf("⏰ Timer fired for order #%d, selecting contractor...\n", orderID)

		// Проверяем что заказ всё ещё в статусе "confirmed"
		order, found := storage.GetByID(orderID)
		if !found {
			fmt.Printf("❌ Order #%d not found\n", orderID)
			return
		}

		if order.Status != "confirmed" {
			fmt.Printf("ℹ️ Order #%d already processed (status: %s)\n", orderID, order.Status)
			return
		}

		// Выбираем подрядчика
		winningBid, err := s.SelectContractor(orderID)
		if err != nil {
			fmt.Printf("❌ Error selecting contractor for order #%d: %v\n", orderID, err)
			return
		}

		// Назначаем
		err = storage.AssignContractor(orderID, winningBid.ContractorID)
		if err != nil {
			fmt.Printf("❌ Error assigning contractor for order #%d: %v\n", orderID, err)
			return
		}

		fmt.Printf("✅ Order #%d assigned to contractor #%d (time: %s)\n", orderID, winningBid.ContractorID, winningBid.ProposedTime)
	}()
}
