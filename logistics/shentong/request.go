package shentong

import (
	"crypto/md5"
	"encoding/base64"
	"net/url"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/pkg/util"
)

const (
	BaseSandboxUrl = "http://cloudinter-linkgatewaytest.sto.cn/gateway/link.do" //测试
	BaseUrl        = "https://cloudinter-linkgateway.sto.cn/gateway/link.do"    //正式
)

const (
	SUCCESS_TRUE  = "true"
	SUCCESS_FALSE = "false"
)

const (
	OMS_EXPRESS_ORDER_CREATE     = "OMS_EXPRESS_ORDER_CREATE"     //创建接口
	EDI_MODIFY_ORDER_CANCEL      = "EDI_MODIFY_ORDER_CANCEL"      //取消接口
	GET_ORDERDISPATCH_INFO       = "GET_ORDERDISPATCH_INFO"       //获取取件码接口
	QUERY_SEND_SERVICE_DETAIL    = "QUERY_SEND_SERVICE_DETAIL"    //查询时效运费
	PERSONAL_ADDRESS_PARSE       = "PERSONAL_ADDRESS_PARSE"       //地址解析
	STO_TRACE_PLATFORM_SUBSCRIBE = "STO_TRACE_PLATFORM_SUBSCRIBE" //物流详情订阅(CP请求STO)
)

type CommonParamsRequest struct {
	ApiName    string `json:"api_name"`
	Content    string `json:"content"`     //结构化的业务报文体，可以是JSON或XML格式的字串（见下文表格及示例）
	FromAppKey string `json:"from_appkey"` //订阅方/请求发起方的应用key
	FromCode   string `json:"from_code"`   //订阅方/请求发起方的应用资源code
	ToAppkey   string `json:"to_appkey"`
	ToCode     string `json:"to_code"`
	DataDigest string `json:"data_digest"` //报文签名
}

type CreateOrderRequest struct {
	OrderNo             string     `json:"orderNo"`     //订单号（客户系统自己生成，唯一）
	OrderSource         string     `json:"orderSource"` //订单来源（订阅服务时填写的来源编码）
	BillType            string     `json:"billType"`
	OrderType           string     `json:"orderType"`
	Sender              Sender     `json:"sender"`
	Receiver            Receiver   `json:"receiver"` //收件人信息
	Cargo               Cargo      `json:"cargo"`    //包裹信息
	Customer            Customer   `json:"customer"` //客户信息，在线下单取运单号必填，代单号下单不需要填写，测试账号传值如下，生产账号联系合作业务方提供
	WaybillNo           string     `json:"waybillNo"`
	CodValue            string     `json:"codValue"`
	FreightCollectValue string     `json:"freightCollectValue"`
	TimelessType        string     `json:"timelessType"`
	ProductType         string     `json:"productType"`
	ExtendFieldMap      goutil.Map `json:"extendFieldMap"`
	Remark              string     `json:"remark"`
	ExpressDirection    string     `json:"expressDirection"`
	CreateChannel       string     `json:"createChannel"`
	RegionType          string     `json:"regionType"`
	ExpectValue         string     `json:"expectValue"`
	PayModel            string     `json:"payModel"`
}

type Sender struct {
	Name     string `json:"name"`     //寄件人名称
	Tel      string `json:"tel"`      //寄件人固定电话
	Mobile   string `json:"mobile"`   //寄件人手机号码
	PostCode string `json:"postCode"` //邮编
	Country  string `json:"country"`  //国家
	Province string `json:"province"` //省
	City     string `json:"city"`     //市
	Area     string `json:"area"`     //区
	Town     string `json:"town"`     //镇
	Address  string `json:"address"`  //详细地址
}

type Receiver struct {
	Name     string `json:"name"`   //收件人名称
	Tel      string `json:"tel"`    //收件人固定电话
	Mobile   string `json:"mobile"` //收件人手机号码
	PostCode string `json:"postCode"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	Town     string `json:"town"`
	Address  string `json:"address"`
	SafeNo   string `json:"safeNo"`
}

type Cargo struct {
	Battery    string `json:"battery"`    //带电标识 （10/未知 20/带电 30/不带电）
	GoodsType  string `json:"goodsType"`  //物品类型（大件、小件、扁平件\文件）
	GoodsName  string `json:"goodsName"`  //物品名称
	GoodsCount int    `json:"goodsCount"` //物品数量
	Weight     int    `json:"weight"`     //kg
}

type Customer struct {
	SiteCode          string `json:"siteCode"`
	CustomerName      string `json:"customerName"`
	SitePwd           string `json:"sitePwd"`
	MonthCustomerCode string `json:"monthCustomerCode"`
}

type CancelOrderRequest struct {
	BillCode    string `json:"billCode"`    //订单号（客户系统自己生成，唯一）
	OrderType   string `json:"orderType"`   //01：普通订单，02：调动订单，该类型传值同下单接口一致
	OrderSource string `json:"orderSource"` //订单来源编码
}

type PickOrderInfoRequest struct {
	BillCode        string `json:"BillCode"`            //订单号（客户系统自己生成，唯一）
	SourceOrderId   string `json:"SourceOrderId"`       //第三方平台订单号
	OrderSourceCode string `json:"OrderSourceCodeCode"` //订单来源编码
}

type QuerySendServiceDetailRequest struct {
	SendName    string `json:"SendName"`    //发件人姓名
	SendMobile  string `json:"SendMobile"`  //发件人手机号
	SendProv    string `json:"SendProv"`    //发件人省份
	SendCity    string `json:"SendCity"`    //发件人城市
	SendArea    string `json:"SendArea"`    //发件人区县
	SendAddress string `json:"SendAddress"` //发件人详细地址
	RecName     string `json:"RecName"`     //收件人姓名
	RecMobile   string `json:"RecMobile"`   //收件人手机号
	RecProv     string `json:"RecProv"`     //收件人省份
	RecCity     string `json:"RecCity"`     //收件人城市
	RecArea     string `json:"RecArea"`     //收件人区县
	RecAddress  string `json:"RecAddress"`  //收件人详细地址
	OpenId      string `json:"OpenId"`      //下单用户唯一标识
	Weight      string `json:"Weight"`      //物品重量（kg）
}

type ParseAddressRequest struct {
	AddressText string `json:"addressText"`
}

func convertFormData(apiName, appKey, fromCode, toAppKey, secretKey string, data interface{}) url.Values {
	// 序列化业务数据
	dataStr := util.S2Json(data)
	// 生成签名
	sign := generateSign(dataStr, secretKey)

	// 构建参数
	// 构造表单数据
	formData := url.Values{}
	formData.Add("api_name", apiName)
	formData.Add("content", dataStr)
	formData.Add("from_appkey", appKey)
	formData.Add("from_code", fromCode)
	formData.Add("to_appkey", toAppKey)
	formData.Add("to_code", toAppKey)
	formData.Add("data_digest", sign)

	return formData
}

// generateSign 生成签名
func generateSign(content, secretKey string) string {
	// 拼接内容和密钥
	data := content + secretKey

	// 计算MD5哈希（二进制格式）
	hash := md5.Sum([]byte(data))

	// Base64编码
	return base64.StdEncoding.EncodeToString(hash[:])
}
