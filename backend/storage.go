package main

import "database/sql"

// OrderStorage работает с PostgreSQL

type OrderStorage struct {
	db *sql.DB
}

// NewOrderStorage создаёт новое хранилище
func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

// GetAll возвращает все заказы из базы
func (s *OrderStorage) GetAll() []Order {
	query := `SELECT id, client_name, phone, device, COALESCE(brand, ''), problem, 
		zip_code, COALESCE(preferred_time, ''), status, price, contractor_id FROM orders`

	rows, err := s.db.Query(query)
	if err != nil {
		return []Order{}
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.ClientName, &order.Phone, &order.Device,
			&order.Brand, &order.Problem, &order.ZipCode, &order.PreferredTime,
			&order.Status, &order.Price, &order.ContractorID,
		)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return orders
}

// Create создаёт новый заказ в базе
func (s *OrderStorage) Create(order Order) Order {
	query := `
        INSERT INTO orders (client_name, phone, device, brand, problem, zip_code, preferred_time, status, price)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `
	err := s.db.QueryRow(
		query,
		order.ClientName, order.Phone, order.Device, order.Brand,
		order.Problem, order.ZipCode, order.PreferredTime, order.Status, order.Price,
	).Scan(&order.ID)

	if err != nil {
		return Order{}
	}
	return order
}

// Update обновляет заказ в базе
func (s *OrderStorage) Update(id int, order Order) bool {
	query := `
        UPDATE orders 
        SET client_name = $1, phone = $2, device = $3, brand = $4, problem = $5, 
            zip_code = $6, preferred_time = $7, status = $8, price = $9
        WHERE id = $10
    `
	result, err := s.db.Exec(
		query,
		order.ClientName, order.Phone, order.Device, order.Brand,
		order.Problem, order.ZipCode, order.PreferredTime, order.Status, order.Price,
		id,
	)
	if err != nil {
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

// Delete удаляет заказ из базы
func (s *OrderStorage) Delete(id int) bool {
	query := "DELETE FROM orders WHERE id = $1"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

// GetByID находит заказ по ID
func (s *OrderStorage) GetByID(id int) (*Order, bool) {
	query := `SELECT id, client_name, phone, device, COALESCE(brand, ''), problem, 
		zip_code, COALESCE(preferred_time, ''), status, price, contractor_id FROM orders WHERE id = $1`

	var order Order
	err := s.db.QueryRow(query, id).Scan(
		&order.ID, &order.ClientName, &order.Phone, &order.Device,
		&order.Brand, &order.Problem, &order.ZipCode, &order.PreferredTime,
		&order.Status, &order.Price, &order.ContractorID,
	)
	if err != nil {
		return nil, false
	}
	return &order, true
}

// AssignContractor назначает подрядчика на заказ
func (s *OrderStorage) AssignContractor(orderID int, contractorID int) error {
	query := `UPDATE orders SET status = 'in_progress', contractor_id = $1 WHERE id = $2`
	_, err := s.db.Exec(query, contractorID, orderID)
	return err
}

// AcceptOrder — атомарный захват заказа (защита от race condition)
func (s *OrderStorage) AcceptOrder(orderID int, contractorID int) (bool, error) {
	query := `
		UPDATE orders
		SET status = 'in_progress', contractor_id = $1, accepted_at = NOW()
		WHERE id = $2 AND status = 'confirmed'
	`
	result, err := s.db.Exec(query, contractorID, orderID)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

// ReassignOrder возвращает заказ в пул если подрядчик не позвонил
func (s *OrderStorage) ReassignOrder(orderID int) error {
	query := `UPDATE orders SET status = 'confirmed', contractor_id = NULL, accepted_at = NULL WHERE id = $1 AND status = 'in_progress'`
	_, err := s.db.Exec(query, orderID)
	return err
}

// MarkClientUnreachable — подрядчик звонил, но клиент не ответил
func (s *OrderStorage) MarkClientUnreachable(orderID int) error {
	query := `UPDATE orders SET status = 'client_unreachable' WHERE id = $1`
	_, err := s.db.Exec(query, orderID)
	return err
}

// MarkLeadSold — звонок 30+ сек, лид продан
func (s *OrderStorage) MarkLeadSold(orderID int) error {
	query := `UPDATE orders SET status = 'lead_sold' WHERE id = $1`
	_, err := s.db.Exec(query, orderID)
	return err
}

// GetExpiredAcceptedOrders — заказы в статусе in_progress старше 15 минут
func (s *OrderStorage) GetExpiredAcceptedOrders() ([]Order, error) {
	query := `
		SELECT id, client_name, phone, device, COALESCE(brand, ''), problem, 
			zip_code, COALESCE(preferred_time, ''), status, price, contractor_id
		FROM orders
		WHERE status = 'in_progress'
		AND accepted_at < NOW() - INTERVAL '15 minutes'
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.ClientName, &order.Phone, &order.Device,
			&order.Brand, &order.Problem, &order.ZipCode, &order.PreferredTime,
			&order.Status, &order.Price, &order.ContractorID,
		)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}
