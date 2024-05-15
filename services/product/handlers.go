package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i101dev/gRPC-kafka-eCommerce/kafka"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	orderClient pb.OrderServiceClient
	userClient  pb.UserServiceClient
)

func (s *ProductServer) ProductTest(ctx context.Context, req *pb.ProductTestReq) (*pb.ProductTestRes, error) {

	// fmt.Println("*** >>> [product-gRPC] - server test message: ", req)

	kafkaErr := pushMsgToKafka(req.Msg, req.Val)

	return &pb.ProductTestRes{
		Msg: req.Msg,
	}, kafkaErr
}

func (s *ProductServer) ProductConn(ctx context.Context, req *pb.ProductConnReq) (*pb.ProductConnRes, error) {

	// --------------------------------------------------------------------------
	// Order service
	//
	if orderClient == nil {

		orderConn, orderConnErr := grpc.Dial(req.OrderSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if orderConnErr != nil {

			fmt.Printf("\n*** >>> Failed to dial the [order] server %v", orderConnErr)

			return &pb.ProductConnRes{
				Msg: "Product service failed to connect to order service",
			}, orderConnErr

		} else {
			orderClient = pb.NewOrderServiceClient(orderConn)
		}
	}

	// --------------------------------------------------------------------------
	// Product service
	//
	if userClient == nil {

		userConn, userConnErr := grpc.Dial(req.UserSrv, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if userConnErr != nil {

			fmt.Printf("\n*** >>> Failed to dial the [product] server %v", userConnErr)

			return &pb.ProductConnRes{
				Msg: "User service failed to connect to product service",
			}, userConnErr

		} else {
			userClient = pb.NewUserServiceClient(userConn)
		}
	}

	// --------------------------------------------------------------------------
	// Return
	return &pb.ProductConnRes{
		Msg: "[Product] service connected to [Order] and [Product] services",
	}, nil
}

// func createProduct(c *fiber.Ctx) error {
// 	// Proxy to product-service
// 	return c.SendStatus(fiber.StatusNotImplemented)
// }

// func deleteProduct(c *fiber.Ctx) error {
// 	// Proxy to product-service
// 	return c.SendStatus(fiber.StatusNotImplemented)
// }

// --------------------------------------------------------------------------
// Ping test routes
//

func (s *ProductServer) ProductPingOrder(ctx context.Context, req *pb.ProductPingOrderReq) (*pb.ProductPingOrderRes, error) {

	pingReq := &pb.OrderPingReq{
		Msg: req.Msg,
	}

	pingRes, pingErr := orderClient.OrderPing(context.Background(), pingReq)
	if pingErr != nil {
		return &pb.ProductPingOrderRes{}, fmt.Errorf("[product] service failed to ping [order] service")

	}

	fmt.Println("\n*** >>> [ProductPingOrder] - Chk - 2")

	return &pb.ProductPingOrderRes{
		Msg: pingRes.Msg,
	}, nil
}

// --------------------------------------------------------------------------
// Kafka
func pushMsgToKafka(msg string, val int64) error {

	msgInBytes, err := json.Marshal(kafka.ProductMsg{
		Msg: msg,
		Val: val,
	})

	if err != nil {
		return kafka.HandleKafkaError(err, "Error marshalling to JSON")
	}

	err = kafka.PushMsgToQueue(KAFKA_URI, KAFKA_TOPIC, msgInBytes)

	return err
}
