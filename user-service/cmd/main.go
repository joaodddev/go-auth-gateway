package main

import (
	"user-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/users", handler.GetUsers)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "user service ok"})
	})

	r.Run(":8001")
}