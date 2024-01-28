package main

import (
	"familytree-gen/dev/assets"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/countries", func(c *gin.Context) {
		c.JSON(200, assets.TFRData)
	})

	router.Run(":8080")
}
