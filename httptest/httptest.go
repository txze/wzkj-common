package httptest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Token   string
	BaseUrl string
	Router  *gin.Engine
}

func (c *Client) Get(path string, params map[string]interface{}) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(params)

	// 构造post请求，json数据以请求body的形式传递
	var uri = c.BaseUrl + path
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Authorization", c.Token)

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	c.Router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

func (c *Client) Post(path string, data map[string]interface{}) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(data)

	var uri = c.BaseUrl + path
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Authorization", c.Token)

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	c.Router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

func (c *Client) Put(path string, data map[string]interface{}) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(data)

	var uri = c.BaseUrl + path
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("PUT", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Authorization", c.Token)

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	c.Router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

func (c *Client) Delete(path string, data map[string]interface{}) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(data)

	var uri = c.BaseUrl + path
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("DELETE", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Authorization", c.Token)

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	c.Router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}
