package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
)

func fiberErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
func generateJWT(userID string, secretKey string) (string, error) {

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
func validateJWT(c *fiber.Ctx) error {

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

func main() {

	if err := godotenv.Load("../config/.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// --------------------------------------------------------------------------
	// INITIALIZE FIBER
	//
	app := fiber.New(fiber.Config{
		ErrorHandler: fiberErrorHandler,
	})

	// --------------------------------------------------------------------------
	// MIDDLEWARES
	//
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// --------------------------------------------------------------------------
	// ROUTES
	//
	api := app.Group("/api")

	api.Get("/", handle_get_index)
	api.Get("/users", validateJWT, handle_get_users)
	api.Get("/orders", validateJWT, haandl_get_orders)
	api.Get("/products", validateJWT, handle_get_products)
	api.Get("/inventory", validateJWT, handle_get_inventory)

	api.Post("/login", handle_post_login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func handle_get_index(c *fiber.Ctx) error {
	return c.SendString("Welcome to the E-commerce Order Processing Platform")
}

func handle_post_login(c *fiber.Ctx) error {

	jwtSecretStr := os.Getenv("JWT_SECRET")

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	tokenString, err := generateJWT(request.Username, jwtSecretStr)

	if err != nil {
		fmt.Print(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create JWT"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func handle_get_users(c *fiber.Ctx) error {
	// Proxy to user-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func handle_get_products(c *fiber.Ctx) error {
	// Proxy to product-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func haandl_get_orders(c *fiber.Ctx) error {
	// Proxy to order-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func handle_get_inventory(c *fiber.Ctx) error {
	// Proxy to inventory-service
	return c.SendStatus(fiber.StatusNotImplemented)
}
