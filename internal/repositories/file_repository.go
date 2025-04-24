package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.File, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.File, error)
	Update(ctx context.Context, file *models.File) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*models.File, error)
	BatchCreate(ctx context.Context, files []*models.File) error
	BatchDelete(ctx context.Context, ids []uuid.UUID) error
}

type fileRepository struct {
	db    *gorm.DB
	cache *fileCache
}

type fileCache struct {
	mu        sync.RWMutex
	items     map[uuid.UUID]*models.File
	userFiles map[uuid.UUID][]*models.File
	expiry    time.Duration
}

func newFileCache() *fileCache {
	return &fileCache{
		items:     make(map[uuid.UUID]*models.File),
		userFiles: make(map[uuid.UUID][]*models.File),
		expiry:    5 * time.Minute,
	}
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{
		db:    db,
		cache: newFileCache(),
	}
}

func (r *fileRepository) Create(ctx context.Context, file *models.File) error {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return err
	}

	// Оновлюємо кеш
	r.cache.mu.Lock()
	r.cache.items[file.ID] = file
	if files, ok := r.cache.userFiles[file.UserID]; ok {
		r.cache.userFiles[file.UserID] = append(files, file)
	}
	r.cache.mu.Unlock()

	return nil
}

func (r *fileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.File, error) {
	// Спочатку перевіряємо кеш
	r.cache.mu.RLock()
	if file, ok := r.cache.items[id]; ok {
		r.cache.mu.RUnlock()
		return file, nil
	}
	r.cache.mu.RUnlock()

	// Якщо немає в кеші, читаємо з БД
	var file models.File
	if err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Зберігаємо в кеш
	r.cache.mu.Lock()
	r.cache.items[file.ID] = &file
	r.cache.mu.Unlock()

	return &file, nil
}

func (r *fileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.File, error) {
	// Спочатку перевіряємо кеш
	r.cache.mu.RLock()
	if files, ok := r.cache.userFiles[userID]; ok {
		r.cache.mu.RUnlock()
		return files, nil
	}
	r.cache.mu.RUnlock()

	// Якщо немає в кеші, читаємо з БД
	var files []*models.File
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}

	// Зберігаємо в кеш
	r.cache.mu.Lock()
	r.cache.userFiles[userID] = files
	for _, file := range files {
		r.cache.items[file.ID] = file
	}
	r.cache.mu.Unlock()

	return files, nil
}

func (r *fileRepository) Update(ctx context.Context, file *models.File) error {
	// Оновлюємо в БД
	if err := r.db.WithContext(ctx).Model(&models.File{}).Where("id = ?", file.ID).Updates(file).Error; err != nil {
		return err
	}

	// Оновлюємо кеш
	r.cache.mu.Lock()
	r.cache.items[file.ID] = file
	if files, ok := r.cache.userFiles[file.UserID]; ok {
		for i, f := range files {
			if f.ID == file.ID {
				files[i] = file
				break
			}
		}
	}
	r.cache.mu.Unlock()

	return nil
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Отримуємо файл для визначення userID
	file, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Видаляємо з БД
	if err := r.db.WithContext(ctx).Delete(&models.File{}, "id = ?", id).Error; err != nil {
		return err
	}

	// Видаляємо з кешу
	r.cache.mu.Lock()
	delete(r.cache.items, id)
	if files, ok := r.cache.userFiles[file.UserID]; ok {
		newFiles := make([]*models.File, 0, len(files)-1)
		for _, f := range files {
			if f.ID != id {
				newFiles = append(newFiles, f)
			}
		}
		r.cache.userFiles[file.UserID] = newFiles
	}
	r.cache.mu.Unlock()

	return nil
}

func (r *fileRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.File{})

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	return count, query.Count(&count).Error
}

func (r *fileRepository) List(ctx context.Context, filter map[string]interface{}) ([]*models.File, error) {
	var files []*models.File
	query := r.db.WithContext(ctx)

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}

func (r *fileRepository) BatchCreate(ctx context.Context, files []*models.File) error {
	if len(files) == 0 {
		return nil
	}

	// Створюємо записи в транзакції
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, file := range files {
			if err := tx.Create(file).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Оновлюємо кеш
	r.cache.mu.Lock()
	for _, file := range files {
		r.cache.items[file.ID] = file
		if existingFiles, ok := r.cache.userFiles[file.UserID]; ok {
			r.cache.userFiles[file.UserID] = append(existingFiles, file)
		} else {
			r.cache.userFiles[file.UserID] = []*models.File{file}
		}
	}
	r.cache.mu.Unlock()

	return nil
}

func (r *fileRepository) BatchDelete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	// Отримуємо файли для визначення userID
	var files []*models.File
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&files).Error; err != nil {
		return err
	}

	// Видаляємо записи в транзакції
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&models.File{}).Error; err != nil {
		return err
	}

	// Оновлюємо кеш
	r.cache.mu.Lock()
	for _, file := range files {
		delete(r.cache.items, file.ID)
		if userFiles, ok := r.cache.userFiles[file.UserID]; ok {
			newFiles := make([]*models.File, 0, len(userFiles))
			for _, f := range userFiles {
				if f.ID != file.ID {
					newFiles = append(newFiles, f)
				}
			}
			r.cache.userFiles[file.UserID] = newFiles
		}
	}
	r.cache.mu.Unlock()

	return nil
}
