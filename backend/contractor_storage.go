package main

import "database/sql"

// ContractorStorage работает с подрядчиками

type ContractorStorage struct {
	db *sql.DB
}

// NewContractorStorage создаёт хранилище
func NewContractorStorage(db *sql.DB) *ContractorStorage {
	return &ContractorStorage{db: db}
}

// Create создаёт нового подрядчика
func (s *ContractorStorage) Create(contractor Contractor) (Contractor, error) {
	query := `
        INSERT INTO contractors (name, email, password_hash, phone)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err := s.db.QueryRow(
		query,
		contractor.Name,
		contractor.Email,
		contractor.PasswordHash,
		contractor.Phone,
	).Scan(&contractor.ID)

	if err != nil {
		return Contractor{}, err
	}

	return contractor, nil
}

// GetByEmail находит подрядчика по email
func (s *ContractorStorage) GetByEmail(email string) (*Contractor, error) {
	query := `
        SELECT id, name, email, password_hash, phone, rating
        FROM contractors
        WHERE email = $1
    `

	var c Contractor
	err := s.db.QueryRow(query, email).Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&c.PasswordHash,
		&c.Phone,
		&c.Rating,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetByID находит подрядчика по ID
func (s *ContractorStorage) GetByID(id int) (*Contractor, error) {
	query := `
        SELECT id, name, email, password_hash, phone, rating
        FROM contractors
        WHERE id = $1
    `

	var c Contractor
	err := s.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&c.PasswordHash,
		&c.Phone,
		&c.Rating,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetAll возвращает всех подрядчиков
func (s *ContractorStorage) GetAll() []Contractor {
	query := "SELECT id, name, email, password_hash, phone, rating FROM contractors"

	rows, err := s.db.Query(query)
	if err != nil {
		return []Contractor{}
	}
	defer rows.Close()

	contractors := []Contractor{}

	for rows.Next() {
		var c Contractor
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.PasswordHash, &c.Phone, &c.Rating)
		if err != nil {
			continue
		}
		contractors = append(contractors, c)
	}

	return contractors
}

// ContractorStats — подрядчик со статистикой для админки
type ContractorStats struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Rating        float64 `json:"rating"`
	CreatedAt     string  `json:"created_at"`
	OrdersTaken   int     `json:"orders_taken"`    // сколько заказов принял (in_progress + lead_sold + completed)
	OrdersSold    int     `json:"orders_sold"`     // сколько лидов продано (lead_sold)
	ActiveOrderID *int    `json:"active_order_id"` // активный заказ прямо сейчас (in_progress)
}

// GetAllWithStats возвращает всех подрядчиков со статистикой
func (s *ContractorStorage) GetAllWithStats() ([]ContractorStats, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.email,
			c.phone,
			c.rating,
			c.created_at,
			COUNT(CASE WHEN o.status IN ('in_progress', 'lead_sold', 'completed') THEN 1 END) AS orders_taken,
			COUNT(CASE WHEN o.status = 'lead_sold' THEN 1 END) AS orders_sold,
			MAX(CASE WHEN o.status = 'in_progress' THEN o.id END) AS active_order_id
		FROM contractors c
		LEFT JOIN orders o ON o.contractor_id = c.id
		GROUP BY c.id, c.name, c.email, c.phone, c.rating, c.created_at
		ORDER BY c.created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ContractorStats
	for rows.Next() {
		var cs ContractorStats
		err := rows.Scan(
			&cs.ID,
			&cs.Name,
			&cs.Email,
			&cs.Phone,
			&cs.Rating,
			&cs.CreatedAt,
			&cs.OrdersTaken,
			&cs.OrdersSold,
			&cs.ActiveOrderID,
		)
		if err != nil {
			continue
		}
		result = append(result, cs)
	}

	return result, nil
}
