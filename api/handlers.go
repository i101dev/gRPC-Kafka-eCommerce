package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/i101dev/gRPC-kafka-eCommerce/auth"
)

func Handle_get_index(c *fiber.Ctx) error {
	return c.SendString("Welcome to the E-commerce Order Processing Platform")
}

func Handle_post_login(c *fiber.Ctx) error {

	jwtSecretStr := os.Getenv("JWT_SECRET")

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Referral string `json:"referral"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	var role string
	switch request.Referral {
	case "0x1":
		role = "admin"
	default:
		role = "customer"
	}

	tokenString, err := auth.GenerateJWT(request.Username, role, jwtSecretStr)

	if err != nil {
		fmt.Print(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create JWT"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func Handle_get_users(c *fiber.Ctx) error {
	// Proxy to user-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func Handle_get_products(c *fiber.Ctx) error {
	// Proxy to product-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func haandl_get_orders(c *fiber.Ctx) error {
	// Proxy to order-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func Handle_get_inventory(c *fiber.Ctx) error {
	// Proxy to inventory-service
	return c.SendStatus(fiber.StatusNotImplemented)
}
