package main

import (
	"github.com/Webservice-Rathje/Main-Backend/endpoints"
	"github.com/Webservice-Rathje/Main-Backend/endpoints/kundenmanagement"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", endpoints.DefaultEndpoint)
	app.Post("/login", endpoints.LoginController)
	app.Post("/checkToken", endpoints.CheckTokenController)
	app.Post("/logout", endpoints.LogoutController)

	// Kunden-Management
	app.Get("/ws/kunden-management/createAccount", kundenmanagement.CreateAccountWebsocket())

	app.Listen(":8080")
}
