package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/txze/wzkj-common/logger"
)

// UserAgentMiddleware 记录请求日志，标记是否可能来自常见软件，并打印所有请求参数（包括 JSON 请求体）
func UserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ua := c.Request.UserAgent()

		// 生成 traceID
		traceID := uuid.New().String()
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "traceID", traceID))
		c.Set("traceID", traceID)

		// 常见软件的 UA 特征
		knownAgents := []string{
			"curl", "Postman", "python-requests", "okhttp", "axios", "java", "wget",
		}

		isSoftware := false
		for _, b := range knownAgents {
			if strings.Contains(strings.ToLower(ua), b) {
				isSoftware = true
				break
			}
		}

		// 打印所有请求参数
		params := c.Request.URL.Query()
		c.Request.ParseForm() // 解析表单参数（包含 POST form-data 和 x-www-form-urlencoded）
		for k, v := range c.Request.Form {
			params[k] = v
		}

		// 读取并打印 JSON 请求体
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			if len(bodyBytes) > 0 {
				var jsonData map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
					for k, v := range jsonData {
						params[k] = []string{toString(v)}
					}
				}
			}
			// 重新放回请求体，避免后续处理读不到
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		var userId = c.GetString("userId")

		// 使用 Zap 记录日志
		duration := time.Since(start)
		status := c.Writer.Status()
		logger.WithTrace(c).Info("HTTP Request",
			logger.String("userId", userId),
			logger.String("traceID", traceID),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("clientIP", c.ClientIP()),
			logger.Int("status", status),
			logger.Any("duration", duration),
			logger.String("userAgent", ua),
			logger.Bool("isSoftware", isSoftware),
			logger.Any("params", params),
			logger.Any("headrs", c.Request.Header),
		)

		c.Next()
	}
}

// 辅助函数：将 interface{} 转换为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		bytes, _ := json.Marshal(val)
		return string(bytes)
	}
}
