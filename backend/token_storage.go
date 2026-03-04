package main

import "database/sql"

// SaveJobToken сохраняет токен в базе
func SaveJobToken(db *sql.DB, orderID int, token string) error {
	query := "INSERT INTO job_tokens (token, order_id) VALUES ($1, $2)"
	_, err := db.Exec(query, token, orderID)
	return err
}

// GetOrderByToken находит заказ по токену
func GetOrderByToken(db *sql.DB, token string) (*Order, error) {
	query := `
        SELECT o.id, o.client_name, o.phone, o.device, o.problem, 
               o.zip_code, o.status, o.price, o.contractor_id
        FROM orders o
        JOIN job_tokens jt ON jt.order_id = o.id
        WHERE jt.token = $1 AND jt.used = FALSE
    `

	var order Order
	err := db.QueryRow(query, token).Scan(
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
		return nil, err
	}

	return &order, nil
}

// MarkTokenUsed помечает токен как использованный
func MarkTokenUsed(db *sql.DB, token string) error {
	query := "UPDATE job_tokens SET used = TRUE WHERE token = $1"
	_, err := db.Exec(query, token)
	return err
}
