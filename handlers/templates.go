package handlers

import (
	"github.com/gin-gonic/gin"
)

func IndexPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{})
	}
}

func ResetPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "passwordreset.tmpl", gin.H{})
	}
}

func KickAllPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "kickall.tmpl", gin.H{})
	}
}

func KickPlayerPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "kickplayer.tmpl", gin.H{})
	}
}

func ForceLoginPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "forcelogin.tmpl", gin.H{})
	}
}