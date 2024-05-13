package main

import (
	"context"
	"fmt"

	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	productClient pb.ProductServiceClient
	userClient    pb.UserServiceClient
)

func (s *OrderServer) OrderTest(ctx context.Context, req *pb.OrderTestReq) (*pb.OrderTestRes, error) {

	fmt.Println("*** >>> [order-gRPC] - server test message: ", req.Msg)

	return &pb.OrderTestRes{
		Msg: req.Msg,
	}, nil
}

func (s *OrderServer) OrderConn(ctx context.Context, req *pb.OrderConnReq) (*pb.OrderConnRes, error) {

	// --------------------------------------------------------------------------
	// Order service
	//
	if productClient == nil {

		productConn, productConnErr := grpc.Dial(req.ProductSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if productConnErr != nil {

			fmt.Printf("\n*** >>> Failed to dial the [order] server %v", productConnErr)

			return &pb.OrderConnRes{
				Msg: "Product service failed to connect to order service",
			}, productConnErr

		} else {
			productClient = pb.NewProductServiceClient(productConn)
		}
	}

	// --------------------------------------------------------------------------
	// Product service
	//
	if userClient == nil {

		userConn, userConnErr := grpc.Dial(req.UserSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if userConnErr != nil {

			fmt.Printf("\n*** >>> Failed to dial the [product] server %v", userConnErr)

			return &pb.OrderConnRes{
				Msg: "User service failed to connect to product service",
			}, userConnErr

		} else {
			userClient = pb.NewUserServiceClient(userConn)
		}
	}

	// --------------------------------------------------------------------------
	// Return
	return &pb.OrderConnRes{
		Msg: "[Order] service connected to [User] and [Product] services",
	}, nil
}

func (s *OrderServer) OrderPing(ctx context.Context, req *pb.OrderPingReq) (*pb.OrderPingRes, error) {

	fmt.Println("\n*** >>> [OrderPing] - Chk - 3")

	return &pb.OrderPingRes{
		Msg: "*** SUCCESS **** USER -> PRODUCT -> ORDER",
	}, nil
}

// func createOrder(c *fiber.Ctx) error {
// 	// Proxy to Order-service
// 	return c.SendStatus(fiber.StatusNotImplemented)
// }

// func deleteOrder(c *fiber.Ctx) error {
// 	// Proxy to Order-service
// 	return c.SendStatus(fiber.StatusNotImplemented)
// }
