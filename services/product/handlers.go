package main

import (
	"context"
	"fmt"

	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

func (s *ProductServer) ProductTest(ctx context.Context, req *pb.ProductTestReq) (*pb.ProductTestRes, error) {

	fmt.Println("*** >>> [product-gRPC] - server test message: ", req.Msg)

	return &pb.ProductTestRes{
		Msg: req.Msg,
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
