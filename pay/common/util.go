package common

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sort"
	"strings"
)

// ToUrlParams 辅助函数，用于拼接URL参数（与PHP中的AppUtil::ToUrlParams对应）
func ToUrlParams(params map[string]interface{}) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprintf("%v", params[k]))
	}
	return buf.String()
}

func RandomString32Custom() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 32)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < 32; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return ""
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}

/**
*执行HTTP请求
 */
func Execute(postUrl string, params map[string]interface{}) (string, error) {
	// 创建HTTP客户端
	client := &http.Client{}
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(params)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", postUrl, bytes.NewReader(jsonByte))
	if err != nil {
		return "", err
	}

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
