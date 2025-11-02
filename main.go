package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/Jake-Schuler/minecraft-management-thingy/config"
	"github.com/Jake-Schuler/minecraft-management-thingy/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed static/*
var static embed.FS

//go:embed templates/*
var templates embed.FS

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	rcon := config.InitRCON()

	// Initialize Gin router
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(templates, "templates/*")))
	r.StaticFS("/static", http.FS(static))

	// Setup routes
	handlers.SetupRoutes(r, rcon)

	// Start server
	if err := r.Run("127.0.0.1:8080"); err != nil {
		panic("failed to start server")
	}
}
