package utils

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
	"time"
)

// TemplateFunctions returns a map of functions that can be used in templates
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"formatMoney":        SafeFormatMoney,
		"formatDate":         FormatDate,
		"formatTime":         FormatTime,
		"formatDateTime":     FormatDateTime,
		"formatReminderTime": FormatReminderTime,
		"firstChar":          FirstChar,
		"add":                func(a, b int) int { return a + b },
		"sub":                func(a, b int) int { return a - b },
		"subtract":           func(a, b int) int { return a - b },
		"mul":                func(a, b int) int { return a * b },
		"div":                func(a, b int) int { return a / b },
		"hasSliceElem":       HasSliceElem,
		"join":               strings.Join,
		"split":              strings.Split,
		"contains":           strings.Contains,
		"title":              strings.Title,
		"lower":              strings.ToLower,
		"upper":              strings.ToUpper,
		"seq":                Seq,
		"date":               DateFormat,
	}
}

// DateFormat formats a date according to the provided format
// It supports both time.Time values and the special string "now"
func DateFormat(format string, value interface{}) string {
	var t time.Time

	switch v := value.(type) {
	case time.Time:
		t = v
	case string:
		if v == "now" {
			t = time.Now()
		} else {
			// Try to parse the string as a time
			parsed, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return ""
			}
			t = parsed
		}
	default:
		return ""
	}

	return t.Format(format)
}

// SafeFormatMoney safely formats a value as money, handling nil and other types
func SafeFormatMoney(amount interface{}) string {
	if amount == nil {
		return "0,00"
	}

	// Convert to float64 if possible
	var floatAmount float64
	switch v := amount.(type) {
	case float64:
		floatAmount = v
	case float32:
		floatAmount = float64(v)
	case int:
		floatAmount = float64(v)
	case int64:
		floatAmount = float64(v)
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "0,00"
		}
		floatAmount = parsed
	default:
		return "0,00"
	}

	return FormatMoney(floatAmount)
}

// FormatMoney formats a float64 as a money string (e.g. 1234.56 -> "1 234,56")
func FormatMoney(amount float64) string {
	// Round to 2 decimal places
	rounded := math.Round(amount*100) / 100

	// Convert to string with 2 decimal places
	str := strconv.FormatFloat(rounded, 'f', 2, 64)

	// Split into integer and decimal parts
	parts := strings.Split(str, ".")

	// Format integer part with thousands separator
	intPart := parts[0]
	var formatted strings.Builder

	for i, char := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			formatted.WriteRune(' ')
		}
		formatted.WriteRune(char)
	}

	// Return formatted amount with comma as decimal separator
	return formatted.String() + "," + parts[1]
}

// FormatDate formats a time.Time as a date string (e.g. "01.01.2023")
func FormatDate(t time.Time) string {
	return t.Format("02.01.2006")
}

// FormatTime formats a time.Time as a time string (e.g. "15:04")
func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

// FormatDateTime formats a time.Time as a date and time string (e.g. "01.01.2023 15:04")
func FormatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

// HasSliceElem checks if a slice contains a specific element
func HasSliceElem(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

// FormatReminderTime formats a reminder time value (e.g. 24, 48, 72) as human-readable text
func FormatReminderTime(hours int) string {
	if hours == 0 {
		return "Не надсилати"
	} else if hours < 24 {
		return fmt.Sprintf("%d годин", hours)
	} else {
		days := hours / 24
		return fmt.Sprintf("%d %s", days, pluralizeDays(days))
	}
}

// pluralizeDays returns the correct form of the word "day" based on the count
func pluralizeDays(days int) string {
	lastDigit := days % 10
	lastTwoDigits := days % 100

	if lastTwoDigits >= 11 && lastTwoDigits <= 19 {
		return "днів"
	}

	switch lastDigit {
	case 1:
		return "день"
	case 2, 3, 4:
		return "дні"
	default:
		return "днів"
	}
}

// FirstChar returns the first character of a string
func FirstChar(s string) string {
	if s == "" {
		return ""
	}

	// Trim whitespace first
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// Get the first rune (to handle UTF-8 characters properly)
	for _, r := range s {
		return string(r)
	}

	return ""
}

// Seq generates a sequence of integers from 1 to n
func Seq(n int) []int {
	if n <= 0 {
		return []int{}
	}

	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i + 1
	}

	return result
}
