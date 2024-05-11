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
	port     string
	adminKey string

	userSrv    string
	userConn   *grpc.ClientConn
	userClient pb.UserServiceClient
)

func loadENV() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	adminKey = os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		log.Fatal("Invalid [ADMIN_KEY] - not found in [.env]")
	}

	userSrv = os.Getenv("USER_SRV")
	if userSrv == "" {
		log.Fatal("Invalid [USER_SRV] - not found in [.env]")
	}
}

func loadGRPC() {

	// --------------------------------------------------------------------------
	// User service
	//
	userConn, err := grpc.Dial(userSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial the [user] server %v", err)
	} else {
		userClient = pb.NewUserServiceClient(userConn)
	}
}

func fiberApp() *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Fiber error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	return app
}

func main() {

	loadENV()
	loadGRPC()

	defer userConn.Close()

	app := fiberApp()

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
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
