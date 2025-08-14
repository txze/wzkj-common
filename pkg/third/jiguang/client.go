package jiguang

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 极光一键登录客户端
type Client struct {
	config     *Config
	httpClient *http.Client
	privateKey *rsa.PrivateKey
}

// LoginTokenRequest 登录Token验证请求
type LoginTokenRequest struct {
	LoginToken string `json:"loginToken"` // 认证SDK获取到的loginToken
	ExID       string `json:"exID"`       // 开发者自定义的id，非必填
}

// LoginTokenResponse 登录Token验证响应
type LoginTokenResponse struct {
	ID      int64  `json:"id"`      // 流水号，请求出错时可能为空
	ExID    string `json:"exID"`    // 开发者自定义的id
	Code    int    `json:"code"`    // 返回码
	Content string `json:"content"` // 返回码说明
	Phone   string `json:"phone"`   // 加密后的手机号码
}

// PhoneInfo 解密后的手机号信息
type PhoneInfo struct {
	Phone   string `json:"phone"`   // 解密后的手机号
	ExID    string `json:"exID"`    // 开发者自定义的id
	Success bool   `json:"success"` // 是否成功
	Error   string `json:"error"`   // 错误信息
}

var defaultClient *Client

// Init 初始化极光客户端（参考其他服务的初始化方式）
func Init() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("加载极光配置失败: %v", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("极光配置验证失败: %v", err)
	}

	client, err := NewClient(config)
	if err != nil {
		return fmt.Errorf("创建极光客户端失败: %v", err)
	}

	defaultClient = client
	return nil
}

// GetClient 获取极光客户端实例
func GetClient() *Client {
	return defaultClient
}

// NewClient 创建新的极光客户端
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 解析RSA私钥
	privateKey, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("解析RSA私钥失败: %v", err)
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		privateKey: privateKey,
	}

	return client, nil
}

// VerifyLoginToken 验证登录Token并获取手机号
func (c *Client) VerifyLoginToken(loginToken, exID string) (*PhoneInfo, error) {
	if loginToken == "" {
		return nil, errors.New("loginToken不能为空")
	}

	// 构建请求
	req := LoginTokenRequest{
		LoginToken: loginToken,
		ExID:       exID,
	}

	// 调用API
	resp, err := c.callAPI(req)
	if err != nil {
		return &PhoneInfo{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// 检查响应码
	if resp.Code != 8000 {
		return &PhoneInfo{
			ExID:    resp.ExID,
			Success: false,
			Error:   fmt.Sprintf("API返回错误: %s (code: %d)", resp.Content, resp.Code),
		}, fmt.Errorf("API返回错误: %s (code: %d)", resp.Content, resp.Code)
	}

	// 解密手机号
	phone, err := c.decryptPhone(resp.Phone)
	if err != nil {
		return &PhoneInfo{
			ExID:    resp.ExID,
			Success: false,
			Error:   fmt.Sprintf("解密手机号失败: %v", err),
		}, fmt.Errorf("解密手机号失败: %v", err)
	}

	return &PhoneInfo{
		Phone:   phone,
		ExID:    resp.ExID,
		Success: true,
	}, nil
}

// callAPI 调用极光API
func (c *Client) callAPI(req LoginTokenRequest) (*LoginTokenResponse, error) {
	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", c.config.APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	// 极光API鉴权方式：Authorization: Basic ${base64(appKey:masterSecret)}
	httpReq.SetBasicAuth(c.config.AppKey, c.config.MasterSecret)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var loginResp LoginTokenResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应体: %s", err, string(respBody))
	}

	return &loginResp, nil
}

// decryptPhone 解密手机号
func (c *Client) decryptPhone(encryptedPhone string) (string, error) {
	if encryptedPhone == "" {
		return "", errors.New("加密的手机号为空")
	}

	// Base64解码
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedPhone)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %v", err)
	}

	// RSA解密
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, c.privateKey, encryptedBytes)
	if err != nil {
		return "", fmt.Errorf("RSA解密失败: %v", err)
	}

	return string(decryptedBytes), nil
}

// parsePrivateKey 解析RSA私钥
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	// 解码PEM格式
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("无法解析PEM格式的私钥")
	}

	// 尝试解析PKCS8格式
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if rsaKey, ok := privateKey.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("私钥不是RSA类型")
	}

	// 尝试解析PKCS1格式
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}

	return rsaKey, nil
}

// 便捷方法：直接验证登录Token（使用默认客户端）
func VerifyLoginToken(loginToken, exID string) (*PhoneInfo, error) {
	if defaultClient == nil {
		return nil, errors.New("极光客户端未初始化，请先调用Init()")
	}
	return defaultClient.VerifyLoginToken(loginToken, exID)
}
