package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, role string) (string, error) {

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT secret not defined")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
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

	user_id, user_id_Ok := (*claims)["user_id"].(string)
	role, roleOk := (*claims)["role"].(string)

	if !user_id_Ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid [user_id]"})
	}

	if !roleOk {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid [role]"})
	}

	c.Locals("user_id", user_id)
	c.Locals("role", role)

	fmt.Printf("*** >>> Username: %s, Role: %s\n", user_id, role)

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
func RequireRole(role string) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		userRole := c.Locals("role").(string)

		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access forbidden: insufficient role"})
		}

		return c.Next()
	}
}
