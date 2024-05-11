package main

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

func parseError(err error) string {
	errMsg := err.Error()
	parts := strings.Split(errMsg, "desc = ")
	if len(parts) > 1 {
		return parts[1]
	}
	return errMsg

}

func Handle_get_index(c *fiber.Ctx) error {
	return c.SendString("Welcome to the E-commerce Order Processing Platform")
}

func POST_user_test(c *fiber.Ctx) error {

	var request struct {
		Msg string `json:"msg"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	testReq := &pb.TestReq{
		Msg: request.Msg,
	}

	testRes, err := userClient.Test(context.Background(), testReq)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	log.Printf("\n*** >>> [call_Test] - client - response - %s", testRes.Msg)

	return c.JSON(testRes)
}

func POST_user_auth(c *fiber.Ctx) error {

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	authRes, authErr := userClient.AuthUser(context.Background(), &pb.AuthReq{
		Username: request.Username,
		Password: request.Password,
	})

	if authErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": parseError(authErr)})
	}

	return c.JSON(authRes)
}

func POST_user_register(c *fiber.Ctx) error {

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Referral string `json:"referral"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	registerRes, registerErr := userClient.RegisterUser(context.Background(), &pb.RegisterReq{
		Username: request.Username,
		Password: request.Password,
		Referral: request.Referral,
		Email:    request.Email,
	})

	if registerErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": parseError(registerErr)})
	}

	return c.JSON(registerRes)
}

func Handle_get_users(c *fiber.Ctx) error {
	// Proxy to user-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func Handle_get_products(c *fiber.Ctx) error {
	// Proxy to product-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func Handle_get_orders(c *fiber.Ctx) error {
	// Proxy to order-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

func Handle_get_inventory(c *fiber.Ctx) error {
	// Proxy to inventory-service
	return c.SendStatus(fiber.StatusNotImplemented)
}
