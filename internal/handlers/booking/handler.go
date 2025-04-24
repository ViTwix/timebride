package booking

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/services/booking"
	"timebride/internal/services/user"
	"timebride/internal/types"
)

// Handler реалізує обробку запитів бронювань
type Handler struct {
	bookingService booking.IBookingService
	userService    user.IUserService
}

// NewHandler створює новий екземпляр обробника бронювань
func NewHandler(bookingService booking.IBookingService, userService user.IUserService) *Handler {
	return &Handler{
		bookingService: bookingService,
		userService:    userService,
	}
}

// List отримує список бронювань
func (h *Handler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.ErrUnauthorized
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	offset := (page - 1) * pageSize

	filter := map[string]interface{}{
		"user_id": userID,
		"limit":   pageSize,
		"offset":  offset,
	}

	bookings, err := h.bookingService.List(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"items":     bookings,
		"page":      page,
		"page_size": pageSize,
		"total":     len(bookings),
	})
}

// Create створює нове бронювання
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var input types.BookingCreate
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	input.UserID = userID.String()

	// Конвертуємо types.BookingCreate в models.BookingCreate
	modelInput := &models.BookingCreate{
		UserID:       input.UserID,
		ClientID:     input.ClientID,
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		Amount:       input.Amount,
		Prepayment:   input.Prepayment,
		Currency:     input.Currency,
		Description:  input.Description,
		Location:     input.Location,
		TeamMembers:  input.TeamMembers,
		CustomFields: input.CustomFields,
	}

	booking, err := h.bookingService.Create(c.Context(), modelInput)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(booking)
}

// Get отримує бронювання за ID
func (h *Handler) Get(c *fiber.Ctx) error {
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	booking, err := h.bookingService.Get(c.Context(), bookingID)
	if err != nil {
		return err
	}

	return c.JSON(booking)
}

// Update оновлює бронювання
func (h *Handler) Update(c *fiber.Ctx) error {
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	var input types.BookingUpdate
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	// Конвертуємо types.BookingUpdate в models.BookingUpdate
	modelInput := &models.BookingUpdate{
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		Status:       (*models.BookingStatus)(input.Status),
		Amount:       input.Amount,
		Prepayment:   input.Prepayment,
		Currency:     input.Currency,
		Description:  input.Description,
		Location:     input.Location,
		TeamMembers:  input.TeamMembers,
		CustomFields: input.CustomFields,
	}

	booking, err := h.bookingService.Update(c.Context(), bookingID, modelInput)
	if err != nil {
		return err
	}

	return c.JSON(booking)
}

// Delete видаляє бронювання
func (h *Handler) Delete(c *fiber.Ctx) error {
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.bookingService.Delete(c.Context(), bookingID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetCalendarEvents отримує події для календаря
func (h *Handler) GetCalendarEvents(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.ErrUnauthorized
	}

	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	bookings, err := h.bookingService.GetByDateRange(c.Context(), userID, startDate, endDate)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// GetStatistics отримує статистику бронювань
func (h *Handler) GetStatistics(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.ErrUnauthorized
	}

	period := c.Query("period", "month")
	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "month":
		startDate = now.AddDate(0, 0, -30)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, 0, -30)
	}
	endDate = now

	// Отримуємо статистику
	bookings, err := h.bookingService.GetByDateRange(c.Context(), userID, startDate, endDate)
	if err != nil {
		return err
	}

	upcoming, err := h.bookingService.CountUpcoming(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"total_bookings": len(bookings),
		"upcoming":       upcoming,
		"period":         period,
		"start_date":     startDate,
		"end_date":       endDate,
	})
}
