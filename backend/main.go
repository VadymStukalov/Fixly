package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

			// 👇 ВОТ СЮДА ВСТАВЛЯЕМ
			if updated.Status == "confirmed" {
				contractors := contractorStorage.GetAll()
				for _, contractor := range contractors {
					token := GenerateToken()
					jobURL := fmt.Sprintf("https://fixly-eta.vercel.app/accept/%s", token)
					message := fmt.Sprintf("New Job - %s repair - %s\nAccept: %s", updated.Device, updated.ZipCode, jobURL)

					err := SaveJobToken(db, updated.ID, contractor.ID, token)
					if err != nil {
						fmt.Printf("❌ Failed to save token for contractor #%d: %v\n", contractor.ID, err)
						continue
					}

					err = SendSMS(contractor.Phone, message)
					if err != nil {
						fmt.Printf("❌ Failed to send SMS to %s: %v\n", contractor.Phone, err)
					} else {
						fmt.Printf("✅ SMS sent to contractor #%d (%s)\n", contractor.ID, contractor.Phone)
					}
				}
			}
			// 👆 конец вставки

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

	// GET /accept/{token} — получить детали заказа по токену
	http.HandleFunc("/accept/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Извлекаем токен из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.Error(w, "Invalid URL", 400)
			return
		}
		token := parts[2]

		if r.Method == "GET" {
			order, _, _, err := GetOrderByToken(db, token)
			if err != nil {
				fmt.Printf("❌ GetOrderByToken error: %v\n", err)
				http.Error(w, err.Error(), 404)
				return
			}
			json.NewEncoder(w).Encode(order)
			return
		}

		if r.Method == "POST" {
			order, contractorID, contractorPhone, err := GetOrderByToken(db, token)
			if err != nil {
				fmt.Printf("❌ GetOrderByToken error: %v\n", err)
				http.Error(w, "Token not found or already used", 404)
				return
			}

			if order.Status != "confirmed" {
				http.Error(w, "Order already taken", 400)
				return
			}

			err = storage.AssignContractor(order.ID, contractorID)
			if err != nil {
				http.Error(w, "Failed to assign contractor", 500)
				return
			}

			MarkTokenUsed(db, token)

			fmt.Printf("✅ Order #%d accepted by contractor #%d via token\n", order.ID, contractorID)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success":          true,
				"contractor_phone": contractorPhone,
			})
			return
		}

		http.Error(w, "Method not supported", 405)
	})

	// POST /api/call — инициируем звонок подрядчику через Twilio
	http.HandleFunc("/api/call", func(w http.ResponseWriter, r *http.Request) {
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

		var data struct {
			ContractorPhone string `json:"contractor_phone"`
			ClientPhone     string `json:"client_phone"`
			OrderID         int    `json:"order_id"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data.ContractorPhone == "" || data.ClientPhone == "" {
			http.Error(w, "contractor_phone and client_phone required", 400)
			return
		}

		err = InitiateCall(data.ContractorPhone, data.ClientPhone, data.OrderID)
		if err != nil {
			fmt.Printf("❌ Call error: %v\n", err)
			http.Error(w, "Failed to initiate call", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	})

	// POST /api/twiml — TwiML инструкция для Twilio (соединить с клиентом)
	http.HandleFunc("/api/twiml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")

		clientPhone := r.URL.Query().Get("client_phone")
		orderID := r.URL.Query().Get("order_id")

		fmt.Printf("📞 TwiML called, connecting to client %s (order %s)\n", clientPhone, orderID)

		w.Write([]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say>Connecting you to your client now.</Say>
    <Dial callerId="%s">%s</Dial>
</Response>`, TWILIO_PHONE, clientPhone)))
	})

	// Запускаем сервер
	//fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	//http.ListenAndServe(":8080", nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🚀 Сервер запущен на порту", port)
	http.ListenAndServe(":"+port, nil)

}
