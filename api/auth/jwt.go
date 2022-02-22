package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var JWTSecretKey = []byte("Secret Key Goes Here")

func IssueJWT(userId string, userRole string) (string, error) {

	claims := jwt.MapClaims{
		"userId": userId,
		"role":   userRole,
		"exp":    time.Now().Add(time.Minute * 10).Unix(), // 60 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(JWTSecretKey)

	return t, err
}

func ValidateJWT(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	authHeader := headers["Authorization"]
	authToken := strings.Split(authHeader, " ")

	if len(authToken) != 2 {
		return c.SendStatus(403)
	}

	token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return c.SendStatus(403)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			c.Request().Header.Set("Role", fmt.Sprintf("%s", claims["role"]))
			c.Request().Header.Set("Userid", fmt.Sprintf("%s", claims["userId"]))
		}
		return c.Next()
	}

	return c.SendStatus(fiber.StatusInternalServerError)
}

func ValidateAdmin(c *fiber.Ctx) error {
	if c.GetReqHeaders()["Role"] != "admin" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
