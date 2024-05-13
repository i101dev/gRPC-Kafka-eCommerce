package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/uuid"
	auth "github.com/i101dev/gRPC-kafka-eCommerce/auth"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

var (
	orderClient   pb.OrderServiceClient
	productClient pb.ProductServiceClient
)

func (s *UserServer) UserJoin(ctx context.Context, req *pb.UserJoinReq) (*pb.UserJoinRes, error) {

	db := GetDB()

	var role string
	switch req.Referral {
	case "0x1":
		role = "admin"
	default:
		role = "cust"
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return &pb.UserJoinRes{}, fmt.Errorf("failed to hash password")
	}

	user := User{
		CreatedAt: time.Now(),
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Role:      role,
		UUID:      uuid.New().String(),
	}

	if err := db.Create(&user).Error; err != nil {
		return &pb.UserJoinRes{}, fmt.Errorf("failed to create user")
	}

	return &pb.UserJoinRes{
		UserId: user.UUID,
	}, nil
}

func (s *UserServer) AuthUser(ctx context.Context, req *pb.UserAuthReq) (*pb.UserAuthRes, error) {

	db := GetDB()

	userDat := new(User)

	if err := db.Where("username = ?", req.Username).First(&userDat).Error; err != nil {
		return &pb.UserAuthRes{}, fmt.Errorf("not authorizeed - %+v", err)
	}

	if !checkPassword(userDat.Password, req.Password) {
		return &pb.UserAuthRes{}, fmt.Errorf("not authorizeed")
	}

	token, err := auth.GenerateJWT(userDat.Username, userDat.Role)

	if err != nil {
		return &pb.UserAuthRes{}, fmt.Errorf("failed to generate token")
	}

	return &pb.UserAuthRes{
		Token: token,
	}, nil
}

func (s *UserServer) UserTest(ctx context.Context, req *pb.UserTestReq) (*pb.UserTestRes, error) {

	fmt.Println("*** >>> [user-gRPC] - server test message: ", req.Msg)

	return &pb.UserTestRes{
		Msg: req.Msg,
	}, nil
}

func (s *UserServer) UserConn(ctx context.Context, req *pb.UserConnReq) (*pb.UserConnRes, error) {

	// --------------------------------------------------------------------------
	// Order service
	//
	if orderClient == nil {

		orderConn, orderConnErr := grpc.Dial(req.OrderSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if orderConnErr != nil {

			fmt.Printf("\n*** >>> Failed to dial the [order] server %v", orderConnErr)

			return &pb.UserConnRes{
				Msg: "User service failed to connect to order service",
			}, orderConnErr

		} else {
			orderClient = pb.NewOrderServiceClient(orderConn)
		}
	}

	// --------------------------------------------------------------------------
	// Product service
	//
	if productClient == nil {

		productConn, productConnErr := grpc.Dial(req.ProductSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if productConnErr != nil {

			fmt.Printf("\n*** >>> Failed to dial the [product] server %v", productConnErr)

			return &pb.UserConnRes{
				Msg: "User service failed to connect to product service",
			}, productConnErr

		} else {
			productClient = pb.NewProductServiceClient(productConn)
		}
	}

	// --------------------------------------------------------------------------
	// Return
	return &pb.UserConnRes{
		Msg: "[User] service connected to [Order] and [Product] services",
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// --------------------------------------------------------------------------
// Ping test routes
//

func (s *UserServer) UserPingProduct(ctx context.Context, req *pb.UserPingProductReq) (*pb.UserPingProductRes, error) {

	pingReq := &pb.ProductPingOrderReq{
		Msg: req.Msg,
	}

	pingRes, pingErr := productClient.ProductPingOrder(context.Background(), pingReq)
	if pingErr != nil {
		return &pb.UserPingProductRes{}, fmt.Errorf("failed to ping [Product] service")

	}

	fmt.Println("\n*** >>> [UserPingProduct] - Chk - 1")

	return &pb.UserPingProductRes{
		Msg: pingRes.Msg,
	}, nil
}
