package main

import (
	"github.com/suvidsahay/Factly/controllers"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	a := controllers.App{}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		)
	a.Run(":5000")
}