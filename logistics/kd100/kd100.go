package kd100

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

type KD100 struct {
	config KD100Config
}

func NewKD100(cfg KD100Config) *KD100 {
	return &KD100{config: cfg}
}

func (K *KD100) QueryLogisticsByNumber(code, number, phone, resultv2 string) (string, error) {
	param := map[string]string{
		"com":      code,
		"num":      number,
		"phone":    phone,
		"resultv2": resultv2,
	}

	// 将参数转换为JSON字符串
	paramJson, _ := json.Marshal(param)
	paramStr := string(paramJson)

	signStr := paramStr + K.config.KEY + K.config.CUSTOMER
	sign := K.generateSign(signStr)
	// 构造form表单数据
	formData := url.Values{}
	formData.Add("customer", K.config.CUSTOMER)
	formData.Add("sign", sign)
	formData.Add("param", paramStr)

	// 发送请求
	return doRequest(QueryURL, formData)
}

func (K *KD100) ParseAddress(addr string) (model.Address, error) {
	param := map[string]string{
		"content": addr,
	}

	// 将参数转换为JSON字符串
	paramJson, _ := json.Marshal(param)
	paramStr := string(paramJson)

	// 拼接签名参数
	t := fmt.Sprintf("%d", time.Now().UnixMilli())
	signStr := paramStr + t + K.config.KEY + K.config.Secret
	sign := K.generateSign(signStr)
	formData := url.Values{}
	formData.Add("key", K.config.KEY)
	formData.Add("t", t)
	formData.Add("sign", sign)
	formData.Add("param", paramStr)
	// 发送请求
	data, err := doRequest(ParseAddressURL, formData)

	var result ParseAddress
	err = util.Json2S(data, &result)
	if err != nil {
		return nil, err
	}

	if result.Code == http.StatusOK {
		city := result.Data.Result[0]
		return city.Xzq, nil
	}

	return nil, ierr.NewIError(ierr.InternalError, result.Message)
}

func (k *KD100) generateSign(signStr string) string {
	hash := md5.New()
	hash.Write([]byte(signStr))
	sign := hex.EncodeToString(hash.Sum(nil))
	sign = strings.ToUpper(sign)
	return sign
}

/**
*执行HTTP请求
 */
func doRequest(postUrl string, formData url.Values) (string, error) {
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
