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
	Device       string `json:"device"`
	Problem      string `json:"problem"`
	ZipCode      string `json:"zip_code"`
	Status       string `json:"status"`
	Price        int    `json:"price"`
	ContractorID *int   `json:"contractor_id"`
}

// Contractor — подрядчик
type Contractor struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	Phone        string  `json:"phone"`
	Rating       float64 `json:"rating"`
}

// Bid — ставка подрядчика на заказ
type Bid struct {
	ID           int    `json:"id"`
	OrderID      int    `json:"order_id"`
	ContractorID int    `json:"contractor_id"`
	ProposedTime string `json:"proposed_time"`
}

func main() {
	db, err := InitDB()
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return
	}
	defer db.Close()

	storage := NewOrderStorage(db)
	contractorStorage := NewContractorStorage(db)
	bidStorage := NewBidStorage(db)
	callLogStorage := NewCallLogStorage(db)

	// ─── Фоновый воркер: каждую минуту проверяет просроченные заказы ───────────
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			expiredOrders, err := storage.GetExpiredAcceptedOrders()
			if err != nil {
				fmt.Println("❌ Worker error:", err)
				continue
			}
			for _, order := range expiredOrders {
				// Проверяем были ли вообще звонки
				hadCall, _ := callLogStorage.HasAnyCallAttempt(order.ID)
				if hadCall {
					// Звонил, но клиент не ответил (нет 30+ сек) — ручная проверка
					storage.MarkClientUnreachable(order.ID)
					fmt.Printf("📵 Order #%d → client_unreachable (called but no answer)\n", order.ID)
				} else {
					// Вообще не звонил — возвращаем в пул
					storage.ReassignOrder(order.ID)
					fmt.Printf("🔄 Order #%d → reassigned (no call attempt in 15 min)\n", order.ID)
				}
			}
		}
	}()
	// ────────────────────────────────────────────────────────────────────────────

	http.HandleFunc("/api/orders/available", func(w http.ResponseWriter, r *http.Request) {
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

		allOrders := storage.GetAll()

		var availableOrders []Order
		for _, order := range allOrders {
			if order.Status == "confirmed" {
				availableOrders = append(availableOrders, order)
			}
		}

		json.NewEncoder(w).Encode(availableOrders)
	})

	http.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		parts := strings.Split(r.URL.Path, "/")

		// POST /api/orders/{id}/complete
		if len(parts) == 5 && parts[4] == "complete" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")

			id, err := strconv.Atoi(parts[3])
			if err != nil {
				http.Error(w, "Invalid ID", 400)
				return
			}

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

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

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

			updated, found := storage.GetByID(id)
			if !found {
				http.Error(w, "Заказ не найден", 404)
				return
			}

			// Рассылка SMS при смене статуса на confirmed
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

					// Задержка между SMS
					time.Sleep(200 * time.Millisecond)
				}
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
			w.WriteHeader(204)
			return
		}

		http.Error(w, "Метод не поддерживается", 405)
	})

	http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			orders := storage.GetAll()
			json.NewEncoder(w).Encode(orders)
			return
		}

		if r.Method == "POST" {
			var newOrder Order
			err := json.NewDecoder(r.Body).Decode(&newOrder)
			if err != nil {
				http.Error(w, "Неверный формат JSON", 400)
				return
			}

			if newOrder.ClientName == "" || newOrder.Phone == "" || newOrder.Device == "" {
				http.Error(w, "Имя, телефон и техника обязательны", 400)
				return
			}

			created := storage.Create(newOrder)
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(created)
			return
		}

		http.Error(w, "Метод не поддерживается", 405)
	})

	// POST /api/contractors/register
	http.HandleFunc("/api/contractors/register", func(w http.ResponseWriter, r *http.Request) {
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
			return
		}

		var data struct {
			Name     string
			Email    string
			Password string
			Phone    string
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Неверный формат JSON", 400)
			return
		}

		if data.Name == "" || data.Email == "" || data.Password == "" {
			http.Error(w, "Имя, email и пароль обязательны", 400)
			return
		}

		existing, _ := contractorStorage.GetByEmail(data.Email)
		if existing != nil {
			http.Error(w, "Email уже зарегистрирован", 400)
			return
		}

		contractor := Contractor{
			Name:         data.Name,
			Email:        data.Email,
			PasswordHash: data.Password,
			Phone:        data.Phone,
		}

		created, err := contractorStorage.Create(contractor)
		if err != nil {
			http.Error(w, "Ошибка создания подрядчика", 500)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(created)
	})

	// POST /api/bids
	http.HandleFunc("/api/bids", func(w http.ResponseWriter, r *http.Request) {
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
			OrderID      int    `json:"order_id"`
			ContractorID int    `json:"contractor_id"`
			ProposedTime string `json:"proposed_time"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON format", 400)
			return
		}

		if data.OrderID == 0 || data.ContractorID == 0 || data.ProposedTime == "" {
			http.Error(w, "order_id, contractor_id and proposed_time required", 400)
			return
		}

		order, found := storage.GetByID(data.OrderID)
		if !found {
			http.Error(w, "Order not found", 404)
			return
		}

		if order.Status != "confirmed" {
			http.Error(w, "This order is not available for bidding", 400)
			return
		}

		hasBid, _ := bidStorage.HasBid(data.OrderID, data.ContractorID)
		if hasBid {
			http.Error(w, "You already placed a bid on this order", 400)
			return
		}

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

		bids, _ := bidStorage.GetByOrderID(created.OrderID)
		if len(bids) == 1 {
			delay := 30 * time.Second
			bidStorage.ScheduleSelection(created.OrderID, storage, delay)
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(created)
	})

	// POST /api/contractors/login
	http.HandleFunc("/api/contractors/login", func(w http.ResponseWriter, r *http.Request) {
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
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON format", 400)
			return
		}

		contractor, err := contractorStorage.GetByEmail(data.Email)
		if err != nil {
			http.Error(w, "Invalid email or password", 401)
			return
		}

		if contractor.PasswordHash != data.Password {
			http.Error(w, "Invalid email or password", 401)
			return
		}

		json.NewEncoder(w).Encode(contractor)
	})

	// GET /api/contractors/{id}/bids
	http.HandleFunc("/api/contractors/", func(w http.ResponseWriter, r *http.Request) {
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

		bids, err := bidStorage.GetByContractorID(contractorID)
		if err != nil {
			http.Error(w, "Error fetching bids", 500)
			return
		}

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

	// GET/POST /accept/{token}
	http.HandleFunc("/accept/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		w.Header().Set("Content-Type", "application/json")

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

			// Атомарный захват заказа — защита от race condition
			accepted, err := storage.AcceptOrder(order.ID, contractorID)
			if err != nil {
				http.Error(w, "Failed to accept order", 500)
				return
			}
			if !accepted {
				http.Error(w, "Order already taken", 400)
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

	// POST /api/call — инициируем звонок через Twilio
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
			ContractorID    int    `json:"contractor_id"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data.ContractorPhone == "" || data.ClientPhone == "" {
			http.Error(w, "contractor_phone, client_phone required", 400)
			return
		}

		// Передаём contractorID в InitiateCall (нужен для StatusCallback)
		err = InitiateCall(data.ContractorPhone, data.ClientPhone, data.OrderID, data.ContractorID)
		if err != nil {
			fmt.Printf("❌ Call error: %v\n", err)
			http.Error(w, "Failed to initiate call", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	})

	// POST /api/call-status — webhook от Twilio после завершения звонка
	http.HandleFunc("/api/call-status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(200)
			return
		}

		// Читаем параметры из URL (order_id и contractor_id мы передали в StatusCallback URL)
		orderIDStr := r.URL.Query().Get("order_id")
		contractorIDStr := r.URL.Query().Get("contractor_id")

		orderID, err1 := strconv.Atoi(orderIDStr)
		contractorID, err2 := strconv.Atoi(contractorIDStr)
		if err1 != nil || err2 != nil {
			fmt.Println("❌ call-status: invalid order_id or contractor_id")
			w.WriteHeader(200)
			return
		}

		// Twilio присылает эти поля в теле POST запроса (form-encoded)
		callSID := r.FormValue("CallSid")
		callStatus := r.FormValue("CallStatus")
		durationStr := r.FormValue("CallDuration")

		duration, _ := strconv.Atoi(durationStr)

		fmt.Printf("📞 Call status webhook: order #%d, contractor #%d, status=%s, duration=%ds\n",
			orderID, contractorID, callStatus, duration)

		// Сохраняем лог звонка
		err := callLogStorage.SaveCallLog(orderID, contractorID, callSID, duration, callStatus)
		if err != nil {
			fmt.Printf("❌ Failed to save call log: %v\n", err)
		}

		// Если звонок длился 30+ секунд — лид продан
		if duration >= 30 {
			err := storage.MarkLeadSold(orderID)
			if err != nil {
				fmt.Printf("❌ Failed to mark lead sold: %v\n", err)
			} else {
				fmt.Printf("💰 Order #%d → lead_sold (call duration: %ds)\n", orderID, duration)
			}
		}

		// Twilio ждёт 200 OK
		w.WriteHeader(200)
	})

	// POST /api/twiml — TwiML инструкция для Twilio
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🚀 Сервер запущен на порту", port)
	http.ListenAndServe(":"+port, nil)
}
