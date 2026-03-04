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
