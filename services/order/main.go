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
	orderDB *gorm.DB

	ADMIN_KEY      string
	ORDER_SRV_HOST string
	ORDER_SRV_PORT string
)

type OrderServer struct {
	pb.OrderServiceServer
}

func GetDB() *gorm.DB {
	return orderDB
}

func loadENV() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading [order-service] .env file")
	}

	ADMIN_KEY = os.Getenv("ADMIN_KEY")
	if ADMIN_KEY == "" {
		log.Fatal("Invalid [ADMIN_KEY] - not found in [.env]")
	}

	ORDER_SRV_HOST = os.Getenv("ORDER_SRV_HOST")
	if ORDER_SRV_HOST == "" {
		log.Fatal("Invalid [ORDER_SRV_HOST] - not found in [.env]")
	}

	ORDER_SRV_PORT = os.Getenv("ORDER_SRV_PORT")
	if ORDER_SRV_PORT == "" {
		log.Fatal("Invalid [ORDER_SRV_PORT] - not found in [.env]")
	}
}

func loadDB() {

	dbUser := os.Getenv("ORDER_DB_USER")
	dbPass := os.Getenv("ORDER_DB_PASS")
	dbName := os.Getenv("ORDER_DB_NAME")
	dbHost := os.Getenv("ORDER_DB_HOST")
	dbPort := os.Getenv("ORDER_DB_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatalf("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to [order-service] database: %+v", err)
	} else {
		orderDB = db
		InitModels(db)
	}
}

func loadSRV() {

	lis, err := net.Listen("tcp", ":"+ORDER_SRV_PORT)

	if err != nil {
		log.Fatalf("Failed to start the [order-gRPC] %+v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, &OrderServer{})

	log.Printf("*** >>> [order-gRPC] server started at %+v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("*** >>> [order-gRPC] failed to start - %+v", err)
	}
}

func main() {
	loadENV()
	loadDB()
	loadSRV()
}
