package ierr

// 自定义的错误Code
type Code int

// 定义的错误码段，不同的系统定义不同的错误码段。和电话号码的区号类似
// 如：公共模板的错误为 10，用户模块为11，
// 公共模块的错误码形如10001,10002;用户模块的错误码形如11001,11002;
var CodeNum Code

var (
	Success Code = 0 // 成功

	NotAllowed       Code = 403
	BadRequest       Code = 400
	TokenExpire      Code = 40001
	PermissionDenied Code = 103 // 权限不足
	NotFound         Code = 404

	// 通用错误
	InvalidId              Code = 10001 // 无效的id
	InvalidState           Code = 10002 // 无效的状态
	ParamErr               Code = 10003 // 参数错误
	ReceiveRequestFileFail Code = 10004
	InvalidType            Code = 10005 // 无效的类型
	NonPointError          Code = 10006 // 数据空指针
	ParseDataFail          Code = 10007 // 解码失败
	InternalError          Code = 10008 // 内部错误
	InternalCallRequest    Code = 10009 // 内部调用错误
	ConcurrentError        Code = 10010 // 并发错误
	NetworkBusyError       Code = 10011 // 网络繁忙，请稍后再试
	NetwordError           Code = 10012 // 网络错误
	PermissionDeined       Code = 10013 // 权限不足
	WriteFileError         Code = 10011 // 写入文件错误
	PhoneRepeatError       Code = 10012 // 电话号码重复
	ThirdCallRequest       Code = 10013 // 外部调用错误
	ReadFileError          Code = 10014 // 读取文件错误
	QueueWriteError        Code = 10015 // 写入队列失败
	QueueReadError         Code = 10016 // 读取队列失败
	UnexpectedError        Code = 10017 // 意外的错误

	// 数据库相关
	CreateDataFail    Code = 50001 // 创建数据失败
	DeleteDataFail    Code = 50002 // 删除数据失败
	UpdateDataFail    Code = 50003 // 更新数据失败
	QueryDataFail     Code = 50004 // 查询数据失败
	CommitTxFail      Code = 50005 // 提交事务失败
	CountDataFail     Code = 50006 // 统计数据失败
	RedisOperatorFail Code = 50007 // redis操作数据失败
	UpsertDataFail    Code = 50008 // upsert数据失败
	ScanDataFail      Code = 50009 // scan数据
	GroupByDataFail   Code = 50010 // groupBy数据失败
)
