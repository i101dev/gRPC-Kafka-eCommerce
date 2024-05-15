package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

var (
	productDB *gorm.DB

	ADMIN_KEY        string
	KAFKA_URI        string
	KAFKA_TOPIC      string
	PRODUCT_SRV_HOST string
	PRODUCT_SRV_PORT string
)

type ProductServer struct {
	pb.ProductServiceServer
}

func GetDB() *gorm.DB {
	return productDB
}

func loadENV() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading [user-service] .env file")
	}

	ADMIN_KEY = os.Getenv("ADMIN_KEY")
	if ADMIN_KEY == "" {
		log.Fatal("Invalid [ADMIN_KEY] - not found in [.env]")
	}

	KAFKA_URI = os.Getenv("KAFKA_URI")
	if KAFKA_URI == "" {
		log.Fatal("Invalid [KAFKA_URI] - not found in [.env]")
	}

	KAFKA_TOPIC = os.Getenv("KAFKA_TOPIC")
	if KAFKA_TOPIC == "" {
		log.Fatal("Invalid [KAFKA_TOPIC] - not found in [.env]")
	}

	PRODUCT_SRV_HOST = os.Getenv("PRODUCT_SRV_HOST")
	if PRODUCT_SRV_HOST == "" {
		log.Fatal("Invalid [PRODUCT_SRV_HOST] - not found in [.env]")
	}

	PRODUCT_SRV_PORT = os.Getenv("PRODUCT_SRV_PORT")
	if PRODUCT_SRV_PORT == "" {
		log.Fatal("Invalid [PRODUCT_SRV_PORT] - not found in [.env]")
	}
}

func loadDB() {

	dbUser := os.Getenv("PRODUCT_DB_USER")
	dbPass := os.Getenv("PRODUCT_DB_PASS")
	dbHost := os.Getenv("PRODUCT_DB_HOST")
	dbPort := os.Getenv("PRODUCT_DB_PORT")
	dbName := os.Getenv("PRODUCT_DB_NAME")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatalf("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to [product-service] database: %+v", err)
	} else {
		productDB = db
		InitModels(db)
	}
}

func loadSRV() {

	lis, err := net.Listen("tcp", ":"+PRODUCT_SRV_PORT)

	if err != nil {
		log.Fatalf("Failed to start the [product-gRPC] %+v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(grpcServer, &ProductServer{})

	log.Printf("*** >>> [product-gRPC] server started at %+v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("*** >>> [product-gRPC] failed to start - %+v", err)
	}
}

func main() {
	loadENV()
	loadDB()
	loadSRV()
}
