package api

import (
	"app/api/routes/userRoutes"

	"github.com/gofiber/fiber/v2"

	"app/api/auth"
)

func testHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func SetRoutes(app *fiber.App) {
	api := app.Group("/api")
	userRoutes.SetUserRoutes(api)

	api.Get("/test", testHandler)

	api.Get("/", auth.ValidateJWT, testHandler)
}

func SetupAPI(app *fiber.App) {
	SetRoutes(app)
}
