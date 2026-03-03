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
	// SQL запрос
	query := "SELECT id, client_name, phone, device, problem, zip_code, status, price, contractor_id FROM orders"

	// Выполняем запрос
	rows, err := s.db.Query(query)
	if err != nil {
		// Если ошибка — возвращаем пустой список
		return []Order{}
	}
	defer rows.Close()

	var orders []Order

	// Перебираем результаты
	for rows.Next() {
		var order Order

		// Читаем данные из строки
		err := rows.Scan(
			&order.ID,
			&order.ClientName,
			&order.Phone,
			&order.Device,
			&order.Problem,
			&order.ZipCode,
			&order.Status,
			&order.Price,
			&order.ContractorID,
		)

		if err != nil {
			continue // Пропускаем строку с ошибкой
		}

		orders = append(orders, order)
	}

	return orders
}

// Create создаёт новый заказ в базе
func (s *OrderStorage) Create(order Order) Order {
	// SQL запрос для вставки
	query := `
        INSERT INTO orders (client_name, phone, device, problem, zip_code, status, price)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	// Выполняем запрос и получаем ID
	err := s.db.QueryRow(
		query,
		order.ClientName,
		order.Phone,
		order.Device,
		order.Problem,
		order.ZipCode,
		order.Status,
		order.Price,
	).Scan(&order.ID)

	if err != nil {
		// Если ошибка — возвращаем пустой заказ
		return Order{}
	}

	return order
}

// Update обновляет заказ в базе
func (s *OrderStorage) Update(id int, order Order) bool {
	query := `
        UPDATE orders 
        SET client_name = $1, phone = $2, device = $3, problem = $4, 
            zip_code = $5, status = $6, price = $7
        WHERE id = $8
    `

	result, err := s.db.Exec(
		query,
		order.ClientName,
		order.Phone,
		order.Device,
		order.Problem,
		order.ZipCode,
		order.Status,
		order.Price,
		id, // $8 — это ID заказа который обновляем
	)

	if err != nil {
		return false
	}

	// Проверяем сколько строк обновилось
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

	// Проверяем сколько строк удалилось
	rows, _ := result.RowsAffected()
	return rows > 0
}

// GetByID находит заказ по ID
func (s *OrderStorage) GetByID(id int) (*Order, bool) {
	query := "SELECT id, client_name, phone, device, problem, zip_code, status, price, contractor_id FROM orders WHERE id = $1"

	var order Order
	err := s.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.ClientName,
		&order.Phone,
		&order.Device,
		&order.Problem,
		&order.ZipCode,
		&order.Status,
		&order.Price,
		&order.ContractorID,
	)

	if err != nil {
		return nil, false
	}

	return &order, true
}

// AssignContractor назначает подрядчика на заказ
func (s *OrderStorage) AssignContractor(orderID int, contractorID int) error {
	query := `
        UPDATE orders 
        SET status = 'in_progress', contractor_id = $1
        WHERE id = $2
    `

	_, err := s.db.Exec(query, contractorID, orderID)
	return err
}
