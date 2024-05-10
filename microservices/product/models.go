package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// --------------------------------------------------------------------
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func InitModels(database *gorm.DB) {

	if err := database.AutoMigrate(&Product{}); err != nil {
		fmt.Println("Error initializing [models/Product.go]")
		log.Fatal(err)
	}

	fmt.Println("\n*** >>> Successfully initialized [models/Product.go]")

}
