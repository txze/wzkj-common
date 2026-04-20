package handler

import (
	"net/http"
	"strings"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/jwt"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type JWTConfig struct {
	Secret       string `mapstructure:"secret"`
	ManageSecret string `mapstructure:"manage_secret"`
	ManagePrefix string `mapstructure:"manage_prefix"`
}

var jwtConfig *JWTConfig

var IsForbidUser func(userId string) bool

func InitJWT(cfg *JWTConfig) {
	jwtConfig = cfg
}

func ParseToken(c *gin.Context) (*jwt.Payload, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	if token == "" {
		return nil, ierr.NewIError(ierr.NotAllowed, "need login")
	}

	if jwtConfig == nil {
		return nil, ierr.NewIError(ierr.NotAllowed, "JWT未初始化")
	}

	var secret = jwtConfig.Secret
	var prefix = jwtConfig.ManagePrefix
	if strings.HasPrefix(token, prefix) {
		token = strings.TrimPrefix(token, prefix)
		secret = jwtConfig.ManageSecret
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

func ParseTokenWithConfig(c *gin.Context, cfg *JWTConfig) (*jwt.Payload, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	if token == "" {
		return nil, ierr.NewIError(ierr.NotAllowed, "need login")
	}

	if cfg == nil {
		return nil, ierr.NewIError(ierr.NotAllowed, "JWT配置不能为空")
	}

	var secret = cfg.Secret
	var prefix = cfg.ManagePrefix
	if strings.HasPrefix(token, prefix) {
		token = strings.TrimPrefix(token, prefix)
		secret = cfg.ManageSecret
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
	return func(c *gin.Context) {
		info, err := ParseToken(c)
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if jwtConfig != nil {
			var prefix = jwtConfig.ManagePrefix
			if !strings.HasPrefix(token, prefix) {
				if IsForbidUser != nil && IsForbidUser(info.UserId) {
					c.AbortWithStatusJSON(http.StatusForbidden, ierr.NewIError(00000000, "permission denied"))
					return
				}
			}
		}

		c.Set("userId", info.UserId)
		c.Set("userName", info.Name)
		c.Set("userRole", info.Role)
		c.Next()
	}
}

func AuthTokenWithConfig(cfg *JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := ParseTokenWithConfig(c, cfg)
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		var prefix = cfg.ManagePrefix
		if !strings.HasPrefix(token, prefix) {
			if IsForbidUser != nil && IsForbidUser(info.UserId) {
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

func AuthTokenRoleWithConfig(cfg *JWTConfig, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := ParseTokenWithConfig(c, cfg)
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
