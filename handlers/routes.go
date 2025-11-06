package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func SetupRoutes(r *gin.Engine, rconConn *rcon.Conn) {
	r.GET("/stop-panel", StopPanelHandler())
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"user": os.Getenv("ADMIN_PASSWORD"),
	}))

	authorized.GET("/", IndexPageHandler())
	authorized.GET("/reset", ResetPageHandler())
	authorized.GET("/users", GetUsersHandler(rconConn))
	authorized.GET("/kickall", KickAllPageHandler())
	authorized.POST("/kickall", KickAllPlayersHandler(rconConn))
	authorized.POST("/reset-password", ResetUserPassHandler(rconConn))
	authorized.GET("/kick-player", KickPlayerPageHandler())
	authorized.POST("/kick-player", KickPlayerHandler(rconConn))
	authorized.GET("/force-login", ForceLoginPageHandler())
	authorized.POST("/force-login", ForceLoginHandler(rconConn))
}
