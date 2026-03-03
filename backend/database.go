package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// InitDB создаёт подключение к PostgreSQL
func InitDB() (*sql.DB, error) {
	// Строка подключения
	connStr := "host=localhost port=5432 user=apple dbname=fixly sslmode=disable"

	// Открываем подключение
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %v", err)
	}

	// Проверяем что подключение работает
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("не могу связаться с БД: %v", err)
	}

	fmt.Println("✅ Подключение к PostgreSQL установлено")

	return db, nil
}
