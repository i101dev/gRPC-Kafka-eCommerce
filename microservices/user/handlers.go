package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/i101dev/gRPC-kafka-eCommerce/auth"
	pb "github.com/i101dev/gRPC-kafka-eCommerce/proto"
)

func (s *UserServer) RegisterUser(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {

	db := GetDB()

	user := new(User)

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return &pb.RegisterRes{}, fmt.Errorf("failed to hash password")
	}

	var role string
	switch req.Referral {
	case "0x1":
		role = "admin"
	default:
		role = "customer"
	}

	user.Email = req.Email
	user.Username = req.Username
	user.CreatedAt = time.Now()
	user.Password = hashedPassword
	user.UUID = uuid.New().String()
	user.Role = role

	if err := db.Create(&user).Error; err != nil {
		return &pb.RegisterRes{}, fmt.Errorf("failed to create user")
	}

	return &pb.RegisterRes{
		UserId: user.UUID,
	}, nil
}

func (s *UserServer) AuthUser(ctx context.Context, req *pb.AuthReq) (*pb.AuthRes, error) {

	db := GetDB()

	userDat := new(User)

	if err := db.Where("username = ?", req.Username).First(&userDat).Error; err != nil {
		return &pb.AuthRes{}, fmt.Errorf("not authorizeed - %+v", err)
	}

	if !checkPassword(userDat.Password, req.Password) {
		return &pb.AuthRes{}, fmt.Errorf("not authorizeed")
	}

	token, err := auth.GenerateJWT(userDat.Username, userDat.Role)

	if err != nil {
		return &pb.AuthRes{}, fmt.Errorf("failed to generate token")
	}

	return &pb.AuthRes{
		Token: token,
	}, nil
}

func (s *UserServer) Test(ctx context.Context, req *pb.TestReq) (*pb.TestRes, error) {

	fmt.Println("*** >>> [user-gRPC] - server test message: ", req.Msg)

	return &pb.TestRes{
		Msg: req.Msg,
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
