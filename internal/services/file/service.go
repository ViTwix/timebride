package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"timebride/internal/config"
	"timebride/internal/models"
	"timebride/internal/repositories"
)

type Service struct {
	repo      repositories.FileRepository
	s3Client  *s3.Client
	s3Config  config.StorageConfig
	cdnConfig config.CDNConfig
}

type Handler struct {
	service *Service
}

func NewService(repo repositories.FileRepository, s3Client *s3.Client, cfg config.StorageConfig) *Service {
	return &Service{
		repo:      repo,
		s3Client:  s3Client,
		s3Config:  cfg,
		cdnConfig: cfg.CDN,
	}
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (s *Service) UploadFile(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*models.File, error) {
	// Generate a unique key for the file while preserving the original name
	ext := filepath.Ext(file.Filename)
	safeName := sanitizeFileName(strings.TrimSuffix(file.Filename, ext))
	uniqueID := uuid.New().String()
	key := fmt.Sprintf("%s/%s/%s-%s%s",
		userID.String(),
		time.Now().Format("2006-01"),
		safeName,
		uniqueID,
		ext,
	)

	// Create a temporary file for upload
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	if _, err = io.Copy(tempFile, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Reset the file pointer to the beginning
	if _, err = tempFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek temp file: %w", err)
	}

	// Set metadata for the file
	metadata := map[string]string{
		"original-name": file.Filename,
		"upload-date":   time.Now().Format(time.RFC3339),
		"content-type":  file.Header.Get("Content-Type"),
	}

	// Marshal metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Upload to S3/Backblaze
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Config.Backblaze.Bucket),
		Key:         aws.String(key),
		Body:        tempFile,
		ContentType: aws.String(file.Header.Get("Content-Type")),
		Metadata:    metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create database record
	fileModel := &models.File{
		ID:       uuid.New(),
		UserID:   userID,
		FileName: file.Filename,
		FileType: file.Header.Get("Content-Type"),
		FileSize: file.Size,
		Bucket:   s.s3Config.Backblaze.Bucket,
		Key:      key,
		CDNURL:   fmt.Sprintf("%s://%s/%s", s.cdnConfig.Protocol, s.cdnConfig.Domain, key),
		Metadata: string(metadataJSON),
	}

	if err := s.repo.Create(ctx, fileModel); err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return fileModel, nil
}

func (s *Service) DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, string, error) {
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get file: %w", err)
	}

	result, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(file.Bucket),
		Key:    aws.String(file.Key),
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to download file: %w", err)
	}

	return result.Body, file.FileName, nil
}

func (s *Service) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	_, err = s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(file.Bucket),
		Key:    aws.String(file.Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	if err := s.repo.Delete(ctx, fileID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

func (s *Service) ListFiles(ctx context.Context, userID uuid.UUID) ([]*models.File, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) GetFileByID(ctx context.Context, id uuid.UUID) (*models.File, error) {
	return s.repo.GetByID(ctx, id)
}

// Helper function to sanitize file names
func sanitizeFileName(name string) string {
	// Replace unsafe characters with a hyphen
	reg := regexp.MustCompile(`[^a-zA-Z0-9-_.]`)
	return reg.ReplaceAllString(name, "-")
}

// GetTotalFiles returns the total number of files
func (s *Service) GetTotalFiles(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx, map[string]interface{}{})
}

// GetTotalSize returns the total size of all files
func (s *Service) GetTotalSize(ctx context.Context) (int64, error) {
	files, err := s.repo.List(ctx, map[string]interface{}{})
	if err != nil {
		return 0, err
	}

	var totalSize int64
	for _, file := range files {
		totalSize += file.FileSize
	}
	return totalSize, nil
}
