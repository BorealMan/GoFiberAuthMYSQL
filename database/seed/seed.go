package seed

import (
	"fmt"

	"app/database"
	"app/models/user"

	"golang.org/x/crypto/bcrypt"
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

	// Converting password to hash
	pld := user.Email + user.Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(pld), 7)
	if err != nil {
		return
	}
	hash := string(bytes)
	user.Password = hash

	database.DB.Save(&user)
	fmt.Println("\tSucessfully Seeded User Table")
}
