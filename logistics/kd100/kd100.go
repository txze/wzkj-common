package kd100

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type KD100Config struct {
	KEY      string
	CUSTOMER string
}

type KD100 struct {
	config KD100Config
}

func NewKD100(cfg KD100Config) *KD100 {
	return &KD100{config: cfg}
}

func (K *KD100) QueryLogisticsByNumber(code, number, phone string) (string, error) {
	param := map[string]string{
		"com":   code,
		"num":   number,
		"phone": phone,
	}

	// 将参数转换为JSON字符串
	paramJson, _ := json.Marshal(param)
	paramStr := string(paramJson)

	// 发送请求
	return K.customerRequest(paramStr, QUERY_URL)
}

// customerRequest 鉴权
func (K *KD100) customerRequest(param string, postUrl string) (string, error) {
	// 计算签名
	signStr := param + K.config.KEY + K.config.CUSTOMER
	hash := md5.New()
	hash.Write([]byte(signStr))
	sign := hex.EncodeToString(hash.Sum(nil))
	sign = strings.ToUpper(sign)

	// 构造form表单数据
	formData := url.Values{}
	formData.Add("customer", K.config.CUSTOMER)
	formData.Add("sign", sign)
	formData.Add("param", param)

	return execute(postUrl, formData)
}

/**
*执行HTTP请求
 */
func execute(postUrl string, formData url.Values) (string, error) {
	// 创建HTTP客户端
	client := &http.Client{}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", postUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(respBody), err
	}

	// 打印响应内容
	return string(respBody), err
}
