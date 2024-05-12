package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/i101dev/gRPC-kafka-eCommerce/auth"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port     string
	adminKey string

	orderSrv    string
	orderConn   *grpc.ClientConn
	orderClient pb.OrderServiceClient

	productSrv    string
	productConn   *grpc.ClientConn
	productClient pb.ProductServiceClient

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

	orderSrv = os.Getenv("ORDER_SRV")
	if orderSrv == "" {
		log.Fatal("Invalid [ORDER_SRV] - not found in [.env]")
	}
	productSrv = os.Getenv("PRODUCT_SRV")
	if productSrv == "" {
		log.Fatal("Invalid [PRODUCT_SRV] - not found in [.env]")
	}
	userSrv = os.Getenv("USER_SRV")
	if userSrv == "" {
		log.Fatal("Invalid [USER_SRV] - not found in [.env]")
	}
}

func loadGRPC() {

	// --------------------------------------------------------------------------
	// order service
	//
	orderConn, err := grpc.Dial(orderSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial the [order] server %v", err)
	} else {
		orderClient = pb.NewOrderServiceClient(orderConn)
	}

	// --------------------------------------------------------------------------
	// inventory service
	//
	productConn, err := grpc.Dial(productSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial the [product] server %v", err)
	} else {
		productClient = pb.NewProductServiceClient(productConn)
	}

	// --------------------------------------------------------------------------
	// user service
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

	defer orderConn.Close()
	defer productConn.Close()
	defer userConn.Close()

	app := fiberApp()

	// --------------------------------------------------------------------------
	// Routers
	app.Get("/", GET_index)
	app.Post("/order/test", POST_order_test)
	app.Post("/product/test", POST_product_test)
	app.Post("/user/test", POST_user_test)

	app.Post("/auth", POST_AuthUser)
	app.Post("/register", POST_RegisterUser)

	// --------------------------------------------------------------------------
	// User API
	api_user := app.Group("/user")
	api_user.Use(auth.ValidateJWT)
	api_user.Get("/products", GET_products)
	api_user.Get("/inventory", GET_inventory)

	// --------------------------------------------------------------------------
	// Admin API
	api_admin := app.Group("/admin")
	api_admin.Use(auth.ValidateJWT)
	api_admin.Get("/users", auth.RequireRole("admin"), GET_users)
	api_admin.Get("/orders", auth.RequireRole("admin"), GET_orders)

	// --------------------------------------------------------------------------
	// Launch server
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
