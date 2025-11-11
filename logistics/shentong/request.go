package shentong

import "github.com/hzxiao/goutil"

const (
	BaseSandboxUrl = "http://cloudinter-linkgatewaytest.sto.cn/gateway/link.do" //测试
	BaseUrl        = "https://cloudinter-linkgateway.sto.cn/gateway/link.do"    //正式
)

const (
	SUCCESS_TRUE  = "true"
	SUCCESS_FALSE = "false"
)

const (
	OMS_EXPRESS_ORDER_CREATE = "OMS_EXPRESS_ORDER_CREATE" //创建接口
	EDI_MODIFY_ORDER_CANCEL  = "EDI_MODIFY_ORDER_CANCEL"  //取消接口
	GET_ORDERDISPATCH_INFO   = "GET_ORDERDISPATCH_INFO"   //获取取件码接口
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
	Battery     string `json:"battery"`    //带电标识 （10/未知 20/带电 30/不带电）
	GoodsType   string `json:"goodsType"`  //物品类型（大件、小件、扁平件\文件）
	GoodsName   string `json:"goodsName"`  //物品名称
	GoodsCount  int    `json:"goodsCount"` //物品数量
	SpaceX      int    `json:"spaceX"`
	SpaceY      int    `json:"spaceY"`
	SpaceZ      int    `json:"spaceZ"`
	Weight      int    `json:"weight"`
	GoodsAmount string `json:"goodsAmount"`
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
