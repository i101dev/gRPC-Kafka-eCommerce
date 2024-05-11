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
	userDB *gorm.DB

	ADMIN_KEY string
	SRV_HOST  string
	SRV_PORT  string
)

type UserServer struct {
	pb.UserServiceServer
}

func GetDB() *gorm.DB {
	return userDB
}

func loadENV() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading [user-service] .env file")
	}

	ADMIN_KEY = os.Getenv("ADMIN_KEY")
	if ADMIN_KEY == "" {
		log.Fatal("Invalid [ADMIN_KEY] - not found in [.env]")
	}

	SRV_HOST = os.Getenv("SRV_HOST")
	if SRV_HOST == "" {
		log.Fatal("Invalid [SRV_HOST] - not found in [.env]")
	}

	SRV_PORT = os.Getenv("SRV_PORT")
	if SRV_PORT == "" {
		log.Fatal("Invalid [SRV_PORT] - not found in [.env]")
	}
}

func loadDB() {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatalf("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to [user-service] database: %+v", err)
	} else {
		userDB = db
		InitModels(db)
	}
}

func loadSRV() {

	lis, err := net.Listen("tcp", ":"+SRV_PORT)

	if err != nil {
		log.Fatalf("Failed to start the [user-gRPC] %+v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	log.Printf("*** >>> [user-gRPC] server started at %+v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("*** >>> [user-gRPC] failed to start - %+v", err)
	}
}

func main() {
	loadENV()
	loadDB()
	loadSRV()
}
