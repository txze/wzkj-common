package allinpay

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-pay/gopay"
	"github.com/go-pay/util"
	"github.com/pkg/errors"

	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pay/config"
)

type AllInPay struct {
	config *AllInPayConfig
}

func NewAllInPay(cfg config.PaymentConfig) *AllInPay {
	return &AllInPay{
		config: cfg.(*AllInPayConfig),
	}
}

func (a *AllInPay) Pay(request *PayRequest) (map[string]interface{}, error) {
	//TODO implement me
	params := make(map[string]interface{})
	params["cusid"] = a.config.CuSID
	params["appid"] = a.config.AppId
	params["version"] = 12
	params["trxamt"] = request.TrxAmt
	params["reqsn"] = request.Reqsn
	params["validtime"] = request.Validtime
	params["notify_url"] = request.NotifyUrl
	params["body"] = request.Body
	params["remark"] = request.Remark
	params["paytype"] = "A02"
	params["randomstr"] = util.RandomString(32)

	sign, err := a.GenerateSign(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = sign

	return params, nil
}

func (a *AllInPay) QueryPayment(orderID string) (*common.UnifiedResponse, error) {
	params := make(map[string]interface{})
	params["cusid"] = a.config.CuSID
	params["appid"] = a.config.AppId
	params["version"] = 12
	params["reqsn"] = orderID
	params["signtype"] = "RSA"
	params["randomstr"] = util.RandomString(32)
	params["sign"], _ = a.GenerateSign(params)
	rspStr, err := common.Execute(a.config.QueryOrderUrl, params)
	if err != nil {
		return nil, err
	}
	var rsp *QueryResponse
	_ = json.Unmarshal([]byte(rspStr), &rsp)
	var checkSign = make(map[string]interface{})
	v := reflect.ValueOf(rsp)
	t := reflect.TypeOf(rsp)
	for i := 0; i < v.NumField(); i++ {
		checkSign[t.Name()] = v.Field(i).Interface()
	}

	isCheck, err := a.VerifySign(checkSign)
	if err != nil {
		return nil, err
	}

	if isCheck == false {
		return nil, errors.New("sign is invalid")
	}
	var status bool
	if gopay.SUCCESS == rsp.TrxStatus {
		status = true
	}
	return &common.UnifiedResponse{
		Platform:    a.GetType(),
		OrderID:     rsp.ReqSn,
		PlatformID:  rsp.ChnlTrxID,
		Amount:      rsp.InitAmt,
		Status:      status,
		TradeStatus: rsp.TrxStatus,
		PaidAmount:  rsp.TrxAmt,
		//PaidTime:    rsp.FinTime,
		Message: rsp,
	}, nil
}

func (a *AllInPay) Refund(ctx context.Context, request *common.RefundRequest) error {
	return nil
}

func (a *AllInPay) GenerateSign(params map[string]interface{}) (string, error) {
	//1. 拼接参数字符串
	bufSignSrc := common.ToUrlParams(params)
	//2. 处理私钥字符串格式
	block, _ := pem.Decode([]byte(a.config.PrivateKey))

	if block == nil {
		return "", errors.Errorf("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 尝试PKCS8格式
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", errors.Errorf("failed to parse private key: %v", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return "", errors.Errorf("not an RSA private key")
		}
	}

	// 3. 使用SHA256进行RSA签名
	hashed := sha1.Sum([]byte(bufSignSrc))
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", errors.Errorf("sign failed: %v", err)
	}

	// 4. Base64编码签名结果
	sign := base64.StdEncoding.EncodeToString(signature)
	return sign, nil
}

func (a *AllInPay) VerifySign(params map[string]interface{}) (bool, error) {
	sign := params["sign"].(string)
	// 获取签名并从参数中移除
	if sign == "" {
		return false, errors.New("sign is required")
	}
	delete(params, "sign")

	bufSignSrc := common.ToUrlParams(params)

	// RSA 公钥
	publicKey := a.config.PublicKey
	// 解码签名
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, errors.Errorf("failed to decode sign: %v", err)
	}

	// 解析公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, fmt.Errorf("failed to parse public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("not an RSA public key")
	}

	// 计算数据的哈希
	hash := crypto.SHA1.New()
	hash.Write([]byte(bufSignSrc))
	hashed := hash.Sum(nil)

	// 验证签名
	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA1, hashed, signBytes)
	if err != nil {
		return false, nil // 签名验证失败，但不返回错误
	}

	return true, nil
}

func (a *AllInPay) VerifyNotification(req *http.Request) (*common.UnifiedResponse, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	var form map[string][]string = req.Form
	params := make(map[string]interface{}, len(form)+1)
	for k, v := range form {
		if len(v) == 1 {
			params[k] = v[0]
		}
	}

	isCheck, err := a.VerifySign(params)
	if err != nil {
		return nil, err
	}

	if isCheck == false {
		return nil, errors.New("sign is invalid")
	}

	return &common.UnifiedResponse{
		Platform:   a.GetType(),
		OrderID:    params["cusorderid"].(string),
		PlatformID: params["chnltrxid"].(string),
		//Amount:     params["initamt"].(int),
		//Status:     params["trxstatus"].(string),
		//PaidAmount: params["trxamt"].(int),
		//PaidTime: params["paytime"].(string),
		Message: params,
	}, nil
}

func (a *AllInPay) Close(orderId string) (bool, error) {
	params := make(map[string]interface{})
	params["cusid"] = a.config.CuSID
	params["appid"] = a.config.AppId
	params["version"] = 12
	params["oldreqsn"] = orderId
	params["signtype"] = "RSA"
	params["randomstr"] = util.RandomString(32)
	params["sign"], _ = a.GenerateSign(params)
	rspStr, err := common.Execute(a.config.CloseOrderUrl, params)
	if err != nil {
		return false, err
	}

	var rsp *CloseResponse
	_ = json.Unmarshal([]byte(rspStr), &rsp)
	var checkSign = make(map[string]interface{})
	v := reflect.ValueOf(rsp)
	t := reflect.TypeOf(rsp)
	for i := 0; i < v.NumField(); i++ {
		checkSign[t.Name()] = v.Field(i).Interface()
	}

	isCheck, err := a.VerifySign(checkSign)
	if err != nil {
		return false, err
	}

	if isCheck == false {
		return false, errors.New("sign is invalid")
	}

	return rsp.TrxStatus == "0000", nil
}
func (a *AllInPay) GetType() string {
	return a.config.GetType()
}
