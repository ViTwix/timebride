package handlers

import (
	"encoding/json"
	"net/http"

	"timebride/internal/services/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"
)

type UserHandler struct {
	userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// List повертає список користувачів
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Create створює нового користувача
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		FullName    string `json:"full_name"`
		CompanyName string `json:"company_name"`
		Role        string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), input.Email, input.Password, input.FullName, input.CompanyName, input.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Get повертає користувача за ID
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Користувача не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Update оновлює дані користувача
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Користувача не знайдено", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	if err := h.userService.Update(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Delete видаляє користувача
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ChangePassword змінює пароль користувача
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.userService.UpdatePassword(r.Context(), user, input.CurrentPassword, input.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateSettings оновлює налаштування користувача
func (h *UserHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Конвертуємо налаштування в JSON
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		http.Error(w, "Помилка конвертації налаштувань", http.StatusInternalServerError)
		return
	}

	user.Settings = datatypes.JSON(settingsJSON)

	if err := h.userService.Update(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleProfile відображає профіль користувача
func (h *UserHandler) HandleProfile(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Отримуємо користувача
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Рендеримо шаблон
	return c.Render("user/profile", fiber.Map{
		"Title": "Мій профіль",
		"User":  user,
	})
}

// HandleSettings відображає налаштування користувача
func (h *UserHandler) HandleSettings(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Отримуємо користувача
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Рендеримо шаблон
	return c.Render("user/settings", fiber.Map{
		"Title": "Налаштування",
		"User":  user,
	})
}

func (h *UserHandler) HandleUpdateProfile(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Отримуємо користувача
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Парсимо дані з форми
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Оновлюємо користувача
	if err := h.userService.Update(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Redirect("/profile")
}

func (h *UserHandler) HandleUpdatePassword(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Отримуємо користувача
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Парсимо дані з форми
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Оновлюємо пароль
	if err := h.userService.UpdatePassword(c.Context(), user, input.CurrentPassword, input.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.Redirect("/settings")
}
