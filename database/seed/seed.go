package seed

import (
	"fmt"

	"app/database"
	"app/models/user"
)

// Inserts Default Values
func SeedDB() {
	fmt.Println("Seeding Database:")
	UserSeed()
}

func UserSeed() {

	// Creating Admin
	var user user.User

	database.DB.First(&user)

	if user.Id != 0 {
		return
	}

	user.Email = "admin@admin.com"
	user.Password = "12345"
	user.Role = "admin"

	database.DB.Save(&user)
	fmt.Println("\tSucessfully Seeded User Table")
}
