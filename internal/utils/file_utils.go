package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// GetFileExtension повертає розширення файлу
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// SaveUploadedFile зберігає завантажений файл на диск
func SaveUploadedFile(file multipart.File, dst string) error {
	// Створюємо директорію, якщо вона не існує
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Створюємо файл
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Копіюємо вміст
	_, err = io.Copy(out, file)
	return err
}

// ParseCSV парсить CSV файл та повертає його вміст як масив рядків
func ParseCSV(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Дозволяє різну кількість полів у рядках

	return reader.ReadAll()
}

// FormatFileSize форматує розмір файлу для зручного відображення
func FormatFileSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%d KB", size/KB)
	case size < GB:
		return fmt.Sprintf("%d MB", size/MB)
	default:
		return fmt.Sprintf("%d GB", size/GB)
	}
}
