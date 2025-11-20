package kd100

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

type KD100 struct {
	config       KD100Config
	adaptorQuery *QueryLogisticsAdaptor
}

func (k *KD100) GetPriceQuote(req *model.GetPriceQuoteReq) (*model.PriceQuote, error) {
	//TODO implement me
	panic("implement me")
}

func (K *KD100) QueryLogistics(req *model.QueryLogisticsRequest) (*model.QueryResp, error) {
	param := K.adaptorQuery.ConvertRequest(req)
	signStr := param + K.config.KEY + K.config.CUSTOMER
	sign := K.generateSign(signStr)
	// 构造form表单数据
	formData := url.Values{}
	formData.Add("customer", K.config.CUSTOMER)
	formData.Add("sign", sign)
	formData.Add("param", param)
	rsp, err := model.DoRequest(QueryURL, formData)
	if err != nil {
		return nil, err
	}
	return K.adaptorQuery.ParseResponse(rsp)
}

func (K *KD100) CreateOrder(req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	//TODO implement me
	panic("implement me")
}

func (K *KD100) CancelOrder(req *model.CancelOrderReq) error {
	//TODO implement me
	panic("implement me")
}

func (K *KD100) ParseWebhook(body []byte) (*model.WebhookData, error) {
	//TODO implement me
	panic("implement me")
}

func (K *KD100) ParseOrderNotify(body []byte) (*model.OrderNotifyResp, error) {
	//TODO implement me
	panic("implement me")
}

func (K *KD100) ParseAddress(addr string) (model.Address, error) {
	param := map[string]string{
		"content": addr,
	}
	// 将参数转换为JSON字符串
	paramStr := util.S2Json(param)
	// 拼接签名参数
	t := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signStr := paramStr + t + K.config.KEY + K.config.Secret
	sign := K.generateSign(signStr)
	formData := url.Values{}
	formData.Add("key", K.config.KEY)
	formData.Add("t", t)
	formData.Add("sign", sign)
	formData.Add("param", paramStr)
	// 发送请求
	data, err := model.DoRequest(ParseAddressURL, formData)

	var result ParseAddress
	err = util.S2S(data, &result)
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

func NewKD100(cfg KD100Config) *KD100 {
	return &KD100{
		config:       cfg,
		adaptorQuery: &QueryLogisticsAdaptor{},
	}
}
