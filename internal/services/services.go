package services

import (
	"timebride/internal/services/auth"
	"timebride/internal/services/booking"
	"timebride/internal/services/client"
	"timebride/internal/services/price"
	"timebride/internal/services/storage"
	"timebride/internal/services/team"
	"timebride/internal/services/template"
	"timebride/internal/services/user"
)

// Services містить всі сервіси програми
type Services struct {
	Auth     auth.IAuthService
	User     user.IUserService
	Booking  booking.IBookingService
	Client   client.IClientService
	Team     team.ITeamService
	Price    price.IPriceService
	Storage  storage.IStorageService
	Template template.ITemplateService
}

// NewServices створює нову структуру Services
func NewServices(
	authSvc auth.IAuthService,
	userSvc user.IUserService,
	bookingSvc booking.IBookingService,
	clientSvc client.IClientService,
	teamSvc team.ITeamService,
	priceSvc price.IPriceService,
	storageSvc storage.IStorageService,
	templateSvc template.ITemplateService,
) *Services {
	return &Services{
		Auth:     authSvc,
		User:     userSvc,
		Booking:  bookingSvc,
		Client:   clientSvc,
		Team:     teamSvc,
		Price:    priceSvc,
		Storage:  storageSvc,
		Template: templateSvc,
	}
}
