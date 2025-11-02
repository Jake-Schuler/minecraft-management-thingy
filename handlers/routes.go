package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func SetupRoutes(r *gin.Engine, rconConn *rcon.Conn) {
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"user": os.Getenv("ADMIN_PASSWORD"),
	}))

	authorized.GET("/", HomeHandler())
	authorized.GET("/users", GetUsersHandler(rconConn))
	authorized.POST("/reset-password", ResetUserPassHandler(rconConn))
}
