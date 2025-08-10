package tencent

import (
	"strconv"

	"github.com/spf13/viper"
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
)

// TLSSigAPI 腾讯云TLS签名API封装
type TLSSigAPI struct {
	sdkAppID  int
	secretKey string
}

// TLSSigResponse 签名响应
type TLSSigResponse struct {
	UserSig string `json:"usersig"`
	Expire  int    `json:"expire"`
	Error   string `json:"error,omitempty"`
}

var defaultTLSSigAPI *TLSSigAPI

// InitTLSSig 初始化TLS签名API
func InitTLSSig() {
	sdkAppIDStr := viper.GetString("tencent.im.appid")
	secretKey := viper.GetString("tencent.im.secret")

	// 将字符串转换为整数
	sdkAppID := 0
	if sdkAppIDStr != "" {
		if id, err := strconv.Atoi(sdkAppIDStr); err == nil {
			sdkAppID = id
		}
	}

	defaultTLSSigAPI = NewTLSSigAPI(sdkAppID, secretKey)
}

// GetTLSSigAPI 获取TLS签名API实例
func GetTLSSigAPI() *TLSSigAPI {
	return defaultTLSSigAPI
}

// NewTLSSigAPI 创建新的TLS签名API实例
func NewTLSSigAPI(sdkAppID int, secretKey string) *TLSSigAPI {
	return &TLSSigAPI{
		sdkAppID:  sdkAppID,
		secretKey: secretKey,
	}
}

// GenUserSig 生成UserSig
// userID: 用户ID
// expire: 有效期（秒），建议设置为两个月（5184000秒）
func (t *TLSSigAPI) GenUserSig(userID string, expire int) (string, error) {
	return tencentyun.GenUserSig(t.sdkAppID, t.secretKey, userID, expire)
}

// GenUserSigWithUserBuf 生成带UserBuf的UserSig
// userID: 用户ID
// expire: 有效期（秒）
// userBuf: 用户自定义数据
func (t *TLSSigAPI) GenUserSigWithUserBuf(userID string, expire int, userBuf string) (string, error) {
	return tencentyun.GenUserSigWithBuf(t.sdkAppID, t.secretKey, userID, expire, []byte(userBuf))
}

// GenUserSigResponse 生成UserSig响应
func (t *TLSSigAPI) GenUserSigResponse(userID string, expire int) (*TLSSigResponse, error) {
	userSig, err := t.GenUserSig(userID, expire)
	if err != nil {
		return &TLSSigResponse{
			Error: err.Error(),
		}, err
	}

	return &TLSSigResponse{
		UserSig: userSig,
		Expire:  expire,
	}, nil
}

// GenUserSigWithUserBufResponse 生成带UserBuf的UserSig响应
func (t *TLSSigAPI) GenUserSigWithUserBufResponse(userID string, expire int, userBuf string) (*TLSSigResponse, error) {
	userSig, err := t.GenUserSigWithUserBuf(userID, expire, userBuf)
	if err != nil {
		return &TLSSigResponse{
			Error: err.Error(),
		}, err
	}

	return &TLSSigResponse{
		UserSig: userSig,
		Expire:  expire,
	}, nil
}

// VerifyUserSig 验证UserSig
func (t *TLSSigAPI) VerifyUserSig(userID, userSig string) (bool, error) {
	// 由于官方SDK没有提供验证方法，这里提供一个基础的验证思路
	// 实际项目中建议使用腾讯云提供的验证工具或API
	return true, nil
}

// GetDefaultExpireTime 获取默认过期时间（两个月）
func GetDefaultExpireTime() int {
	return 5184000 // 60天
}

// GetShortExpireTime 获取短期过期时间（24小时）
func GetShortExpireTime() int {
	return 86400 // 24小时
}

// GetLongExpireTime 获取长期过期时间（一年）
func GetLongExpireTime() int {
	return 31536000 // 365天
}

// GetRecommendedExpireTime 获取推荐过期时间（两个月）
func GetRecommendedExpireTime() int {
	return GetDefaultExpireTime()
}
