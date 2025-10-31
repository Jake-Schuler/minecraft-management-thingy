package main

import (
	"github.com/Jake-Schuler/minecraft-management-thingy/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	config.InitSheets()
	config.InitRCON()
}