package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, secretKey string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"admin":   true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func ValidateJWT(c *fiber.Ctx) error {

	jwtSecret := os.Getenv("JWT_SECRET")
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	if len(authHeader) < len("Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization format"})
	}

	token, err := parseJWT(authHeader, jwtSecret)

	if err != nil {
		fmt.Printf("JWT parse error: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT"})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT"})
	}

	claims, ok := token.Claims.(*jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
	}

	username, usernameOk := (*claims)["user_id"].(string)
	expiry, expiryOk := (*claims)["exp"].(float64)
	admin, adminOk := (*claims)["admin"].(bool)

	if !usernameOk || !expiryOk || !adminOk {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to extract claims"})
	}

	fmt.Printf("*** >>> Username: %s, Expiry: %d\n", username, int64(expiry))
	fmt.Printf("*** >>> Admin: %+v\n", admin)

	return c.Next()
}
func parseJWT(authHeader string, jwtSecret string) (*jwt.Token, error) {

	tokenStr := authHeader[len("Bearer "):]

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	return token, err
}
