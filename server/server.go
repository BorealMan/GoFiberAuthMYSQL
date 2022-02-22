package server

import (
	"fmt"
	"log"

	"app/api"
	"app/database"
	"app/database/seed"
	"app/models/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	var err error
	conn := database.DbURL(database.BuildDBConfig())
	database.DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("\nDatabase:\nSucessfully Connected")
	database.DB.AutoMigrate(&user.User{})
	fmt.Println("Sucessfully Migrated")
	seed.SeedDB()
}

func Start() {
	app := fiber.New()

	// Connect to DB
	initDB()

	// Set Routes & Middleware
	app.Use(recover.New())
	api.SetupAPI(app)

	// Start API
	fmt.Println("\nStarting app at http://localhost:5000")
	log.Fatal(app.Listen(":5000"))
}
