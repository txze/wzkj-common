package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, data ...interface{}) {
	var ret interface{}
	if len(data) == 0 {
		ret = "OK"
	} else {
		ret = data[0]
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": ret,
	})
}

// 使用这个方法需要返回自己定义的错误
func ResponseErr(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, err)
}

// 使用这个方法需要返回自己定义的错误
func ResponseErrWithCode(c *gin.Context, code int, err error) {
	c.JSON(code, err)
}

func ResponseXMLSuccess(c *gin.Context, data interface{}) {
	c.XML(http.StatusOK, data)
}

func ResponseXMLErr(c *gin.Context, err error) {
	c.XML(http.StatusOK, err)
}
