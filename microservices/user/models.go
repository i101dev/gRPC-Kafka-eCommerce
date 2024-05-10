package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// --------------------------------------------------------------------
type Orders []Order
type Order struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	UUID string `json:"uuid"`
	Item string `json:"item"`
	Qty  int    `json:"qty"`
}

func (s Orders) Value() (driver.Value, error) {
	return json.Marshal(s)
}
func (s *Orders) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, s)
	}
	return errors.New("unsupported data type for scanning into [Orders]")
}

// --------------------------------------------------------------------
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `json:"uuid"`
	Role      string    `json:"role"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"password"`
	Email     string    `gorm:"unique" json:"email"`
	Orders    Orders    `gorm:"type:jsonb" json:"orders"`
	CreatedAt time.Time `json:"created_at"`
}
type UserLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func InitModels(database *gorm.DB) {

	if err := database.AutoMigrate(&User{}); err != nil {
		fmt.Println("Error initializing [models/User.go]")
		log.Fatal(err)
	}

	fmt.Println("\n*** >>> Successfully initialized [models/User.go]")

	if err := database.AutoMigrate(&Order{}); err != nil {
		fmt.Println("Error initializing [models/Order.go]")
		log.Fatal(err)
	}

	fmt.Println("*** >>> Successfully initialized [models/Order.go]")
}
