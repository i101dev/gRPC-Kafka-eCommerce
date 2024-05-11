package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/i101dev/gRPC-kafka-eCommerce/auth"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

var (
	userSrv    string
	userClient pb.UserServiceClient
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	userSrv = os.Getenv("USER_SRV")
	userConn, err := grpc.Dial(userSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to dial the [user] server %v", err)
	} else {
		defer userConn.Close()
		userClient = pb.NewUserServiceClient(userConn)
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
	// Routers
	app.Get("/", Handle_get_index)
	app.Post("/test", POST_user_test)
	app.Post("/login", POST_user_auth)
	app.Post("/register", POST_user_register)

	// --------------------------------------------------------------------------
	// User API
	api_user := app.Group("/user")
	api_user.Use(auth.ValidateJWT)
	api_user.Get("/products", Handle_get_products)
	api_user.Get("/inventory", Handle_get_inventory)

	// --------------------------------------------------------------------------
	// Admin API
	api_admin := app.Group("/admin")
	api_admin.Use(auth.ValidateJWT)
	api_admin.Get("/users", auth.RequireRole("admin"), Handle_get_users)
	api_admin.Get("/orders", auth.RequireRole("admin"), Handle_get_orders)

	// --------------------------------------------------------------------------
	// Launch server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func fiberErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"Fiber error": err.Error(),
	})
}
