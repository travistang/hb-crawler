package api

import "github.com/gin-gonic/gin"

func reportError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

func sendJSONPayload(c *gin.Context, statusCode int, payload any) {
	c.JSON(statusCode, gin.H{
		"ok":   true,
		"data": payload,
	})
}
