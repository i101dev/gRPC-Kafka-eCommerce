package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

func parseError(err error) map[string]interface{} {

	errMsg := err.Error()
	parts := strings.Split(errMsg, "desc = ")

	if len(parts) > 1 {
		return fiber.Map{"error": parts[1]}
	}

	return fiber.Map{"error": errMsg}

}

// --------------------------------------------------------------------------
// Testing
func POST_inventory_test(c *fiber.Ctx) error {

	var request struct {
		Msg string `json:"msg"`
	}

	fmt.Println("Chk - 1")

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	fmt.Println("Chk - 2")
	testReq := &pb.InventoryTestReq{
		Msg: request.Msg,
	}

	testRes, testErr := inventoryClient.InventoryTest(context.Background(), testReq)

	fmt.Println("Chk - 3")
	if testErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(parseError(testErr))
	}

	fmt.Println("Chk - 4")
	return c.JSON(testRes)
}
func POST_order_test(c *fiber.Ctx) error {

	var request struct {
		Msg string `json:"msg"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	testReq := &pb.OrderTestReq{
		Msg: request.Msg,
	}

	testRes, testErr := orderClient.OrderTest(context.Background(), testReq)

	if testErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(parseError(testErr))
	}

	return c.JSON(testRes)
}
func POST_product_test(c *fiber.Ctx) error {

	var request struct {
		Msg string `json:"msg"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	testReq := &pb.ProductTestReq{
		Msg: request.Msg,
	}

	testRes, testErr := productClient.ProductTest(context.Background(), testReq)

	if testErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(parseError(testErr))
	}

	return c.JSON(testRes)
}
func POST_user_test(c *fiber.Ctx) error {

	var request struct {
		Msg string `json:"msg"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	testReq := &pb.UserTestReq{
		Msg: request.Msg,
	}

	testRes, testErr := userClient.UserTest(context.Background(), testReq)

	if testErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(parseError(testErr))
	}

	return c.JSON(testRes)
}

// --------------------------------------------------------------------------
// User handlers

func GET_index(c *fiber.Ctx) error {
	return c.SendString("Welcome to the E-commerce Order Processing Platform")
}

func POST_AuthUser(c *fiber.Ctx) error {

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
		return c.Status(fiber.StatusBadRequest).JSON(parseError(authErr))
	}

	return c.JSON(authRes)
}

func POST_RegisterUser(c *fiber.Ctx) error {

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
		return c.Status(fiber.StatusBadRequest).JSON(parseError(registerErr))
	}

	return c.JSON(registerRes)
}

func GET_users(c *fiber.Ctx) error {
	// Proxy to user-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

// --------------------------------------------------------------------------
// Product handlers
func GET_products(c *fiber.Ctx) error {
	// Proxy to product-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

// --------------------------------------------------------------------------
// Order handlers
func GET_orders(c *fiber.Ctx) error {
	// Proxy to order-service
	return c.SendStatus(fiber.StatusNotImplemented)
}

// --------------------------------------------------------------------------
// Inventory handlers
func GET_inventory(c *fiber.Ctx) error {
	// Proxy to inventory-service
	return c.SendStatus(fiber.StatusNotImplemented)
}
