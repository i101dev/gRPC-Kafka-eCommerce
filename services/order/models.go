package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// --------------------------------------------------------------------
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func InitModels(database *gorm.DB) {

	if err := database.AutoMigrate(&Order{}); err != nil {
		fmt.Println("Error initializing [models/Order.go]")
		log.Fatal(err)
	}

	fmt.Println("\n*** >>> Successfully initialized [models/Order.go]")

}

func DropModels(database *gorm.DB) {

	if err := database.Migrator().DropTable(&Order{}); err != nil {
		fmt.Println("Error dropping Order table")
		log.Fatal(err)
	}

	fmt.Println("*** >>> Successfully dropped Order table")
}
