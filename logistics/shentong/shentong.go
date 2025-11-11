package shentong

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hzxiao/goutil"
	"github.com/mitchellh/mapstructure"

	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

type STOClient struct {
	cfg *Config
}

func NewSTOClient(cfg *Config) *STOClient {
	return &STOClient{
		cfg: cfg,
	}
}

// generateSign 生成签名
func (c *STOClient) generateSign(content string) string {
	// 拼接内容和密钥
	data := content + c.cfg.SecretKey

	// 计算MD5哈希（二进制格式）
	hash := md5.Sum([]byte(data))

	// Base64编码
	return base64.StdEncoding.EncodeToString(hash[:])
}

// CreateOrder 创建订单
func (c *STOClient) CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
	req.OrderSource = c.cfg.SourceCode
	baseResp, err := c.doRequest(OMS_EXPRESS_ORDER_CREATE, "sto_oms", req)
	if err != nil {
		return nil, err
	}
	if baseResp.GetString("success") == SUCCESS_FALSE {
		return nil, ierr.NewIError(ierr.ParamErr, fmt.Sprintf("API错误: %s(%s)", baseResp.Get("errorMsg"), baseResp.Get("errorCode")))
	}

	var rsp *CreateOrderResponse
	err = mapstructure.Decode(baseResp.GetMapP("data"), &rsp)
	if err != nil {
		return nil, ierr.NewIErrorf(ierr.ParseDataFail, err.Error())
	}

	return rsp, nil
}

// CancelOrder 取消订单
func (c *STOClient) CancelOrder(req *CancelOrderRequest) (*PickInfoResponse, error) {
	req.OrderSource = c.cfg.SourceCode
	baseResp, err := c.doRequest(GET_ORDERDISPATCH_INFO, "ORDERMS_API", req)
	if err != nil {
		return nil, err
	}
	if baseResp.GetString("success") == SUCCESS_FALSE {
		return nil, ierr.NewIErrorf(ierr.ParamErr, "API错误: %s(%s)", baseResp.Get("errorMsg"), baseResp.Get("errorCode"))
	}

	var rsp *PickInfoResponse
	err = mapstructure.Decode(baseResp.GetMapP("data"), &rsp)
	if err != nil {
		return nil, ierr.NewIErrorf(ierr.ParseDataFail, err.Error())
	}

	return rsp, nil
}

// PickOrderInfo 获取取件信息
func (c *STOClient) PickOrderInfo(req *PickOrderInfoRequest) error {
	req.OrderSourceCode = c.cfg.SourceCode
	baseResp, err := c.doRequest(EDI_MODIFY_ORDER_CANCEL, "edi_modify_order", req)
	if err != nil {
		return err
	}
	if baseResp.GetString("success") == SUCCESS_FALSE {
		return ierr.NewIErrorf(ierr.ParamErr, "API错误: %s(%s)", baseResp.Get("errorMsg"), baseResp.Get("errorCode"))
	}

	return nil
}

// doRequest 执行API请求
func (c *STOClient) doRequest(apiName, toAppKey string, data interface{}) (goutil.Map, error) {
	// 序列化业务数据
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, ierr.NewIErrorf(ierr.ParseDataFail, "序列化数据失败: %w", err)
	}
	dataStr := string(dataBytes)
	// 生成签名
	sign := c.generateSign(dataStr)

	// 构建参数
	// 构造表单数据
	formData := url.Values{}
	formData.Add("api_name", apiName)
	formData.Add("content", dataStr)
	formData.Add("from_appkey", c.cfg.AppKey)
	formData.Add("from_code", c.cfg.ResourceCode)
	formData.Add("to_appkey", toAppKey)
	formData.Add("to_code", toAppKey)
	formData.Add("data_digest", sign)

	resp, err := util.HttpFormDataPost(c.cfg.GetBaseUrl(), formData)
	if err != nil {
		return nil, ierr.NewIErrorf(ierr.InternalError, "请求失败: %w", err)
	}

	return resp, nil
}
