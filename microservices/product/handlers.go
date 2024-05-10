package main

import (
	"github.com/gofiber/fiber/v2"
)

func createProduct(c *fiber.Ctx) error {
	// Proxy to product-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func deleteProduct(c *fiber.Ctx) error {
	// Proxy to product-service
	return c.SendStatus(fiber.StatusNotImplemented)
}
