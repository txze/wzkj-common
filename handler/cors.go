package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Cors = func(c *gin.Context) {
	fmt.Println("-----cors")
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "*")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
