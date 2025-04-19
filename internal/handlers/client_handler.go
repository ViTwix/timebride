package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/uptrace/bun"

	"timebride/internal/logger"
	"timebride/internal/models"
	"timebride/internal/utils"
)

type ClientHandler struct {
	DB     *bun.DB
	Logger *logger.Logger
}

func NewClientHandler(db *bun.DB, logger *logger.Logger) *ClientHandler {
	return &ClientHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *ClientHandler) RegisterRoutes(r *mux.Router, auth func(http.Handler) http.Handler) {
	clientsRouter := r.PathPrefix("/clients").Subrouter()
	clientsRouter.Use(auth)

	clientsRouter.HandleFunc("", h.List).Methods(http.MethodGet)
	clientsRouter.HandleFunc("/new", h.NewPage).Methods(http.MethodGet)
	clientsRouter.HandleFunc("", h.Create).Methods(http.MethodPost)
	clientsRouter.HandleFunc("/{id:[0-9]+}", h.Get).Methods(http.MethodGet)
	clientsRouter.HandleFunc("/{id:[0-9]+}/edit", h.EditPage).Methods(http.MethodGet)
	clientsRouter.HandleFunc("/{id:[0-9]+}", h.Update).Methods(http.MethodPut)
	clientsRouter.HandleFunc("/{id:[0-9]+}/delete", h.Delete).Methods(http.MethodPost)
	clientsRouter.HandleFunc("/import", h.Import).Methods(http.MethodPost)
}

// List обробляє запит на отримання списку клієнтів з можливістю фільтрації
func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)

	// Параметри запиту
	query := r.URL.Query().Get("query")
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	// Кількість записів на сторінці
	perPage := 10
	offset := (page - 1) * perPage

	// Основний запит для вибірки клієнтів
	dbQuery := h.DB.NewSelect().
		Model((*models.Client)(nil)).
		Where("user_id = ?", userId)

	// Додаємо умови фільтрації
	if query != "" {
		dbQuery = dbQuery.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR company ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}

	if status == "active" {
		dbQuery = dbQuery.Where("is_active = ?", true)
	} else if status == "inactive" {
		dbQuery = dbQuery.Where("is_active = ?", false)
	}

	// Підрахунок загальної кількості клієнтів для пагінації
	var count int
	_, err := dbQuery.Clone().Count(r.Context())
	if err != nil {
		h.Logger.Error("Error counting clients", err)
		http.Error(w, "Помилка отримання даних", http.StatusInternalServerError)
		return
	}

	// Сортування
	switch sort {
	case "name_asc":
		dbQuery = dbQuery.Order("full_name ASC")
	case "name_desc":
		dbQuery = dbQuery.Order("full_name DESC")
	case "created_asc":
		dbQuery = dbQuery.Order("created_at ASC")
	case "created_desc", "": // За замовчуванням сортування за датою створення за спаданням
		dbQuery = dbQuery.Order("created_at DESC")
	}

	// Пагінація
	dbQuery = dbQuery.Limit(perPage).Offset(offset)

	// Виконання запиту
	var clients []models.Client
	if err := dbQuery.Scan(r.Context(), &clients); err != nil {
		h.Logger.Error("Error fetching clients", err)
		http.Error(w, "Помилка отримання даних", http.StatusInternalServerError)
		return
	}

	// Підготовка даних для відображення
	type ClientViewModel struct {
		ID                 int64     `json:"id"`
		FullName           string    `json:"fullName"`
		Email              string    `json:"email"`
		Phone              string    `json:"phone"`
		Company            string    `json:"company"`
		Category           string    `json:"category"`
		IsActive           bool      `json:"isActive"`
		Avatar             string    `json:"avatar"`
		CreatedAt          time.Time `json:"createdAt"`
		CreatedAtFormatted string    `json:"createdAtFormatted"`
		BookingsCount      int       `json:"bookingsCount"`
		LastBookingDate    string    `json:"lastBookingDate"`
	}

	// Завантаження додаткових даних для кожного клієнта
	clientViewModels := make([]ClientViewModel, 0, len(clients))
	for _, client := range clients {
		// Кількість бронювань та дата останнього бронювання
		var bookingsCount int
		_, err := h.DB.NewSelect().
			Model((*models.Booking)(nil)).
			Where("client_id = ?", client.ID).
			Count(r.Context())
		if err != nil {
			h.Logger.Error("Error counting bookings for client", err, "client_id", client.ID)
			bookingsCount = 0
		}

		// Отримуємо останнє бронювання
		lastBookingDate := ""
		if bookingsCount > 0 {
			var lastBooking models.Booking
			err := h.DB.NewSelect().
				Model(&lastBooking).
				Where("client_id = ?", client.ID).
				Order("start_time DESC").
				Limit(1).
				Scan(r.Context())

			if err == nil {
				lastBookingDate = lastBooking.StartTime.Format("02.01.2006")
			}
		}

		// Форматуємо дату створення
		createdAtFormatted := client.CreatedAt.Format("02.01.2006")

		// Формуємо об'єкт для відображення
		clientVM := ClientViewModel{
			ID:                 client.ID,
			FullName:           client.FullName,
			Email:              client.Email,
			Phone:              client.Phone,
			Company:            client.Company,
			Category:           client.Category,
			IsActive:           client.IsActive,
			Avatar:             client.Avatar,
			CreatedAt:          client.CreatedAt,
			CreatedAtFormatted: createdAtFormatted,
			BookingsCount:      bookingsCount,
			LastBookingDate:    lastBookingDate,
		}

		clientViewModels = append(clientViewModels, clientVM)
	}

	// Підготовка пагінації
	totalPages := (count + perPage - 1) / perPage
	pagination := struct {
		CurrentPage  int    `json:"currentPage"`
		TotalPages   int    `json:"totalPages"`
		TotalItems   int    `json:"totalItems"`
		ItemsPerPage int    `json:"itemsPerPage"`
		Pages        []int  `json:"pages"`
		PrevPageURL  string `json:"prevPageURL"`
		NextPageURL  string `json:"nextPageURL"`
		GetPageURL   func(int) string
	}{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   count,
		ItemsPerPage: perPage,
		GetPageURL: func(p int) string {
			return fmt.Sprintf("/clients?page=%d&query=%s&category=%s&status=%s&sort=%s", p, query, category, status, sort)
		},
	}

	// Формуємо список сторінок для навігації
	pagination.Pages = make([]int, 0)
	startPage := page - 2
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + 4
	if endPage > totalPages {
		endPage = totalPages
		startPage = endPage - 4
		if startPage < 1 {
			startPage = 1
		}
	}
	for i := startPage; i <= endPage; i++ {
		pagination.Pages = append(pagination.Pages, i)
	}

	// URL для попередньої та наступної сторінок
	if page > 1 {
		pagination.PrevPageURL = pagination.GetPageURL(page - 1)
	}
	if page < totalPages {
		pagination.NextPageURL = pagination.GetPageURL(page + 1)
	}

	// Передаємо дані в шаблон
	data := map[string]interface{}{
		"Clients":    clientViewModels,
		"Pagination": pagination,
		"Query":      query,
		"Category":   category,
		"Status":     status,
		"Sort":       sort,
		"Categories": models.DefaultCategories(),
		"CSRFToken":  "dummy_token", // В реальному додатку тут має бути справжній CSRF токен
	}

	// Перетворюємо дані в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.Logger.Error("Error marshaling JSON", err)
		http.Error(w, "Помилка підготовки даних", http.StatusInternalServerError)
		return
	}

	// Встановлюємо заголовки та повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// NewPage відображає форму для створення нового клієнта
func (h *ClientHandler) NewPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Новий клієнт",
		"Action":     "create",
		"Categories": models.DefaultCategories(),
		"Sources":    models.DefaultSources(),
		"CSRFToken":  "dummy_token", // В реальному додатку тут має бути справжній CSRF токен
	}

	// Перетворюємо дані в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.Logger.Error("Error marshaling JSON", err)
		http.Error(w, "Помилка підготовки даних", http.StatusInternalServerError)
		return
	}

	// Встановлюємо заголовки та повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Create обробляє запит на створення нового клієнта
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)

	// Парсимо форму
	err := r.ParseMultipartForm(10 << 20) // 10 MB максимальний розмір файлу
	if err != nil {
		h.Logger.Error("Error parsing form", err)
		http.Error(w, "Помилка обробки форми", http.StatusBadRequest)
		return
	}

	// Отримуємо дані з форми
	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	company := r.FormValue("company")
	category := r.FormValue("category")
	source := r.FormValue("source")
	address := r.FormValue("address")
	notes := r.FormValue("notes")
	isActive := r.FormValue("is_active") == "true" || r.FormValue("is_active") == "on"

	// Базова валідація
	if fullName == "" {
		http.Error(w, "Ім'я клієнта обов'язкове", http.StatusBadRequest)
		return
	}

	// Створюємо нового клієнта
	client := models.Client{
		UserID:    userId,
		FullName:  fullName,
		Email:     email,
		Phone:     phone,
		Company:   company,
		Category:  category,
		Source:    source,
		Address:   address,
		Notes:     notes,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Зберігаємо клієнта в базу даних
	_, err = h.DB.NewInsert().Model(&client).Exec(r.Context())
	if err != nil {
		h.Logger.Error("Error creating client", err)
		http.Error(w, "Помилка збереження клієнта", http.StatusInternalServerError)
		return
	}

	// Обробка завантаження файлу аватара
	file, header, err := r.FormFile("avatar")
	if err == nil && header != nil {
		defer file.Close()

		// Визначаємо шлях до файлу
		avatarPath := fmt.Sprintf("static/img/clients/%d%s", client.ID, utils.GetFileExtension(header.Filename))

		// Зберігаємо файл
		err = utils.SaveUploadedFile(file, avatarPath)
		if err != nil {
			h.Logger.Error("Error saving avatar", err)
			// Не повертаємо помилку, просто логуємо
		} else {
			// Оновлюємо шлях до аватара в базі даних
			client.Avatar = "/" + avatarPath
			_, err = h.DB.NewUpdate().Model(&client).Column("avatar").WherePK().Exec(r.Context())
			if err != nil {
				h.Logger.Error("Error updating avatar path", err)
			}
		}
	}

	// Логуємо активність
	activity := models.Activity{
		UserID:      userId,
		EntityType:  "client",
		EntityID:    client.ID,
		Action:      "create",
		Description: fmt.Sprintf("Створено клієнта: %s", client.FullName),
		CreatedAt:   time.Now(),
	}

	_, err = h.DB.NewInsert().Model(&activity).Exec(r.Context())
	if err != nil {
		h.Logger.Error("Error logging activity", err)
	}

	// Повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Клієнта успішно створено",
		"client":  client,
	})
}

// Get отримує дані клієнта та його бронювання
func (h *ClientHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	vars := mux.Vars(r)
	clientId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некоректний ID клієнта", http.StatusBadRequest)
		return
	}

	// Отримуємо дані клієнта
	client := models.Client{}
	err = h.DB.NewSelect().
		Model(&client).
		Where("id = ? AND user_id = ?", clientId, userId).
		Scan(r.Context())

	if err != nil {
		h.Logger.Error("Error fetching client", err, "client_id", clientId)
		http.Error(w, "Клієнта не знайдено", http.StatusNotFound)
		return
	}

	// Отримуємо бронювання клієнта
	var bookings []models.Booking
	err = h.DB.NewSelect().
		Model(&bookings).
		Where("client_id = ?", clientId).
		Order("start_time DESC").
		Limit(5).
		Scan(r.Context())

	if err != nil {
		h.Logger.Error("Error fetching bookings", err, "client_id", clientId)
		bookings = []models.Booking{}
	}

	// Структура для відображення бронювань
	type BookingViewModel struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		EventType   string    `json:"eventType"`
		Status      string    `json:"status"`
		StartTime   time.Time `json:"startTime"`
		EndTime     time.Time `json:"endTime"`
		DisplayDate string    `json:"displayDate"`
		DisplayTime string    `json:"displayTime"`
	}

	// Підготовка даних бронювань для відображення
	bookingViewModels := make([]BookingViewModel, 0, len(bookings))
	for _, booking := range bookings {
		bookingVM := BookingViewModel{
			ID:          booking.ID.String(),
			Title:       booking.Title,
			EventType:   booking.EventType,
			Status:      booking.Status,
			StartTime:   booking.StartTime,
			EndTime:     booking.EndTime,
			DisplayDate: booking.StartTime.Format("02.01.2006"),
			DisplayTime: booking.StartTime.Format("15:04") + " - " + booking.EndTime.Format("15:04"),
		}
		bookingViewModels = append(bookingViewModels, bookingVM)
	}

	// Отримуємо активності, пов'язані з клієнтом
	var activities []models.Activity
	err = h.DB.NewSelect().
		Model(&activities).
		Where("entity_type = 'client' AND entity_id = ?", clientId).
		Order("created_at DESC").
		Limit(10).
		Scan(r.Context())

	if err != nil {
		h.Logger.Error("Error fetching activities", err, "client_id", clientId)
		activities = []models.Activity{}
	}

	// Структура для відображення активностей
	type ActivityViewModel struct {
		ID            int64     `json:"id"`
		Action        string    `json:"action"`
		Description   string    `json:"description"`
		CreatedAt     time.Time `json:"createdAt"`
		FormattedDate string    `json:"formattedDate"`
		FormattedTime string    `json:"formattedTime"`
	}

	// Підготовка даних активностей для відображення
	activityViewModels := make([]ActivityViewModel, 0, len(activities))
	for _, activity := range activities {
		activityVM := ActivityViewModel{
			ID:            activity.ID,
			Action:        activity.Action,
			Description:   activity.Description,
			CreatedAt:     activity.CreatedAt,
			FormattedDate: activity.CreatedAt.Format("02.01.2006"),
			FormattedTime: activity.CreatedAt.Format("15:04"),
		}
		activityViewModels = append(activityViewModels, activityVM)
	}

	// Передаємо дані для JSON відповіді
	data := map[string]interface{}{
		"Client":     client,
		"Bookings":   bookingViewModels,
		"Activities": activityViewModels,
		"CSRFToken":  "dummy_token", // В реальному додатку тут має бути справжній CSRF токен
	}

	// Перетворюємо дані в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.Logger.Error("Error marshaling JSON", err)
		http.Error(w, "Помилка підготовки даних", http.StatusInternalServerError)
		return
	}

	// Встановлюємо заголовки та повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// EditPage відображає форму для редагування клієнта
func (h *ClientHandler) EditPage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	vars := mux.Vars(r)
	clientId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некоректний ID клієнта", http.StatusBadRequest)
		return
	}

	// Отримуємо дані клієнта
	client := models.Client{}
	err = h.DB.NewSelect().
		Model(&client).
		Where("id = ? AND user_id = ?", clientId, userId).
		Scan(r.Context())

	if err != nil {
		h.Logger.Error("Error fetching client", err, "client_id", clientId)
		http.Error(w, "Клієнта не знайдено", http.StatusNotFound)
		return
	}

	// Підготовка даних для форми
	data := map[string]interface{}{
		"Title":      "Редагування клієнта: " + client.FullName,
		"Action":     "edit",
		"Client":     client,
		"Categories": models.DefaultCategories(),
		"Sources":    models.DefaultSources(),
		"CSRFToken":  "dummy_token", // В реальному додатку тут має бути справжній CSRF токен
	}

	// Перетворюємо дані в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.Logger.Error("Error marshaling JSON", err)
		http.Error(w, "Помилка підготовки даних", http.StatusInternalServerError)
		return
	}

	// Встановлюємо заголовки та повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Update обробляє запит на оновлення інформації про клієнта
func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	vars := mux.Vars(r)
	clientId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некоректний ID клієнта", http.StatusBadRequest)
		return
	}

	// Перевіряємо, чи існує клієнт і чи належить він цьому користувачу
	var existingClient models.Client
	err = h.DB.NewSelect().
		Model(&existingClient).
		Where("id = ? AND user_id = ?", clientId, userId).
		Scan(r.Context())

	if err != nil {
		h.Logger.Error("Error fetching client for update", err, "client_id", clientId)
		http.Error(w, "Клієнта не знайдено", http.StatusNotFound)
		return
	}

	// Парсимо форму
	err = r.ParseMultipartForm(10 << 20) // 10 MB максимальний розмір файлу
	if err != nil {
		h.Logger.Error("Error parsing form", err)
		http.Error(w, "Помилка обробки форми", http.StatusBadRequest)
		return
	}

	// Отримуємо дані з форми
	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	company := r.FormValue("company")
	category := r.FormValue("category")
	source := r.FormValue("source")
	address := r.FormValue("address")
	notes := r.FormValue("notes")
	isActive := r.FormValue("is_active") == "true" || r.FormValue("is_active") == "on"

	// Базова валідація
	if fullName == "" {
		http.Error(w, "Ім'я клієнта обов'язкове", http.StatusBadRequest)
		return
	}

	// Оновлюємо дані клієнта
	client := models.Client{
		ID:        clientId,
		UserID:    userId,
		FullName:  fullName,
		Email:     email,
		Phone:     phone,
		Company:   company,
		Category:  category,
		Source:    source,
		Address:   address,
		Notes:     notes,
		IsActive:  isActive,
		Avatar:    existingClient.Avatar, // Зберігаємо існуючий аватар
		UpdatedAt: time.Now(),
	}

	// Обробка завантаження файлу аватара
	file, header, err := r.FormFile("avatar")
	if err == nil && header != nil {
		defer file.Close()

		// Визначаємо шлях до файлу
		avatarPath := fmt.Sprintf("static/img/clients/%d%s", client.ID, utils.GetFileExtension(header.Filename))

		// Зберігаємо файл
		err = utils.SaveUploadedFile(file, avatarPath)
		if err != nil {
			h.Logger.Error("Error saving avatar", err)
			// Не повертаємо помилку, просто логуємо
		} else {
			// Оновлюємо шлях до аватара
			client.Avatar = "/" + avatarPath
		}
	}

	// Оновлюємо клієнта в базі даних
	_, err = h.DB.NewUpdate().
		Model(&client).
		WherePK().
		Exec(r.Context())

	if err != nil {
		h.Logger.Error("Error updating client", err)
		http.Error(w, "Помилка оновлення клієнта", http.StatusInternalServerError)
		return
	}

	// Логуємо активність
	activity := models.Activity{
		UserID:      userId,
		EntityType:  "client",
		EntityID:    client.ID,
		Action:      "update",
		Description: fmt.Sprintf("Оновлено клієнта: %s", client.FullName),
		CreatedAt:   time.Now(),
	}

	_, err = h.DB.NewInsert().Model(&activity).Exec(r.Context())
	if err != nil {
		h.Logger.Error("Error logging activity", err)
	}

	// Повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Клієнта успішно оновлено",
		"client":  client,
	})
}

// Delete обробляє запит на видалення клієнта
func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	vars := mux.Vars(r)
	clientId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некоректний ID клієнта", http.StatusBadRequest)
		return
	}

	// Перевіряємо, чи існує клієнт і чи належить він цьому користувачу
	var client models.Client
	err = h.DB.NewSelect().
		Model(&client).
		Where("id = ? AND user_id = ?", clientId, userId).
		Scan(r.Context())

	if err != nil {
		h.Logger.Error("Error fetching client for delete", err, "client_id", clientId)
		http.Error(w, "Клієнта не знайдено", http.StatusNotFound)
		return
	}

	// Перевіряємо, чи можна видалити клієнта
	// Наприклад, перевіряємо чи є активні бронювання
	var bookingsCount int
	_, err = h.DB.NewSelect().
		Model((*models.Booking)(nil)).
		Where("client_id = ? AND status NOT IN (?, ?)", clientId, "canceled", "completed").
		Count(r.Context())

	if err != nil {
		h.Logger.Error("Error counting bookings", err)
	} else if bookingsCount > 0 {
		http.Error(w, "Неможливо видалити клієнта з активними бронюваннями", http.StatusBadRequest)
		return
	}

	// Видаляємо всі активності, пов'язані з клієнтом
	_, err = h.DB.NewDelete().
		Model((*models.Activity)(nil)).
		Where("entity_type = 'client' AND entity_id = ?", clientId).
		Exec(r.Context())
	if err != nil {
		h.Logger.Error("Error deleting client activities", err)
		// Не повертаємо помилку, просто логуємо
	}

	// Видаляємо клієнта
	_, err = h.DB.NewDelete().
		Model(&client).
		WherePK().
		Exec(r.Context())

	if err != nil {
		h.Logger.Error("Error deleting client", err)
		http.Error(w, "Помилка видалення клієнта", http.StatusInternalServerError)
		return
	}

	// Логуємо активність
	activity := models.Activity{
		UserID:      userId,
		EntityType:  "user",
		EntityID:    userId,
		Action:      "delete",
		Description: fmt.Sprintf("Видалено клієнта: %s", client.FullName),
		CreatedAt:   time.Now(),
	}

	_, err = h.DB.NewInsert().Model(&activity).Exec(r.Context())
	if err != nil {
		h.Logger.Error("Error logging activity", err)
	}

	// Повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Клієнта успішно видалено",
	})
}

// Import обробляє імпорт клієнтів з файлів CSV або Excel
func (h *ClientHandler) Import(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)

	// Парсимо форму
	err := r.ParseMultipartForm(20 << 20) // 20 MB максимальний розмір файлу
	if err != nil {
		h.Logger.Error("Error parsing form", err)
		http.Error(w, "Помилка обробки форми", http.StatusBadRequest)
		return
	}

	// Отримуємо файл
	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Error("Error getting uploaded file", err)
		http.Error(w, "Помилка отримання файлу", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Перевіряємо тип файлу
	fileExt := utils.GetFileExtension(header.Filename)
	if fileExt != ".csv" && fileExt != ".xlsx" && fileExt != ".xls" {
		http.Error(w, "Підтримуються лише файли CSV та Excel (.xlsx, .xls)", http.StatusBadRequest)
		return
	}

	// Створюємо тимчасовий файл для обробки
	tempFilePath := fmt.Sprintf("temp/%s%s", uuid.New().String(), fileExt)

	// Зберігаємо файл на диск
	err = utils.SaveUploadedFile(file, tempFilePath)
	if err != nil {
		h.Logger.Error("Error saving uploaded file", err)
		http.Error(w, "Помилка збереження файлу", http.StatusInternalServerError)
		return
	}

	// Обробляємо файл залежно від його типу
	fileExt = utils.GetFileExtension(header.Filename)
	var clients []models.Client
	var importErr error

	switch fileExt {
	case ".csv":
		// Реалізуйте цей метод у вашому utils пакеті
		clients, importErr = parseCSVClients(tempFilePath, userId)
	case ".xlsx", ".xls":
		// Реалізуйте цей метод у вашому utils пакеті
		clients, importErr = parseExcelClients(tempFilePath, userId)
	}

	// Видаляємо тимчасовий файл
	os.Remove(tempFilePath)

	if importErr != nil {
		h.Logger.Error("Error importing clients", importErr)
		http.Error(w, "Помилка імпорту клієнтів: "+importErr.Error(), http.StatusInternalServerError)
		return
	}

	// Імпортуємо клієнтів
	importedCount := 0
	duplicateCount := 0

	for _, client := range clients {
		// Перевіряємо, чи існує клієнт з таким email
		var existingClient models.Client
		err = h.DB.NewSelect().
			Model(&existingClient).
			Where("user_id = ? AND email = ? AND email != ''", userId, client.Email).
			Scan(r.Context())

		if err == nil {
			// Клієнт з таким email вже існує, оновлюємо його
			existingClient.FullName = client.FullName
			existingClient.Phone = client.Phone
			existingClient.Company = client.Company
			existingClient.Category = client.Category
			existingClient.Source = client.Source
			existingClient.Address = client.Address
			existingClient.Notes = client.Notes
			existingClient.UpdatedAt = time.Now()

			_, err = h.DB.NewUpdate().
				Model(&existingClient).
				WherePK().
				Exec(r.Context())

			if err != nil {
				h.Logger.Error("Error updating client during import", err)
			} else {
				duplicateCount++
			}
		} else {
			// Створюємо нового клієнта
			client.UserID = userId
			client.CreatedAt = time.Now()
			client.UpdatedAt = time.Now()

			_, err = h.DB.NewInsert().
				Model(&client).
				Exec(r.Context())

			if err != nil {
				h.Logger.Error("Error creating client during import", err)
			} else {
				importedCount++
			}
		}
	}

	// Логуємо активність
	activity := models.Activity{
		UserID:      userId,
		EntityType:  "user",
		EntityID:    userId,
		Action:      "import_clients",
		Description: fmt.Sprintf("Імпортовано %d нових клієнтів, оновлено %d існуючих", importedCount, duplicateCount),
		CreatedAt:   time.Now(),
	}

	_, err = h.DB.NewInsert().Model(&activity).Exec(r.Context())
	if err != nil {
		h.Logger.Error("Error logging activity", err)
	}

	// Повертаємо відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("Імпорт завершено. Додано %d нових клієнтів, оновлено %d існуючих.", importedCount, duplicateCount),
		"importedCount": importedCount,
		"updatedCount":  duplicateCount,
	})
}

// parseCSVClients розбирає клієнтів з CSV файлу
func parseCSVClients(filePath string, userId int64) ([]models.Client, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("помилка відкриття файлу: %w", err)
	}
	defer file.Close()

	// Читаємо CSV
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Читаємо заголовки
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("помилка читання заголовків: %w", err)
	}

	// Маппінг колонок
	columnMap := make(map[string]int)
	for i, header := range headers {
		columnMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	// Перевіряємо наявність обов'язкових колонок
	_, hasName := columnMap["ім'я"]
	_, hasImya := columnMap["імя"]
	_, hasEnglishName := columnMap["name"]

	if !hasName && !hasImya && !hasEnglishName {
		return nil, fmt.Errorf("відсутня обов'язкова колонка з ім'ям клієнта")
	}

	var clients []models.Client

	// Читаємо рядки
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("помилка читання рядка: %w", err)
		}

		// Створюємо клієнта з рядка
		client := models.Client{
			UserID:    userId,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Заповнюємо поля
		for key, idx := range columnMap {
			if idx >= len(row) {
				continue
			}
			value := strings.TrimSpace(row[idx])
			if value == "" {
				continue
			}

			switch key {
			case "ім'я", "імя", "name", "full_name", "fullname":
				client.FullName = value
			case "email", "електронна пошта", "пошта":
				client.Email = value
			case "телефон", "phone":
				client.Phone = value
			case "компанія", "company":
				client.Company = value
			case "категорія", "category":
				client.Category = value
			case "джерело", "source":
				client.Source = value
			case "адреса", "address":
				client.Address = value
			case "примітки", "notes":
				client.Notes = value
			}
		}

		// Перевіряємо наявність імені
		if client.FullName == "" {
			continue
		}

		clients = append(clients, client)
	}

	return clients, nil
}

// parseExcelClients розбирає клієнтів з Excel файлу
// Цей метод потрібно реалізувати з використанням бібліотеки для роботи з Excel
func parseExcelClients(filePath string, userId int64) ([]models.Client, error) {
	// В майбутньому тут буде код для читання Excel файлу
	// та конвертації даних в масив клієнтів
	// За допомогою userId будуть створюватись клієнти
	// з правильним ідентифікатором користувача

	// Наразі просто повертаємо заглушку з помилкою
	return nil, fmt.Errorf("імпорт з Excel файлів ще не реалізовано")
}
