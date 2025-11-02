package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func StopPanelHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		os.Exit(0)
	}
}
