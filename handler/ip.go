package handler

import (
	"github.com/gin-gonic/gin"
)

var ClientIP = func(c *gin.Context) { // 获取真实IP
	realIP := c.Request.Header.Get("X-Real-IP")
	if realIP == "" {
		// 如果没有X-Real-IP，可以尝试获取X-Forwarded-For
		realIP = c.Request.Header.Get("X-Forwarded-For")
	}
	if realIP == "" {
		// 如果仍然没有，则获取客户端的IP
		realIP = c.ClientIP()
	}

	c.Set("clientIP", realIP)
	c.Next()
}
