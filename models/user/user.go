package user

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"app/api/auth"
	"app/database"
)

type User struct {
	Id          uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Email       string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password    string `json:"password" gorm:"not null" validate:"required,min=4,max=100"`
	CreatedAt   uint64 `json:"createdat" gorm:"autoCreateTime"`
	LastUpdated uint64 `json:"lastupdated" gorm:"autoUpdateTime"`
	Role        string `json:"role" gorm:"required"`
}

func Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user User
	user.Email = email
	user.Password = password

	validate := validator.New()
	err := validate.Struct(user)
	// Check Form
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// Check if User Exists
	database.DB.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user doesn't exist",
		})
	}
	// Check User Password
	pld := user.Email + password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pld)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid password",
		})
	}
	// Create JWT Token For User
	userId := user.Id
	userRole := user.Role
	t, err := auth.IssueJWT(fmt.Sprintf("%d", userId), userRole)
	// Check The Token Didn't Explode
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Clear Password
	user.Password = ""
	return c.JSON(fiber.Map{"token": t, "user": user})
}

func CreateUser(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user User
	user.Email = email
	user.Password = password
	user.Role = "free"

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// Hash Password
	pld := email + password
	bytes, err := bcrypt.GenerateFromPassword([]byte(pld), 6)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	hash := string(bytes)
	user.Password = hash
	// Saving User To DB
	database.DB.Save(&user)
	if user.Id == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "account already exists",
		})
	}
	// Create JWT Token For User
	userId := user.Id
	userRole := user.Role
	t, err := auth.IssueJWT(fmt.Sprintf("%d", userId), userRole)
	// Check The Token Didn't Explode
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Clear Password
	user.Password = ""
	return c.JSON(fiber.Map{"token": t, "user": user})
}

func GetAll(c *fiber.Ctx) error {
	var users []User
	database.DB.Find(&users)
	return c.Status(200).JSON(fiber.Map{
		"users": users,
	})
}

func Update(c *fiber.Ctx) error {
	userId := c.GetReqHeaders()["Userid"]
	email := c.FormValue("email")

	var user User
	database.DB.First(&user, userId)

	user.Email = email

	validate := validator.New()
	err := validate.Struct(user)
	// Check Form
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// Try to save the new fields
	database.DB.Save(&user)
	// Check if it updated unique values
	database.DB.First(&user, userId)
	if email != user.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is already in use",
		})
	}

	user.Password = ""
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"user": user,
	})
}
