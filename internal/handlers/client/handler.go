package client

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/services/booking"
	"timebride/internal/services/client"

	"gorm.io/datatypes"
)

// Handler обробляє запити для роботи з клієнтами
type Handler struct {
	clientService  client.IClientService
	bookingService booking.IBookingService
}

// NewHandler створює новий обробник клієнтів
func NewHandler(clientService client.IClientService, bookingService booking.IBookingService) *Handler {
	return &Handler{
		clientService:  clientService,
		bookingService: bookingService,
	}
}

// List returns a list of clients
func (h *Handler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 10

	// Get filter parameters
	query := c.Query("q", "")
	category := c.Query("category", "")
	source := c.Query("source", "")

	// Get clients with pagination
	options := models.ClientListOptions{
		Page:     page,
		PageSize: pageSize,
		Search:   query,
		Category: category,
		Source:   source,
		SortBy:   "created_at",
		SortDesc: true,
	}

	result, err := h.clientService.ListClients(c.Context(), userUUID, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch clients",
		})
	}

	return c.JSON(fiber.Map{
		"items":       result.Items,
		"total":       result.TotalItems,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": (result.TotalItems + int64(pageSize) - 1) / int64(pageSize),
	})
}

// Create creates a new client
func (h *Handler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var input struct {
		FullName string                `json:"full_name"`
		Phone    string                `json:"phone"`
		Email    string                `json:"email"`
		Notes    string                `json:"notes"`
		Settings models.ClientSettings `json:"settings"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	settingsJSON, err := json.Marshal(input.Settings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process settings",
		})
	}

	client := &models.Client{
		ID:       uuid.New(),
		UserID:   userUUID,
		FullName: input.FullName,
		Phone:    input.Phone,
		Email:    input.Email,
		Notes:    input.Notes,
		Settings: datatypes.JSON(settingsJSON),
	}

	if err := client.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createdClient, err := h.clientService.CreateClient(c.Context(), client)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create client",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdClient.ToPublic())
}

// Get returns a client by ID
func (h *Handler) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	clientID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	client, err := h.clientService.GetClient(c.Context(), userUUID, clientID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Client not found",
		})
	}

	// Get client's bookings
	bookings, err := h.bookingService.GetByClientID(c.Context(), clientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch client's bookings",
		})
	}

	return c.JSON(fiber.Map{
		"client":   client.ToPublic(),
		"bookings": bookings,
	})
}

// Update updates a client
func (h *Handler) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	clientID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	var input struct {
		FullName string                `json:"full_name"`
		Phone    string                `json:"phone"`
		Email    string                `json:"email"`
		Notes    string                `json:"notes"`
		Settings models.ClientSettings `json:"settings"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	settingsJSON, err := json.Marshal(input.Settings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process settings",
		})
	}

	client := &models.Client{
		ID:       clientID,
		UserID:   userUUID,
		FullName: input.FullName,
		Phone:    input.Phone,
		Email:    input.Email,
		Notes:    input.Notes,
		Settings: datatypes.JSON(settingsJSON),
	}

	if err := client.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	updatedClient, err := h.clientService.UpdateClient(c.Context(), client)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update client",
		})
	}

	return c.JSON(updatedClient.ToPublic())
}

// Delete deletes a client
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	clientID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	err = h.clientService.DeleteClient(c.Context(), userUUID, clientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete client",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CreateClient створює нового клієнта
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input models.CreateClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := json.Marshal(input.Settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := models.Client{
		ID:       uuid.New(),
		UserID:   input.UserID,
		FullName: input.FullName,
		Email:    input.Email,
		Phone:    input.Phone,
		Notes:    input.Notes,
		Settings: datatypes.JSON(settings),
	}

	if err := h.clientService.Create(r.Context(), &client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

// UpdateClient оновлює існуючого клієнта
func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var input models.UpdateClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := json.Marshal(input.Settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := models.Client{
		ID:       input.ID,
		UserID:   input.UserID,
		FullName: input.FullName,
		Email:    input.Email,
		Phone:    input.Phone,
		Notes:    input.Notes,
		Settings: datatypes.JSON(settings),
	}

	if err := h.clientService.Update(r.Context(), &client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

// GetClientBookings повертає бронювання клієнта
func (h *Handler) GetClientBookings(w http.ResponseWriter, r *http.Request) {
	clientID, err := uuid.Parse(r.URL.Query().Get("client_id"))
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	bookings, err := h.bookingService.GetByClient(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bookings)
}
