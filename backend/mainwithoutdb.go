package main

//
//import (
//"encoding/json"
//"fmt"
//"net/http"
//"os"
//"strconv"
//"strings"
//)
//
//// Order — заказ на ремонт
//
//type Order struct {
//	ID         int    `json:"id"`
//	ClientName string `json:"client_name"`
//	Phone      string `json:"phone"`
//	Device     string `json:"device"`  // Что ремонтируем
//	Problem    string `json:"problem"` // Описание проблемы
//	ZipCode    string `json:"zip_code"`
//	Status     string `json:"status"` // "новый", "в работе", "готов"
//	Price      int    `json:"price"`  // Цена ремонта
//}
//
//type OrderStorage struct {
//	orders   []Order
//	nextID   int
//	filename string
//}
//
//// NewOrderStorage — создаёт хранилище и загружает данные из файла
//
//func newOrderStorage(filename string) *OrderStorage {
//	storage := &OrderStorage{
//		orders:   []Order{},
//		nextID:   1,
//		filename: filename,
//	}
//	// Пытаемся загрузить из файла
//
//	storage.load()
//
//	return storage
//}
//
//// load — загружает заказы из файла
//
//func (s *OrderStorage) load() {
//	// Проверяем, существует ли файл
//
//	_, err := os.Stat(s.filename)
//	if os.IsNotExist(err) {
//		// Файла нет — ничего не делаем
//		fmt.Println("📁 Файл не найден, начинаем с пустого списка")
//		return
//	}
//
//	// Читаем файл
//	data, err := os.ReadFile(s.filename)
//
//	if err != nil {
//		fmt.Println("❌ Ошибка чтения файла:", err)
//		return
//	}
//
//	//	Парсим JSON
//
//	err = json.Unmarshal(data, &s.orders)
//	if err != nil {
//		fmt.Println("❌ Ошибка парсинга JSON:", err)
//		return
//	}
//
//	// Находим максимальный ID для nextID
//
//	for _, order := range s.orders {
//		if order.ID >= s.nextID {
//			s.nextID = order.ID + 1
//		}
//	}
//
//	fmt.Printf("✅ Загружено заказов: %d\n", len(s.orders))
//}
//
//// save — сохраняет заказы в файл
//
//func (s *OrderStorage) save() {
//	data, err := json.MarshalIndent(s.orders, "", "  ")
//	if err != nil {
//		fmt.Println("❌ Ошибка парсинга JSON:", err)
//		return
//	}
//
//	err = os.WriteFile(s.filename, data, 0644)
//	if err != nil {
//		fmt.Println("❌ Ошибка записи файла:", err)
//		return
//	}
//
//	fmt.Println("💾 Данные сохранены")
//}
//
//// GetAll — возвращает все заказы
//
//func (s OrderStorage) GetAll() []Order {
//	return s.orders
//}
//
//// GetByID — находит заказ по ID
//
//func (s OrderStorage) GetByID(id int) (*Order, bool) {
//
//	for i := range s.orders {
//		if s.orders[i].ID == id {
//			return &s.orders[i], true
//		}
//
//	}
//	return nil, false
//}
//
//// Create — создаёт новый заказ
//
//func (s *OrderStorage) Create(order Order) Order {
//	order.ID = s.nextID
//	s.nextID++
//	s.orders = append(s.orders, order)
//	s.save() // Сохраняем после создания
//	return order
//}
//
//// Update — обновляет заказ
//
//func (s *OrderStorage) Update(id int, updated Order) bool {
//	for i := range s.orders {
//		if s.orders[i].ID == id {
//			updated.ID = id // Сохраняем ID
//			s.orders[i] = updated
//			s.save() // Сохраняем после обновления
//			return true
//		}
//	}
//	return false
//}
//
//// Delete — удаляет заказ
//func (s *OrderStorage) Delete(id int) bool {
//	for i := range s.orders {
//		if s.orders[i].ID == id {
//			s.orders = append(s.orders[:i], s.orders[i+1:]...)
//			s.save() // Сохраняем после удаления
//			return true
//		}
//	}
//	return false
//}
//
//func main() {
//	// Создаём хранилище
//	storage := newOrderStorage("orders.json")
//
//	// Регистрируем обработчик
//
//	//fmt.Println("✅ Регистрируем /api/orders/")
//
//	http.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
//
//		//fmt.Println("📍 Обработчик /api/orders/ вызван! URL:", r.URL.Path, "Метод:", r.Method)
//		// CORS
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(200)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//
//		// Извлекаем ID из URL
//		parts := strings.Split(r.URL.Path, "/")
//		if len(parts) < 4 {
//			http.Error(w, "Неверный URL", 400)
//			return
//		}
//
//		id, err := strconv.Atoi(parts[3])
//
//		if err != nil {
//			http.Error(w, "ID должен быть числом", 400)
//			return
//		}
//
//		// PUT - обновить заказ
//
//		if r.Method == "PUT" {
//			var updateData Order
//			err := json.NewDecoder(r.Body).Decode(&updateData)
//
//			if err != nil {
//				http.Error(w, "Неверный формат JSON", 400)
//				return
//			}
//
//			success := storage.Update(id, updateData)
//
//			if !success {
//				http.Error(w, "Заказ не найден", 404)
//				return
//			}
//
//			// Получаем обновлённый заказ
//
//			updated, found := storage.GetByID(id)
//
//			if !found {
//				http.Error(w, "Заказ не найден", 404)
//				return
//			}
//
//			json.NewEncoder(w).Encode(updated)
//			return
//
//		}
//
//		// DELETE - удалить заказ
//
//		if r.Method == "DELETE" {
//			success := storage.Delete(id)
//			if !success {
//				http.Error(w, "Заказ не найден", 404)
//				return
//			}
//
//			w.WriteHeader(204) // 204 No Content
//
//			return
//		}
//		http.Error(w, "Метод не поддерживается", 405)
//
//	})
//
//	//fmt.Println("✅ Регистрируем /api/orders")
//
//	http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
//
//		// CORS заголовки (добавили в самое начало!)
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//
//		// Обработка preflight запроса
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(200)
//			return
//		}
//
//		// Устанавливаем заголовок
//		w.Header().Set("Content-Type", "application/json")
//
//		// GET - получить заказы
//
//		if r.Method == "GET" {
//			// Получаем все заказы
//
//			orders := storage.GetAll()
//
//			// Превращаем в JSON и отправляем
//			json.NewEncoder(w).Encode(orders)
//
//			return
//		}
//
//		if r.Method == "POST" {
//			// Создать новый заказ
//
//			var newOrder Order
//
//			// Читаем JSON из тела запроса
//
//			err := json.NewDecoder(r.Body).Decode(&newOrder)
//			if err != nil {
//				http.Error(w, "Неверный формат JSON", 400)
//				return
//			}
//
//			// Валидация
//			//if newOrder.ClientName == "" {
//			//	http.Error(w, "Имя клиента обязательно", 400)
//			//	return
//			//}
//
//			if newOrder.ClientName == "" || newOrder.Phone == "" || newOrder.Device == "" {
//				http.Error(w, "Имя, телефон и техника обязательны", 400)
//				return
//			}
//
//			// Создаём заказ
//			created := storage.Create(newOrder)
//
//			// Отправляем созданный заказ обратно
//			w.WriteHeader(201) // 201 Created
//			json.NewEncoder(w).Encode(created)
//			return
//		}
//
//		////Методы PUT and DELETE from GPT
//		//if r.Method == "PUT" {
//		//	idStr := r.URL.Query().Get("id")
//		//	if idStr == "" {
//		//		http.Error(w, "Missing id", 400)
//		//		return
//		//	}
//		//
//		//	id, err := strconv.Atoi(idStr)
//		//	if err != nil {
//		//		http.Error(w, "Invalid id", 400)
//		//		return
//		//	}
//		//
//		//	var updated Order
//		//	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
//		//		http.Error(w, "Invalid JSON", 400)
//		//		return
//		//	}
//		//
//		//	ok := storage.Update(id, updated)
//		//	if !ok {
//		//		http.Error(w, "Order not found", 404)
//		//		return
//		//	}
//		//
//		//	json.NewEncoder(w).Encode(updated)
//		//	return
//		//}
//		//
//		//if r.Method == "DELETE" {
//		//	idStr := r.URL.Query().Get("id")
//		//	if idStr == "" {
//		//		http.Error(w, "Missing id", 400)
//		//		return
//		//	}
//		//
//		//	id, err := strconv.Atoi(idStr)
//		//	if err != nil {
//		//		http.Error(w, "Invalid id", 400)
//		//		return
//		//	}
//		//
//		//	ok := storage.Delete(id)
//		//	if !ok {
//		//		http.Error(w, "Order not found", 404)
//		//		return
//		//	}
//		//
//		//	w.WriteHeader(http.StatusNoContent)
//		//	return
//		//}
//		// Если метод не GET и не POST
//		http.Error(w, "Метод не поддерживается", 405)
//
//	})
//
//	// Запускаем сервер
//	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
//	http.ListenAndServe(":8080", nil)
//
//}
//
