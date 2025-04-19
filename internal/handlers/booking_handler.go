package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"

	"timebride/internal/models"
	"timebride/internal/services/booking"
	"timebride/internal/services/user"
)

type BookingHandler struct {
	bookingService *booking.Service
	userService    *user.Service
}

func NewBookingHandler(bookingService *booking.Service, userService *user.Service) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		userService:    userService,
	}
}

// List повертає список бронювань
func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту
	userID := r.Context().Value("user_id").(uuid.UUID)

	bookings, err := h.bookingService.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// Create створює нове бронювання
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title        string                 `json:"title"`
		EventType    string                 `json:"event_type"`
		StartTime    time.Time              `json:"start_time"`
		EndTime      time.Time              `json:"end_time"`
		ClientName   string                 `json:"client_name"`
		ClientPhone  string                 `json:"client_phone"`
		ClientEmail  string                 `json:"client_email"`
		CustomFields map[string]interface{} `json:"custom_fields"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	// Конвертуємо CustomFields в JSON
	customFieldsJSON, err := json.Marshal(input.CustomFields)
	if err != nil {
		http.Error(w, "Помилка конвертації кастомних полів", http.StatusBadRequest)
		return
	}

	// Створюємо нове бронювання
	booking := &models.Booking{
		UserID:        userID,
		Title:         input.Title,
		EventType:     input.EventType,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		ClientName:    input.ClientName,
		ClientPhone:   input.ClientPhone,
		ClientEmail:   input.ClientEmail,
		Status:        models.BookingStatusPending,
		PaymentStatus: models.PaymentStatusPending,
		CustomFields:  datatypes.JSON(customFieldsJSON),
	}

	// Викликаємо сервіс для створення бронювання
	if err := h.bookingService.Create(r.Context(), booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

// Get повертає бронювання за ID
func (h *BookingHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Бронювання не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// GetByDateRange повертає бронювання в заданому діапазоні дат
func (h *BookingHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	bookings, err := h.bookingService.GetByDateRange(r.Context(), userID, input.StartTime, input.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// Update оновлює бронювання
func (h *BookingHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Бронювання не знайдено", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(booking); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	if err := h.bookingService.Update(r.Context(), booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// Delete видаляє бронювання
func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	if err := h.bookingService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateStatus оновлює статус бронювання
func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	if err := h.bookingService.UpdateStatus(r.Context(), id, input.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdatePaymentStatus оновлює статус оплати бронювання
func (h *BookingHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	var input struct {
		PaymentStatus string `json:"payment_status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	if err := h.bookingService.UpdatePaymentStatus(r.Context(), id, input.PaymentStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddTeamMember додає учасника команди до бронювання
func (h *BookingHandler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	var input struct {
		MemberID uuid.UUID `json:"member_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	teamMember := map[string]interface{}{
		"member_id": input.MemberID,
	}

	if err := h.bookingService.AddTeamMember(r.Context(), id, teamMember); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveTeamMember видаляє учасника команди з бронювання
func (h *BookingHandler) RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	memberID := vars["member-id"]

	if err := h.bookingService.RemoveTeamMember(r.Context(), id, memberID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleBookingList відображає список бронювань
func (h *BookingHandler) HandleBookingList(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Отримуємо бронювання користувача
	bookings, err := h.bookingService.List(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get bookings",
		})
	}

	// Рендеримо шаблон
	return c.Render("booking/list", fiber.Map{
		"Title":    "Мої бронювання",
		"Bookings": bookings,
	})
}

func (h *BookingHandler) HandleBookingCreate(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Парсимо дані з форми
	var booking models.Booking
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Встановлюємо ID користувача
	booking.UserID = userID

	// Створюємо бронювання
	if err := h.bookingService.Create(c.Context(), &booking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	return c.Redirect("/bookings")
}

func (h *BookingHandler) HandleBookingUpdate(c *fiber.Ctx) error {
	// Отримуємо ID бронювання
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	// Отримуємо бронювання
	booking, err := h.bookingService.GetByID(c.Context(), bookingID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	// Парсимо дані з форми
	if err := c.BodyParser(booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Оновлюємо бронювання
	if err := h.bookingService.Update(c.Context(), booking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update booking",
		})
	}

	return c.Redirect("/bookings")
}

func (h *BookingHandler) HandleBookingDelete(c *fiber.Ctx) error {
	// Отримуємо ID бронювання
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	// Видаляємо бронювання
	if err := h.bookingService.Delete(c.Context(), bookingID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete booking",
		})
	}

	return c.Redirect("/bookings")
}
