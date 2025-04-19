package handlers

import (
	"context"
	"log"
	"time"

	"timebride/internal/models"
	"timebride/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	userService     services.UserService
	bookingService  services.BookingService
	templateService services.TemplateService
	fileService     services.StorageService
}

func NewDashboardHandler(
	userService services.UserService,
	bookingService services.BookingService,
	templateService services.TemplateService,
	fileService services.StorageService,
) *DashboardHandler {
	return &DashboardHandler{
		userService:     userService,
		bookingService:  bookingService,
		templateService: templateService,
		fileService:     fileService,
	}
}

func (h *DashboardHandler) HandleDashboard(c *fiber.Ctx) error {
	log.Println("Dashboard handler started")

	// Отримуємо ID користувача з контексту
	userIDRaw := c.Locals("user_id")
	log.Printf("Raw user_id from context: %v (type: %T)", userIDRaw, userIDRaw)

	var userID uuid.UUID
	var err error

	// Перевіряємо тип userIDRaw і конвертуємо в uuid.UUID
	switch id := userIDRaw.(type) {
	case string:
		log.Printf("user_id is string: %s", id)
		userID, err = uuid.Parse(id)
		if err != nil {
			log.Printf("Error parsing user_id as UUID: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		log.Printf("Successfully parsed user_id as UUID: %s", userID.String())
	case uuid.UUID:
		log.Printf("user_id is already UUID: %s", id.String())
		userID = id
	default:
		log.Printf("Unknown user_id type: %T, value: %v", userIDRaw, userIDRaw)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID type",
		})
	}

	// Отримуємо дані користувача
	log.Printf("Fetching user data for ID: %s", userID.String())
	user, err := h.userService.GetByID(context.Background(), userID)
	if err != nil {
		log.Printf("Error fetching user data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user data",
		})
	}
	log.Printf("User found: %s", user.Email)

	// Отримуємо бронювання та іншу статистику
	log.Println("Fetching user bookings and statistics")
	bookings, err := h.bookingService.List(context.Background(), userID)
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		// Продовжуємо виконання, навіть якщо бронювання не знайдені
		bookings = []*models.Booking{}
	}

	// Розрахунок статистики
	totalBookings := len(bookings)
	log.Printf("Total bookings: %d", totalBookings)

	activeBookings := 0
	upcomingEvents := 0

	now := time.Now()
	for _, b := range bookings {
		if b.Status == "active" {
			activeBookings++
		}

		if b.StartTime.After(now) {
			upcomingEvents++
		}
	}

	log.Printf("Active bookings: %d, Upcoming events: %d", activeBookings, upcomingEvents)

	// Отримуємо кількість шаблонів
	templates, err := h.templateService.GetByUserID(context.Background(), userID)
	totalTemplates := 0
	if err == nil {
		totalTemplates = len(templates)
	}
	log.Printf("Total templates: %d", totalTemplates)

	// Створюємо модель для відображення
	viewModel := fiber.Map{
		"Title": "Головна панель",
		"User": fiber.Map{
			"FullName": user.FullName,
			"Email":    user.Email,
			"Role":     user.Role,
			"Initials": getInitials(user.FullName),
		},
		"Bookings":       bookings,
		"RecentBookings": formatBookings(bookings),
		"Stats": fiber.Map{
			"TotalBookings":   totalBookings,
			"ActiveBookings":  activeBookings,
			"UpcomingEvents":  upcomingEvents,
			"EventsThisMonth": 0, // Поки що нульове значення
			"TotalTemplates":  totalTemplates,
			"ActiveTemplates": 0,      // Поки що нульове значення
			"TotalFiles":      0,      // Поки що нульове значення
			"TotalSize":       "0 MB", // Поки що нульове значення
		},
	}

	log.Println("Dashboard handler completed successfully")
	return c.Render("simple_dashboard", viewModel)
}

// getStatusClass returns CSS class for booking status
func getStatusClass(status string) string {
	switch status {
	case "pending":
		return "bg-yellow-100 text-yellow-800"
	case "confirmed":
		return "bg-green-100 text-green-800"
	case "cancelled":
		return "bg-red-100 text-red-800"
	case "completed":
		return "bg-blue-100 text-blue-800"
	default:
		return "bg-gray-100 text-gray-800"
	}
}

// getInitials returns initials from full name
func getInitials(fullName string) string {
	if len(fullName) == 0 {
		return ""
	}

	initials := string(fullName[0])
	for i := 0; i < len(fullName)-1; i++ {
		if fullName[i] == ' ' {
			initials += string(fullName[i+1])
		}
	}

	return initials
}

// форматує бронювання для відображення
func formatBookings(bookings []*models.Booking) []fiber.Map {
	result := make([]fiber.Map, 0, len(bookings))
	for _, b := range bookings {
		result = append(result, fiber.Map{
			"ID":          b.ID,
			"ClientName":  b.ClientName,
			"EventType":   b.EventType,
			"StartTime":   b.StartTime.Format("02.01.2006 15:04"),
			"Status":      b.Status,
			"StatusClass": getStatusClass(b.Status),
		})
	}
	return result
}
