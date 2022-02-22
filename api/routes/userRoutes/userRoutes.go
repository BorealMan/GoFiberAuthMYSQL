package userRoutes

import (
	"github.com/gofiber/fiber/v2"

	"app/api/auth"
	"app/models/user"
)

func SetUserRoutes(api fiber.Router) {
	userGroup := api.Group("/user")
	userGroup.Post("/create", user.CreateUser)
	userGroup.Post("/login", user.Login)
	userGroup.Get("/getall", auth.ValidateJWT, auth.ValidateAdmin, user.GetAll)
	userGroup.Patch("/update", auth.ValidateJWT, user.Update)
}
