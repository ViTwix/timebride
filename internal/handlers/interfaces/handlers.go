package interfaces

import "github.com/gofiber/fiber/v2"

// IAuthHandler визначає інтерфейс для обробки запитів аутентифікації
type IAuthHandler interface {
	ShowLoginPage(c *fiber.Ctx) error
	HandleLogin(c *fiber.Ctx) error
	ShowRegisterPage(c *fiber.Ctx) error
	HandleRegister(c *fiber.Ctx) error
	HandleLogout(c *fiber.Ctx) error
	OAuthRedirect(c *fiber.Ctx) error
	OAuthCallback(c *fiber.Ctx) error
}

// IUserHandler визначає інтерфейс для обробки запитів користувачів
type IUserHandler interface {
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// IBookingHandler визначає інтерфейс для обробки запитів бронювань
type IBookingHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// IClientHandler визначає інтерфейс для обробки запитів клієнтів
type IClientHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// ITeamHandler визначає інтерфейс для обробки запитів команди
type ITeamHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// IPriceHandler визначає інтерфейс для обробки запитів прайс-листів
type IPriceHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// IStorageHandler визначає інтерфейс для обробки запитів сховища
type IStorageHandler interface {
	List(c *fiber.Ctx) error
	Upload(c *fiber.Ctx) error
	Download(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
