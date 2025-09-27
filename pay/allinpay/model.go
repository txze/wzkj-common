package allinpay

// 关闭订单接口响应
type CloseResponse struct {
	// 通信标识
	RetCode string `json:"retcode" xml:"retcode"` // 返回码 SUCCESS/FAIL
	RetMsg  string `json:"retmsg" xml:"retmsg"`   // 返回码说明

	CusID     string `json:"cusid" `     // 商户号 - 平台分配的商户号(15位)
	AppID     string `json:"appid" `     // 应用ID - 平台分配的APPID(8位)
	TrxStatus string `json:"trxstatus" ` // 交易状态 - 交易的状态(4位)
	RandomStr string `json:"randomstr"`  // 随机字符串 - 随机生成的字符串(32位)
	ErrMsg    string `json:"errmsg" `    // 错误原因 - 失败的原因说明(100位)
	Sign      string `json:"sign" `      // 签名(32位)
}

// 回掉信息结构
type Notify struct {
	AppID       string `json:"appid"`       // 收银宝APPID
	OutTrxID    string `json:"outtrxid"`    // 第三方交易号(暂未启用)
	TrxCode     string `json:"trxcode"`     // 交易类型
	TrxID       string `json:"trxid"`       // 收银宝交易单号
	InitAmt     int    `json:"initamt"`     // 原始下单金额
	TrxAmt      int    `json:"trxamt"`      // 交易金额(单位：分)
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

// 查询订单响应
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
	TrxAmt     int    `json:"trxamt"`     // 交易金额(分)
	TrxStatus  string `json:"trxstatus"`  // 交易状态
	Acct       string `json:"acct"`       // 支付平台用户标识
	FinTime    string `json:"fintime"`    // 交易完成时间 yyyyMMddHHmmss
	RandomStr  string `json:"randomstr"`  // 随机字符串
	ErrMsg     string `json:"errmsg"`     // 错误原因
	CmID       string `json:"cmid" `      // 渠道子商户号
	ChnlID     string `json:"chnlid"`     // 渠道号
	InitAmt    int    `json:"initamt"`    // 原交易金额(分)
	Fee        int    `json:"fee" `       // 手续费(分)
	ChnlData   string `json:"chnldata" `  // 渠道信息
	AcctType   string `json:"accttype" `  // 借贷标识
	BankCode   string `json:"bankcode" `  // 所属银行
	LogonID    string `json:"logonid" `   // 买家账号
	TlOpenID   string `json:"tlopenid" `  // 通联渠道侧OPENID
	TrxReserve string `json:"trxreserve"` // 交易备注
	Sign       string `json:"sign" `      // 签名
}

// 生成支付信息
type PayRequest struct {
	TrxAmt    int    `json:"trxamt"` //单位为分
	Reqsn     string `json:"reqsn"`  //商户订单号
	Validtime string `json:"validtime"`
	NotifyUrl string `json:"notify_url"`
	Body      string `json:"body"` //订单标题
	Remark    string `json:"remark"`
}
