package handlers

import (
	"encoding/json"
	"net/http"

	"timebride/internal/models"
	"timebride/internal/services/template"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"
)

type TemplateHandler struct {
	templateService *template.Service
}

func NewTemplateHandler(templateService *template.Service) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

// List повертає список шаблонів
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	templates, err := h.templateService.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// Get повертає шаблон за ID
func (h *TemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	template, err := h.templateService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Шаблон не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// Create створює новий шаблон
func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string                 `json:"name"`
		EventType      string                 `json:"event_type"`
		FieldsTemplate map[string]interface{} `json:"fields_template"`
		TeamTemplate   map[string]interface{} `json:"team_template"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	// Конвертуємо шаблони в JSON
	fieldsJSON, err := json.Marshal(input.FieldsTemplate)
	if err != nil {
		http.Error(w, "Помилка конвертації шаблону полів", http.StatusInternalServerError)
		return
	}

	teamJSON, err := json.Marshal(input.TeamTemplate)
	if err != nil {
		http.Error(w, "Помилка конвертації шаблону команди", http.StatusInternalServerError)
		return
	}

	template := &models.Template{
		UserID:         userID,
		Name:           input.Name,
		EventType:      input.EventType,
		FieldsTemplate: datatypes.JSON(fieldsJSON),
		TeamTemplate:   datatypes.JSON(teamJSON),
	}

	if err := h.templateService.Create(r.Context(), template); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetByID повертає шаблон за ID
func (h *TemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	template, err := h.templateService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Шаблон не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// GetByEventType повертає шаблони за типом події
func (h *TemplateHandler) GetByEventType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := vars["event_type"]

	// Отримуємо ID користувача з контексту (після авторизації)
	userID := r.Context().Value("user_id").(uuid.UUID)

	templates, err := h.templateService.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Фільтруємо шаблони за типом події
	var filteredTemplates []*models.Template
	for _, template := range templates {
		if template.EventType == eventType {
			filteredTemplates = append(filteredTemplates, template)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredTemplates)
}

// Update оновлює шаблон
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	template, err := h.templateService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Шаблон не знайдено", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(template); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	if err := h.templateService.Update(r.Context(), template); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// Delete видаляє шаблон
func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Невірний формат ID", http.StatusBadRequest)
		return
	}

	if err := h.templateService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleTemplateList відображає список шаблонів
func (h *TemplateHandler) HandleTemplateList(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Отримуємо шаблони користувача
	templates, err := h.templateService.GetByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get templates",
		})
	}

	// Рендеримо шаблон
	return c.Render("template/list", fiber.Map{
		"Title":     "Мої шаблони",
		"Templates": templates,
	})
}

func (h *TemplateHandler) HandleTemplateCreate(c *fiber.Ctx) error {
	// Отримуємо ID користувача з контексту
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Парсимо дані з форми
	var template models.Template
	if err := c.BodyParser(&template); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Встановлюємо ID користувача
	template.UserID = userID

	// Створюємо шаблон
	if err := h.templateService.Create(c.Context(), &template); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create template",
		})
	}

	return c.Redirect("/templates")
}

func (h *TemplateHandler) HandleTemplateUpdate(c *fiber.Ctx) error {
	// Отримуємо ID шаблону
	templateID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	// Отримуємо шаблон
	template, err := h.templateService.GetByID(c.Context(), templateID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Template not found",
		})
	}

	// Парсимо дані з форми
	if err := c.BodyParser(template); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Оновлюємо шаблон
	if err := h.templateService.Update(c.Context(), template); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update template",
		})
	}

	return c.Redirect("/templates")
}

func (h *TemplateHandler) HandleTemplateDelete(c *fiber.Ctx) error {
	// Отримуємо ID шаблону
	templateID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	// Видаляємо шаблон
	if err := h.templateService.Delete(c.Context(), templateID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete template",
		})
	}

	return c.Redirect("/templates")
}
