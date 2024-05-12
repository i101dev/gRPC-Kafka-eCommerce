package main

import (
	"context"
	"fmt"

	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

func (s *OrderServer) OrderTest(ctx context.Context, req *pb.OrderTestReq) (*pb.OrderTestRes, error) {

	fmt.Println("*** >>> [order-gRPC] - server test message: ", req.Msg)

	return &pb.OrderTestRes{
		Msg: req.Msg,
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
