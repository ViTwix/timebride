package client

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/repositories"
	"timebride/internal/services/storage"
)

// Service реалізує інтерфейс IClientService
type Service struct {
	clientRepo repositories.ClientRepository
	fileRepo   repositories.FileRepository
	storage    storage.IStorageService
}

// NewService створює новий екземпляр сервісу клієнтів
func NewService(
	clientRepo repositories.ClientRepository,
	fileRepo repositories.FileRepository,
	storage storage.IStorageService,
) *Service {
	return &Service{
		clientRepo: clientRepo,
		fileRepo:   fileRepo,
		storage:    storage,
	}
}

// Get отримує клієнта за ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	return s.clientRepo.GetByID(ctx, id)
}

// Create створює нового клієнта
func (s *Service) Create(ctx context.Context, client *models.Client) error {
	return s.clientRepo.Create(ctx, client)
}

// Update оновлює існуючого клієнта
func (s *Service) Update(ctx context.Context, client *models.Client) error {
	return s.clientRepo.Update(ctx, client)
}

// Delete видаляє клієнта
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.clientRepo.Delete(ctx, id)
}

// List отримує список клієнтів за фільтром
func (s *Service) List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*models.ClientListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filter := map[string]interface{}{"user_id": userID}
	clients, err := s.clientRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := s.clientRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Застосовуємо пагінацію
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(clients) {
		clients = []*models.Client{}
	} else {
		if end > len(clients) {
			end = len(clients)
		}
		clients = clients[start:end]
	}

	return &models.ClientListResult{
		Items:      clients,
		TotalItems: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetByUserID отримує всіх клієнтів користувача
func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Client, error) {
	return s.clientRepo.GetByUserID(ctx, userID)
}

// UploadAvatar завантажує аватар для клієнта
func (s *Service) UploadAvatar(ctx context.Context, clientID uuid.UUID, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("file is required")
	}

	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return "", err
	}

	// Генеруємо унікальне ім'я файлу
	ext := filepath.Ext(file.Filename)
	filename := "avatars/" + uuid.New().String() + ext

	// Завантажуємо файл
	url, err := s.storage.UploadFile(ctx, file, filename)
	if err != nil {
		return "", err
	}

	// Створюємо запис про файл
	fileModel := &models.File{
		UserID:      client.UserID,
		Name:        file.Filename,
		Path:        filename,
		URL:         url,
		Size:        file.Size,
		ContentType: file.Header.Get("Content-Type"),
		Type:        models.FileTypeAvatar,
	}

	if err := s.fileRepo.Create(ctx, fileModel); err != nil {
		// Якщо не вдалося створити запис, видаляємо файл
		_ = s.storage.DeleteFile(ctx, filename)
		return "", err
	}

	// Оновлюємо аватар клієнта
	client.Avatar = url
	if err := s.clientRepo.Update(ctx, client); err != nil {
		return "", err
	}

	return url, nil
}

// DeleteAvatar видаляє аватар клієнта
func (s *Service) DeleteAvatar(ctx context.Context, clientID uuid.UUID) error {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	if client.Avatar == "" {
		return nil
	}

	// Знаходимо файл за URL
	filter := map[string]interface{}{"url": client.Avatar}
	files, err := s.fileRepo.List(ctx, filter)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		file := files[0]
		if err := s.storage.DeleteFile(ctx, file.Path); err != nil {
			return err
		}
		if err := s.fileRepo.Delete(ctx, file.ID); err != nil {
			return err
		}
	}

	client.Avatar = ""
	return s.clientRepo.Update(ctx, client)
}

// GetCategories отримує всі категорії клієнтів
func (s *Service) GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.clientRepo.GetCategories(ctx, userID)
}

// GetSources отримує всі джерела клієнтів
func (s *Service) GetSources(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.clientRepo.GetSources(ctx, userID)
}

// ListClients отримує список клієнтів за опціями
func (s *Service) ListClients(ctx context.Context, userID uuid.UUID, options models.ClientListOptions) (*models.ClientListResult, error) {
	if options.Page < 1 {
		options.Page = 1
	}
	if options.PageSize < 1 {
		options.PageSize = 10
	}

	filter := map[string]interface{}{
		"user_id": userID,
	}

	if options.Category != "" {
		filter["category"] = options.Category
	}
	if options.Source != "" {
		filter["source"] = options.Source
	}
	if options.Search != "" {
		filter["search"] = options.Search
	}

	clients, err := s.clientRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := s.clientRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Застосовуємо пагінацію
	start := (options.Page - 1) * options.PageSize
	end := start + options.PageSize
	if start >= len(clients) {
		clients = []*models.Client{}
	} else {
		if end > len(clients) {
			end = len(clients)
		}
		clients = clients[start:end]
	}

	return &models.ClientListResult{
		Items:      clients,
		TotalItems: total,
		Page:       options.Page,
		PageSize:   options.PageSize,
	}, nil
}

// GetClient отримує клієнта за ID з перевіркою userID
func (s *Service) GetClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) (*models.Client, error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if client.UserID != userID {
		return nil, errors.New("client does not belong to user")
	}

	return client, nil
}

// CreateClient створює нового клієнта і повертає його
func (s *Service) CreateClient(ctx context.Context, client *models.Client) (*models.Client, error) {
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

// UpdateClient оновлює існуючого клієнта і повертає оновлену версію
func (s *Service) UpdateClient(ctx context.Context, client *models.Client) (*models.Client, error) {
	existing, err := s.clientRepo.GetByID(ctx, client.ID)
	if err != nil {
		return nil, err
	}

	if existing.UserID != client.UserID {
		return nil, errors.New("client does not belong to user")
	}

	if err := s.clientRepo.Update(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

// DeleteClient видаляє клієнта з перевіркою userID
func (s *Service) DeleteClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) error {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	if client.UserID != userID {
		return errors.New("client does not belong to user")
	}

	return s.clientRepo.Delete(ctx, clientID)
}
