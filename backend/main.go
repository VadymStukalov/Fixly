package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Order — заказ на ремонт

type Order struct {
	ID           int    `json:"id"`
	ClientName   string `json:"client_name"`
	Phone        string `json:"phone"`
	Device       string `json:"device"`  // Что ремонтируем
	Problem      string `json:"problem"` // Описание проблемы
	ZipCode      string `json:"zip_code"`
	Status       string `json:"status"` // "новый", "в работе", "готов"
	Price        int    `json:"price"`  // Цена ремонта
	ContractorID *int   `json:"contractor_id"`
}

// Contractor — подрядчик

type Contractor struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"` // Не отправляем пароль в JSON
	Phone        string  `json:"phone"`
	Rating       float64 `json:"rating"`
}

// Bid — ставка подрядчика на заказ

type Bid struct {
	ID           int    `json:"id"`
	OrderID      int    `json:"order_id"`
	ContractorID int    `json:"contractor_id"`
	ProposedTime string `json:"proposed_time"` // "Сегодня", "Завтра", "Через 2 дня"
}

func main() {
	// Подключаемся к базе данных
	db, err := InitDB() // Подключились к PostgreSQL
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return
	}
	defer db.Close()

	// Создаём хранилище
	storage := NewOrderStorage(db)                // Создали storage с базой
	contractorStorage := NewContractorStorage(db) // Добавили
	bidStorage := NewBidStorage(db)               // Добавили

	http.HandleFunc("/api/orders/available", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			http.Error(w, "Метод не поддерживается", 405)
			return
		}

		// Получаем все заказы

		allOrders := storage.GetAll()

		// Фильтруем только "подтверждённые"

		var availableOrders []Order

		for _, order := range allOrders {
			if order.Status == "confirmed" {
				availableOrders = append(availableOrders, order)
			}
		}

		json.NewEncoder(w).Encode(availableOrders)

	})

	http.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {

		//fmt.Println("📍 Обработчик /api/orders/ вызван! URL:", r.URL.Path, "Метод:", r.Method)
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Извлекаем ID
		parts := strings.Split(r.URL.Path, "/")

		// Извлекаем ID

		// Проверяем /complete endpoint
		if len(parts) == 5 && parts[4] == "complete" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")

			idStr := parts[3]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID", 400)
				return
			}

			// Меняем статус на completed
			order, found := storage.GetByID(id)
			if !found {
				http.Error(w, "Order not found", 404)
				return
			}

			order.Status = "completed"
			success := storage.Update(id, *order)

			if !success {
				http.Error(w, "Failed to update order", 500)
				return
			}

			fmt.Printf("✅ Order #%d marked as completed\n", id)
			json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		// Дальше идёт ваш существующий код для PUT/DELETE
		if r.Method == "PUT" {
			// ...
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Извлекаем ID из URL
		//parts := strings.Split(r.URL.Path, "/")
		//if len(parts) < 4 {
		//	http.Error(w, "Неверный URL", 400)
		//	return
		//}

		id, err := strconv.Atoi(parts[3])

		if err != nil {
			http.Error(w, "ID должен быть числом", 400)
			return
		}

		// PUT - обновить заказ

		if r.Method == "PUT" {
			var updateData Order
			err := json.NewDecoder(r.Body).Decode(&updateData)

			if err != nil {
				http.Error(w, "Неверный формат JSON", 400)
				return
			}

			success := storage.Update(id, updateData)

			if !success {
				http.Error(w, "Заказ не найден", 404)
				return
			}

			// Получаем обновлённый заказ

			updated, found := storage.GetByID(id)

			if !found {
				http.Error(w, "Заказ не найден", 404)
				return
			}

			json.NewEncoder(w).Encode(updated)
			return

		}

		// DELETE - удалить заказ

		if r.Method == "DELETE" {
			success := storage.Delete(id)
			if !success {
				http.Error(w, "Заказ не найден", 404)
				return
			}

			w.WriteHeader(204) // 204 No Content

			return
		}
		http.Error(w, "Метод не поддерживается", 405)

	})

	//fmt.Println("✅ Регистрируем /api/orders")

	http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {

		// CORS заголовки (добавили в самое начало!)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обработка preflight запроса
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		// Устанавливаем заголовок
		w.Header().Set("Content-Type", "application/json")

		// GET - получить заказы

		if r.Method == "GET" {
			// Получаем все заказы

			orders := storage.GetAll()

			// Превращаем в JSON и отправляем
			json.NewEncoder(w).Encode(orders)

			return
		}

		if r.Method == "POST" {
			// Создать новый заказ

			var newOrder Order

			// Читаем JSON из тела запроса

			err := json.NewDecoder(r.Body).Decode(&newOrder)
			if err != nil {
				http.Error(w, "Неверный формат JSON", 400)
				return
			}

			// Валидация
			//if newOrder.ClientName == "" {
			//	http.Error(w, "Имя клиента обязательно", 400)
			//	return
			//}

			if newOrder.ClientName == "" || newOrder.Phone == "" || newOrder.Device == "" {
				http.Error(w, "Имя, телефон и техника обязательны", 400)
				return
			}

			// Создаём заказ
			created := storage.Create(newOrder)

			// Отправляем созданный заказ обратно
			w.WriteHeader(201) // 201 Created
			json.NewEncoder(w).Encode(created)
			return
		}

		////Методы PUT and DELETE from GPT
		//if r.Method == "PUT" {
		//	idStr := r.URL.Query().Get("id")
		//	if idStr == "" {
		//		http.Error(w, "Missing id", 400)
		//		return
		//	}
		//
		//	id, err := strconv.Atoi(idStr)
		//	if err != nil {
		//		http.Error(w, "Invalid id", 400)
		//		return
		//	}
		//
		//	var updated Order
		//	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		//		http.Error(w, "Invalid JSON", 400)
		//		return
		//	}
		//
		//	ok := storage.Update(id, updated)
		//	if !ok {
		//		http.Error(w, "Order not found", 404)
		//		return
		//	}
		//
		//	json.NewEncoder(w).Encode(updated)
		//	return
		//}
		//
		//if r.Method == "DELETE" {
		//	idStr := r.URL.Query().Get("id")
		//	if idStr == "" {
		//		http.Error(w, "Missing id", 400)
		//		return
		//	}
		//
		//	id, err := strconv.Atoi(idStr)
		//	if err != nil {
		//		http.Error(w, "Invalid id", 400)
		//		return
		//	}
		//
		//	ok := storage.Delete(id)
		//	if !ok {
		//		http.Error(w, "Order not found", 404)
		//		return
		//	}
		//
		//	w.WriteHeader(http.StatusNoContent)
		//	return
		//}
		// Если метод не GET и не POST
		http.Error(w, "Метод не поддерживается", 405)

	})

	// POST /api/contractors/register — регистрация подрядчика

	http.HandleFunc("/api/contractors/register", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.Error(w, "Метод не поддерживается", 405)
		}

		// Читаем данные
		var data struct {
			Name     string
			Email    string
			Password string
			Phone    string
		}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			if err != nil {
				http.Error(w, "Неверный формат JSON", 400)
				return
			}
		}

		// Валидация

		if data.Name == "" || data.Email == "" || data.Password == "" {
			http.Error(w, "Имя, email и пароль обязательны", 400)
			return
		}

		// Проверяем что email не занят

		existing, _ := contractorStorage.GetByEmail(data.Email)

		if existing != nil {
			http.Error(w, "Email уже зарегистрирован", 400)
		}

		// Хешируем пароль (пока просто сохраняем как есть — небезопасно!)
		// TODO: использовать bcrypt для хеширования

		// Создаём подрядчика

		contractor := Contractor{
			Name:         data.Name,
			Email:        data.Email,
			PasswordHash: data.Password, // В реальности нужно хешировать!
			Phone:        data.Phone,
		}

		created, err := contractorStorage.Create(contractor)

		if err != nil {
			http.Error(w, "Ошибка создания подрядчика", 500)
			return
		}
		// Возвращаем созданного подрядчика
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(created)

	})

	// GET /api/orders/available — доступные заказы для подрядчиков

	// POST /api/bids — создать ставку на заказ

	http.HandleFunc("/api/bids", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.Error(w, "Method not supported", 405)
			return
		}

		// Читаем данные
		var data struct {
			OrderID      int    `json:"order_id"`
			ContractorID int    `json:"contractor_id"`
			ProposedTime string `json:"proposed_time"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON format", 400)
			return
		}

		// Валидация
		if data.OrderID == 0 || data.ContractorID == 0 || data.ProposedTime == "" {
			http.Error(w, "order_id, contractor_id and proposed_time required", 400)
			return
		}

		// Проверяем что заказ существует
		order, found := storage.GetByID(data.OrderID)
		if !found {
			http.Error(w, "Order not found", 404)
			return
		}

		// Проверяем статус заказа
		if order.Status != "confirmed" {
			http.Error(w, "This order is not available for bidding", 400)
			return
		}

		// Проверяем что подрядчик ещё не делал ставку
		hasBid, _ := bidStorage.HasBid(data.OrderID, data.ContractorID)
		if hasBid {
			http.Error(w, "You already placed a bid on this order", 400)
			return
		}

		// Создаём ставку
		bid := Bid{
			OrderID:      data.OrderID,
			ContractorID: data.ContractorID,
			ProposedTime: data.ProposedTime,
		}

		created, err := bidStorage.Create(bid)
		if err != nil {
			http.Error(w, "Error creating bid", 500)
			return
		}

		// Если ставка "Today" — назначаем СРАЗУ
		if created.ProposedTime == "Today" {
			err := storage.AssignContractor(created.OrderID, created.ContractorID)
			if err != nil {
				fmt.Println("❌ Error assigning contractor:", err)
			} else {
				fmt.Printf("✅ Order #%d assigned to contractor #%d (Today — instant)\n", created.OrderID, created.ContractorID)
			}

			w.WriteHeader(201)
			json.NewEncoder(w).Encode(created)
			return
		}

		// Если НЕ "Today" — проверяем есть ли уже таймер
		// Запускаем таймер только для ПЕРВОЙ ставки
		bids, _ := bidStorage.GetByOrderID(created.OrderID)
		if len(bids) == 1 {
			// Это первая ставка — запускаем таймер
			delay := 30 * time.Second // 30 секунд для теста (потом измените на time.Minute)
			bidStorage.ScheduleSelection(created.OrderID, storage, delay)
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(created)
	})
	//http.HandleFunc("/api/bids", func(w http.ResponseWriter, r *http.Request) {
	//	// CORS
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//
	//	if r.Method == "OPTIONS" {
	//		w.WriteHeader(200)
	//		return
	//	}
	//
	//	w.Header().Set("Content-Type", "application/json")
	//
	//	if r.Method != "POST" {
	//		http.Error(w, "Метод не поддерживается", 405)
	//		return
	//	}
	//
	//	// Читаем данные
	//	var data struct {
	//		OrderID      int    `json:"order_id"`
	//		ContractorID int    `json:"contractor_id"`
	//		ProposedTime string `json:"proposed_time"`
	//	}
	//	err := json.NewDecoder(r.Body).Decode(&data)
	//
	//	if err != nil {
	//		http.Error(w, "Неверный формат JSON", 400)
	//		return
	//	}
	//
	//	// Валидация
	//
	//	if data.OrderID == 0 || data.ContractorID == 0 || data.ProposedTime == "" {
	//		http.Error(w, "order_id, contractor_id и proposed_time обязательны", 400)
	//		return
	//	}
	//
	//	// Проверяем что заказ существует
	//
	//	order, found := storage.GetByID(data.OrderID)
	//
	//	if !found {
	//		http.Error(w, "Заказ не найден", 404)
	//		return
	//	}
	//
	//	// Проверяем статус заказа
	//
	//	if order.Status != "confirmed" {
	//		http.Error(w, "Этот заказ недоступен для ставок", 400)
	//		return
	//	}
	//
	//	// Проверяем что подрядчик ещё не делал ставку
	//
	//	hasBid, _ := bidStorage.HasBid(data.OrderID, data.ContractorID)
	//
	//	if hasBid {
	//		http.Error(w, "Вы уже сделали ставку на этот заказ", 400)
	//		return
	//	}
	//
	//	// Создаём ставку
	//
	//	bid := Bid{
	//		OrderID:      data.OrderID,
	//		ContractorID: data.ContractorID,
	//		ProposedTime: data.ProposedTime,
	//	}
	//
	//	created, err := bidStorage.Create(bid)
	//
	//	if err != nil {
	//		http.Error(w, "Ошибка создания ставки", 500)
	//		return
	//	}
	//
	//	// Проверяем статус заказа
	//	order, found := storage.GetByID(created.OrderID)
	//	if !found {
	//		w.WriteHeader(201)
	//		json.NewEncoder(w).Encode(created)
	//		return
	//	}
	//
	//	// Если заказ не "confirmed" — ничего не делаем
	//	if order.Status != "confirmed" {
	//		w.WriteHeader(201)
	//		json.NewEncoder(w).Encode(created)
	//		return
	//	}
	//
	//	// Если ставка "Today" — назначаем СРАЗУ
	//	if created.ProposedTime == "Today" {
	//		err := storage.AssignContractor(created.OrderID, created.ContractorID)
	//		if err != nil {
	//			fmt.Println("❌ Ошибка назначения подрядчика:", err)
	//		} else {
	//			fmt.Printf("✅ Заказ #%d назначен подрядчику #%d (Today — мгновенно)\n", created.OrderID, created.ContractorID)
	//		}
	//
	//		w.WriteHeader(201)
	//		json.NewEncoder(w).Encode(created)
	//		return
	//	}
	//
	//	// Если НЕ "Today" — проверяем есть ли уже таймер
	//	// Для MVP: запускаем таймер только для ПЕРВОЙ ставки
	//	bids, _ := bidStorage.GetByOrderID(created.OrderID)
	//	if len(bids) == 1 {
	//		// Это первая ставка — запускаем таймер
	//		delay := 30 * time.Minute // 30 минут
	//		bidStorage.ScheduleSelection(created.OrderID, storage, delay)
	//	}
	//
	//	w.WriteHeader(201)
	//
	//	json.NewEncoder(w).Encode(created)
	//
	//})

	// POST /api/contractors/login — вход подрядчика
	http.HandleFunc("/api/contractors/login", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.Error(w, "Method not supported", 405)
			return
		}

		// Читаем данные
		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON format", 400)
			return
		}

		// Ищем подрядчика по email
		contractor, err := contractorStorage.GetByEmail(data.Email)
		if err != nil {
			http.Error(w, "Invalid email or password", 401)
			return
		}

		// Проверяем пароль (пока просто сравниваем, без хеша)
		if contractor.PasswordHash != data.Password {
			http.Error(w, "Invalid email or password", 401)
			return
		}

		// Возвращаем данные подрядчика
		json.NewEncoder(w).Encode(contractor)
	})

	// GET /api/contractors/{id}/bids — получить все ставки подрядчика
	http.HandleFunc("/api/contractors/", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			http.Error(w, "Method not supported", 405)
			return
		}

		// Извлекаем ID из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 || parts[4] != "bids" {
			http.Error(w, "Invalid URL", 400)
			return
		}

		contractorID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid contractor ID", 400)
			return
		}

		// Получаем ставки подрядчика
		bids, err := bidStorage.GetByContractorID(contractorID)
		if err != nil {
			http.Error(w, "Error fetching bids", 500)
			return
		}

		// Для каждой ставки добавляем информацию о заказе
		type BidWithOrder struct {
			Bid   Bid   `json:"bid"`
			Order Order `json:"order"`
		}

		var result []BidWithOrder

		for _, bid := range bids {
			order, found := storage.GetByID(bid.OrderID)
			if found {
				result = append(result, BidWithOrder{
					Bid:   bid,
					Order: *order,
				})
			}
		}

		json.NewEncoder(w).Encode(result)
	})

	// Запускаем сервер
	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
