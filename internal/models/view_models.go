package models

import "time"

// PageData базова структура для всіх сторінок
type PageData struct {
	Title string
	User  *UserViewModel
}

// UserViewModel структура для відображення даних користувача
type UserViewModel struct {
	ID        string
	Email     string
	FullName  string
	Role      string
	Initials  string
	AvatarURL string
}

// BookingViewModel структура для відображення даних бронювання
type BookingViewModel struct {
	ID          string
	ClientName  string
	EventType   string
	StartTime   time.Time
	Status      string
	StatusClass string
}

// DashboardStats структура для відображення статистики
type DashboardStats struct {
	TotalBookings   int
	ActiveBookings  int
	UpcomingEvents  int
	EventsThisMonth int
	TotalTemplates  int
	ActiveTemplates int
	TotalFiles      int
	TotalSize       string
}

// DashboardData структура для відображення даних на дашборді
type DashboardData struct {
	PageData
	Stats          DashboardStats
	RecentBookings []BookingViewModel
}

// GetStatusClass повертає CSS клас для статусу бронювання
func (b *BookingViewModel) GetStatusClass() string {
	switch b.Status {
	case "pending":
		return "bg-yellow"
	case "confirmed":
		return "bg-green"
	case "cancelled":
		return "bg-red"
	case "completed":
		return "bg-blue"
	default:
		return "bg-gray"
	}
}
