package handlers

import (
	"net/http"

	"github.com/Jake-Schuler/minecraft-management-thingy/services"
	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func GetUsersHandler(rconConn *rcon.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := services.GetUsers(rconConn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
func ResetUserPassHandler(rconConn *rcon.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := services.ResetUserPassword(rconConn, req.Username, req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
func KickPlayerHandler(rconConn *rcon.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := services.KickPlayer(rconConn, req.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
func KickAllPlayersHandler(rconConn *rcon.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := services.KickAllPlayers(rconConn); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}