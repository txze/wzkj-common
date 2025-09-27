package allinpay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pay/config"
)

type AllInPay struct {
	config *config.AllInPayConfig
}

type CloseResponse struct {
	CusID     string `json:"cusid" `     // 商户号 - 平台分配的商户号(15位)
	AppID     string `json:"appid" `     // 应用ID - 平台分配的APPID(8位)
	TrxStatus string `json:"trxstatus" ` // 交易状态 - 交易的状态(4位)
	RandomStr string `json:"randomstr"`  // 随机字符串 - 随机生成的字符串(32位)
	ErrMsg    string `json:"errmsg" `    // 错误原因 - 失败的原因说明(100位)
	Sign      string `json:"sign" `      // 签名(32位)
}

type Notify struct {
	AppID       string `json:"appid"`       // 收银宝APPID
	OutTrxID    string `json:"outtrxid"`    // 第三方交易号(暂未启用)
	TrxCode     string `json:"trxcode"`     // 交易类型
	TrxID       string `json:"trxid"`       // 收银宝交易单号
	InitAmt     int64  `json:"initamt"`     // 原始下单金额
	TrxAmt      int64  `json:"trxamt"`      // 交易金额(单位：分)
	TrxDate     string `json:"trxdate"`     // 交易请求日期(yyyymmdd)
	PayTime     string `json:"paytime"`     // 交易完成时间(yyyymmddhhmmss)
	ChnlTrxID   string `json:"chnltrxid"`   // 渠道流水号
	TrxStatus   string `json:"trxstatus"`   // 交易结果码
	CusID       string `json:"cusid"`       // 商户编号
	TermNo      string `json:"termno"`      // 终端编号
	TermBatchID string `json:"termbatchid"` // 终端批次号
	TermTraceNo string `json:"termtraceno"` // 终端流水号
	TermAuthNo  string `json:"termauthno"`  // 终端授权码
	TermRefNum  string `json:"termrefnum"`  // 终端参考号
	TrxReserved string `json:"trxreserved"` // 业务关联内容
	SrcTrxID    string `json:"srctrxid"`    // 原交易流水
	CusOrderID  string `json:"cusorderid"`  // 业务流水(统一下单对应的reqsn订单号)
	Acct        string `json:"acct"`        // 交易账号
	Fee         string `json:"fee"`         // 手续费(单位：分)
	SignType    string `json:"signtype"`    // 签名类型
	CmID        string `json:"cmid"`        // 渠道子商户号
	ChnlID      string `json:"chnlid"`      // 渠道号
	ChnlData    string `json:"chnldata"`    // 渠道信息
	AcctType    string `json:"accttype"`    // 借贷标识
	BankCode    string `json:"bankcode"`    // 发卡行
	LogonID     string `json:"logonid"`     // 支付宝买家账号
	Sign        string `json:"sign"`        // sign校验码
	TlOpenID    string `json:"tlopenid"`    // 通联渠道侧OPENID
}

type PayRequest struct {
	TrxAmt    int    `json:"trxamt"` //单位为分
	Reqsn     string `json:"reqsn"`  //商户订单号
	Validtime string `json:"validtime"`
	NotifyUrl string `json:"notify_url"`
	Body      string `json:"body"` //订单标题
	Remark    string `json:"remark"`
}

type QueryResponse struct {
	// 通信标识
	RetCode string `json:"retcode" xml:"retcode"` // 返回码 SUCCESS/FAIL
	RetMsg  string `json:"retmsg" xml:"retmsg"`   // 返回码说明

	// 业务字段 (仅当RetCode为SUCCESS时返回)
	CusID      string `json:"cusid"`      // 商户号
	AppID      string `json:"appid"`      // 应用ID
	TrxID      string `json:"trxid"`      // 平台交易单号
	ChnlTrxID  string `json:"chnltrxid"`  // 支付渠道交易单号
	ReqSn      string `json:"reqsn"`      // 商户订单号
	TrxCode    string `json:"trxcode"`    // 交易类型
	TrxAmt     int64  `json:"trxamt"`     // 交易金额(分)
	TrxStatus  string `json:"trxstatus"`  // 交易状态
	Acct       string `json:"acct"`       // 支付平台用户标识
	FinTime    string `json:"fintime"`    // 交易完成时间 yyyyMMddHHmmss
	RandomStr  string `json:"randomstr"`  // 随机字符串
	ErrMsg     string `json:"errmsg"`     // 错误原因
	CmID       string `json:"cmid" `      // 渠道子商户号
	ChnlID     string `json:"chnlid"`     // 渠道号
	InitAmt    int64  `json:"initamt"`    // 原交易金额(分)
	Fee        int64  `json:"fee" `       // 手续费(分)
	ChnlData   string `json:"chnldata" `  // 渠道信息
	AcctType   string `json:"accttype" `  // 借贷标识
	BankCode   string `json:"bankcode" `  // 所属银行
	LogonID    string `json:"logonid" `   // 买家账号
	TlOpenID   string `json:"tlopenid" `  // 通联渠道侧OPENID
	TrxReserve string `json:"trxreserve"` // 交易备注
	Sign       string `json:"sign" `      // 签名
}

func NewAllInPay(cfg config.PaymentConfig) *AllInPay {
	return &AllInPay{
		config: cfg.(*config.AllInPayConfig),
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
	params["randomstr"] = common.RandomString32Custom()

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
	params["randomstr"] = common.RandomString32Custom()
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
	var sign string
	for i := 0; i < v.NumField(); i++ {
		if t.Name() == "sign" {
			sign = v.Field(i).String()
			continue
		}
		checkSign[t.Name()] = v.Field(i).Interface()
	}

	isCheck, err := a.VerifySign(checkSign, sign)
	if err != nil {
		return nil, err
	}

	if isCheck == false {
		return nil, errors.New("sign is invalid")
	}

	return &common.UnifiedResponse{
		Platform:   a.GetType(),
		OrderID:    rsp.ReqSn,
		PlatformID: rsp.ChnlTrxID,
		Amount:     rsp.InitAmt,
		Status:     rsp.TrxStatus,
		PaidAmount: rsp.TrxAmt,
		PaidTime:   rsp.FinTime,
		Message:    rsp,
	}, nil
}

func (a *AllInPay) Refund(orderID string, amount float64) error {
	//TODO implement me
	panic("implement me")
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

func (a *AllInPay) VerifySign(params map[string]interface{}, sign string) (bool, error) {
	// 获取签名并从参数中移除
	if sign == "" {
		return false, errors.New("sign is required")
	}

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

func (a *AllInPay) VerifyNotification(notification *Notify) (*common.UnifiedResponse, error) {
	var params = make(map[string]interface{})
	v := reflect.ValueOf(notification)
	t := reflect.TypeOf(notification)
	var sign string
	for i := 0; i < v.NumField(); i++ {
		if t.Name() == "sign" {
			sign = v.Field(i).String()
			continue
		}
		params[t.Name()] = v.Field(i).Interface()
	}

	isCheck, err := a.VerifySign(params, sign)
	if err != nil {
		return nil, err
	}

	if isCheck == false {
		return nil, errors.New("sign is invalid")
	}

	return &common.UnifiedResponse{
		Platform:   a.GetType(),
		OrderID:    notification.CusOrderID,
		PlatformID: notification.ChnlTrxID,
		Amount:     notification.InitAmt,
		Status:     notification.TrxStatus,
		PaidAmount: notification.TrxAmt,
		PaidTime:   notification.PayTime,
		Message:    notification,
	}, nil
}

func (a *AllInPay) Close(orderId string) (bool, error) {
	params := make(map[string]interface{})
	params["cusid"] = a.config.CuSID
	params["appid"] = a.config.AppId
	params["version"] = 12
	params["oldreqsn"] = orderId
	params["signtype"] = "RSA"
	params["randomstr"] = common.RandomString32Custom()
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
	var sign string
	for i := 0; i < v.NumField(); i++ {
		if t.Name() == "sign" {
			sign = v.Field(i).String()
			continue
		}
		checkSign[t.Name()] = v.Field(i).Interface()
	}

	isCheck, err := a.VerifySign(checkSign, sign)
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
