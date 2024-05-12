package main

import (
	"context"
	"fmt"

	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

func (s *InventoryServer) InventoryTest(ctx context.Context, req *pb.InventoryTestReq) (*pb.InventoryTestRes, error) {

	fmt.Println("*** >>> [Inventory-gRPC] - server test message: ", req.Msg)

	return &pb.InventoryTestRes{
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
