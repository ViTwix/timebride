package models

import "time"

// PageData базова структура для всіх сторінок
type PageData struct {
	Title string
	User  *UserViewModel
}

// UserViewModel структура для відображення даних користувача
type UserViewModel struct {
	ID          string
	Email       string
	FullName    string
	CompanyName string
	Phone       string
	Language    string
	Theme       string
	AvatarURL   string
}

// BookingViewModel структура для відображення даних бронювання
type BookingViewModel struct {
	ID          string
	Title       string
	EventType   string
	EventDate   time.Time
	Location    string
	Status      string
	StatusClass string
	PackageName string

	// Client details
	ClientName      string
	ClientPhone     string
	ClientInstagram string
	ClientNotes     string

	// Financial details
	Currency        string
	PriceTotal      float64
	PricePrepayment float64
	PriceExtra      float64
	PriceProfit     float64
	PriceLeftToPay  float64

	// Team details
	TeamPayments []TeamPaymentViewModel

	// Additional details
	DeadlineDays    int
	ContractFileURL string
	DeliveryPageURL string
	DiskCode        string
}

// TeamMemberViewModel структура для відображення даних члена команди
type TeamMemberViewModel struct {
	ID          string
	Name        string
	Email       string
	Role        string
	AccessLevel string
	CanEdit     bool
}

// TeamPaymentViewModel структура для відображення оплати члену команди
type TeamPaymentViewModel struct {
	TeamMemberID string
	Name         string
	Role         string
	Amount       float64
}

// PriceTemplateViewModel структура для відображення цінового шаблону
type PriceTemplateViewModel struct {
	ID              string
	Name            string
	EventType       string
	Currency        string
	PriceTotal      float64
	PricePrepayment float64
	DeadlineDays    int
	Comment         string
	TeamRoles       []TeamRoleAmount
}

// TeamRoleAmount структура для відображення ролі та оплати
type TeamRoleAmount struct {
	Role   string
	Amount float64
}

// DashboardStats структура для відображення статистики
type DashboardStats struct {
	TotalBookings   int
	ActiveBookings  int
	UpcomingEvents  int
	EventsThisMonth int
	TotalEarned     float64
	PendingPayments float64
	TeamMembers     int
	StorageUsedGB   float64
	StorageLimitGB  int
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
	case string(BookingStatusPending):
		return "bg-yellow"
	case string(BookingStatusBooked):
		return "bg-blue"
	case string(BookingStatusEditing):
		return "bg-purple"
	case string(BookingStatusReady):
		return "bg-green"
	case string(BookingStatusArchived):
		return "bg-gray"
	default:
		return "bg-gray"
	}
}
