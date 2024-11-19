package main

import (
	"os"

	"arch/ikeppu/github.com/internal/link"
	"arch/ikeppu/github.com/internal/stat"
	"arch/ikeppu/github.com/internal/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
}
