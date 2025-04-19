package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"mime/multipart"

	appConfig "timebride/internal/config"
	"timebride/internal/models"
	"timebride/internal/repositories"
)

func createTestMultipartFile(filePath string) (*multipart.FileHeader, error) {
	// Create buffer for multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create form part for file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// Copy file content
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Close writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create reader from buffer
	reader := multipart.NewReader(body, writer.Boundary())

	// Read form
	form, err := reader.ReadForm(32 << 20) // 32MB max
	if err != nil {
		return nil, err
	}

	// Get FileHeader
	if files := form.File["file"]; len(files) > 0 {
		return files[0], nil
	}

	return nil, fmt.Errorf("no file found in form")
}

func TestBackblazeConnection(t *testing.T) {
	// Load configuration
	cfg, err := appConfig.Load()
	require.NoError(t, err, "Failed to load configuration")

	// Print configuration values for debugging
	fmt.Printf("Backblaze B2 Configuration:\n")
	fmt.Printf("Account ID: %s\n", cfg.Storage.Backblaze.AccountID)
	fmt.Printf("Application Key: %s\n", cfg.Storage.Backblaze.ApplicationKey)
	fmt.Printf("Bucket: %s\n", cfg.Storage.Backblaze.Bucket)
	fmt.Printf("Endpoint: %s\n", cfg.Storage.Backblaze.Endpoint)
	fmt.Printf("Region: %s\n", cfg.Storage.Backblaze.Region)

	// Create S3 client
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Storage.Backblaze.Endpoint,
			SigningRegion:     cfg.Storage.Backblaze.Region,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Storage.Backblaze.AccountID,
			cfg.Storage.Backblaze.ApplicationKey,
			"",
		)),
	)
	require.NoError(t, err, "Failed to create AWS config")

	s3Client := s3.NewFromConfig(awsCfg)

	// Create temporary database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	require.NoError(t, err, "Failed to create database")

	// Migrate models
	err = db.AutoMigrate(&models.File{})
	require.NoError(t, err, "Failed to migrate database")

	// Get absolute path to test file
	workDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	testFilePath := filepath.Join(workDir, "..", "..", "testdata", "test_video.mp4")
	fmt.Printf("Looking for file at path: %s\n", testFilePath)

	// Check if file exists
	_, err = os.Stat(testFilePath)
	require.NoError(t, err, "Test file not found: %s", testFilePath)

	// Create multipart file for test
	fileHeader, err := createTestMultipartFile(testFilePath)
	require.NoError(t, err, "Failed to create test multipart file")

	// Create repository and service
	repo := repositories.NewFileRepository(db)
	storageService := NewService(repo, s3Client, cfg.Storage)

	// Create test user
	userID := uuid.New()

	var uploadedFile *models.File

	t.Run("Upload File", func(t *testing.T) {
		// Upload file
		uploadedFile, err = storageService.UploadFile(context.Background(), fileHeader, userID)
		require.NoError(t, err, "Failed to upload file")
		assert.NotNil(t, uploadedFile, "File was not uploaded")
		assert.Equal(t, userID, uploadedFile.UserID, "Invalid UserID")
		assert.NotEmpty(t, uploadedFile.Key, "File key is empty")
		assert.NotEmpty(t, uploadedFile.CDNURL, "CDN URL is empty")
	})

	t.Run("Download File", func(t *testing.T) {
		// Download file
		reader, fileName, err := storageService.DownloadFile(context.Background(), uploadedFile.ID)
		require.NoError(t, err, "Failed to download file")
		defer reader.Close()

		assert.Equal(t, fileHeader.Filename, fileName, "Invalid file name")

		// Read downloaded file content
		downloadedData, err := io.ReadAll(reader)
		require.NoError(t, err, "Failed to read downloaded file")

		// Get original file size
		originalFile, err := os.Open(testFilePath)
		require.NoError(t, err, "Failed to open original file")
		defer originalFile.Close()

		originalData, err := io.ReadAll(originalFile)
		require.NoError(t, err, "Failed to read original file")

		assert.Equal(t, len(originalData), len(downloadedData), "File size mismatch")
		assert.Equal(t, originalData, downloadedData, "File content mismatch")
	})

	t.Run("List Files", func(t *testing.T) {
		// Get list of files
		files, err := storageService.ListFiles(context.Background(), userID)
		require.NoError(t, err, "Failed to get file list")
		assert.Len(t, files, 1, "Invalid number of files")
		assert.Equal(t, uploadedFile.ID, files[0].ID, "File ID mismatch")
	})

	t.Run("Delete File", func(t *testing.T) {
		// Delete file
		err = storageService.DeleteFile(context.Background(), uploadedFile.ID)
		require.NoError(t, err, "Failed to delete file")

		// Check if file is deleted
		files, err := storageService.ListFiles(context.Background(), userID)
		require.NoError(t, err, "Failed to get file list")
		assert.Len(t, files, 0, "File was not deleted")

		// Check that we can't download deleted file
		_, _, err = storageService.DownloadFile(context.Background(), uploadedFile.ID)
		assert.Error(t, err, "Expected error when trying to download deleted file")
	})
}
