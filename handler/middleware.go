package handler

import (
	"net/http"
	"strings"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/jwt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var IsForbidUser func(userId string) bool

func ParseToken(c *gin.Context) (*jwt.Payload, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	if token == "" {
		return nil, ierr.NewIError(ierr.NotAllowed, "need login")
	}
	var secret = viper.GetString("jwt.secret")
	var prefix = viper.GetString("jwt.manage_prefix")
	if strings.HasPrefix(token, prefix) {
		token = strings.TrimPrefix(token, prefix)
		secret = viper.GetString("jwt.manage_secret")
	}
	payload, err := jwt.ParseToken(token, secret)
	if err != nil {
		logger.Error("ParseToken parse token error",
			zap.String("token", token))
		return nil, ierr.NewIError(ierr.NotAllowed, err.Error())
	}

	err = payload.Valid()
	if err != nil {
		return nil, ierr.NewIError(ierr.TokenExpire, err.Error())
	}

	return payload, nil
}

func AuthToken() gin.HandlerFunc {
	//这里不管是什么身份只要登录了就可以
	return func(c *gin.Context) {
		info, err := ParseToken(c)
		if err != nil {
			//ResponseErr(c, err)
			c.AbortWithStatusJSON(401, err)
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		var prefix = viper.GetString("jwt.manage_prefix")
		if !strings.HasPrefix(token, prefix) {
			if IsForbidUser != nil && IsForbidUser(info.UserId) {
				//ResponseErrWithCode(c, http.StatusForbidden, ierr.NewIError(00000000, "permission denied"))
				c.AbortWithStatusJSON(http.StatusForbidden, ierr.NewIError(00000000, "permission denied"))
				return
			}
		}

		c.Set("userId", info.UserId)
		c.Set("userName", info.Name)
		c.Set("userRole", info.Role)
		c.Next()
	}
}

func AuthTokenRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := ParseToken(c)
		if err != nil {
			ResponseErr(c, err)
			c.AbortWithStatus(401)
			return
		}
		for _, role := range roles {
			if info.Role == role {
				c.Set("uid", info.UserId)
				c.Set("role", info.Role)
				c.Next()
				return
			}
		}

		c.AbortWithStatus(403)
		ResponseErr(c, ierr.NewIError(00000000, "permission denied"))
	}
}

func CheckRole(roles []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var role = c.GetString("role")
		for _, v := range roles {
			if v == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(403)
		ResponseErr(c, ierr.NewIError(ierr.PermissionDeined, "permission denied"))
	}
}
