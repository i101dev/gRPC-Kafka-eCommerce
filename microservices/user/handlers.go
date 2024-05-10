package main

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/i101dev/gRPC-kafka-eCommerce/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func registerUser(c *fiber.Ctx) error {

	db := c.Locals("db").(*gorm.DB)

	user := new(User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user.CreatedAt = time.Now()
	user.Password = hashedPassword
	user.UUID = uuid.New().String()
	user.Role = "customer"

	if err := db.Create(&user).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Username or email already exists"})
	}

	return c.Status(http.StatusCreated).JSON(user)
}

func loginUser(c *fiber.Ctx) error {

	db := c.Locals("db").(*gorm.DB)

	payload := new(UserLoginParams)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	var userDat User
	if err := db.Where("username = ?", payload.Username).First(&userDat).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	if !checkPassword(userDat.Password, payload.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := auth.GenerateJWT(userDat.Username, userDat.Role)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create JWT token"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
