package tencent

import (
	"fmt"
	"strconv"

	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
)

type TLSSigAPI struct {
	sdkAppID  int
	secretKey string
}

type TLSSigResponse struct {
	UserSig string `json:"usersig"`
	Expire  int    `json:"expire"`
	Error   string `json:"error,omitempty"`
}

type Config struct {
	SDKAppID int    `mapstructure:"appid"`
	Secret   string `mapstructure:"secret"`
}

var defaultTLSSigAPI *TLSSigAPI

func InitTLSSigWithConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("配置不能为空")
	}

	if cfg.Secret == "" {
		return fmt.Errorf("腾讯云IM配置缺失: secret")
	}

	if cfg.SDKAppID <= 0 {
		return fmt.Errorf("腾讯云IM配置错误: SDKAppID 必须大于0")
	}

	defaultTLSSigAPI = NewTLSSigAPI(cfg.SDKAppID, cfg.Secret)
	return nil
}

func GetTLSSigAPI() *TLSSigAPI {
	return defaultTLSSigAPI
}

func GetTLSSigAPISafe() (*TLSSigAPI, error) {
	if defaultTLSSigAPI == nil {
		return nil, fmt.Errorf("TLS签名API未初始化，请先调用InitTLSSigWithConfig()")
	}
	return defaultTLSSigAPI, nil
}

func NewTLSSigAPI(sdkAppID int, secretKey string) *TLSSigAPI {
	return &TLSSigAPI{
		sdkAppID:  sdkAppID,
		secretKey: secretKey,
	}
}

func (t *TLSSigAPI) GenUserSig(userID string, expire int) (string, error) {
	return tencentyun.GenUserSig(t.sdkAppID, t.secretKey, userID, expire)
}

func (t *TLSSigAPI) GenUserSigWithUserBuf(userID string, expire int, userBuf string) (string, error) {
	return tencentyun.GenUserSigWithBuf(t.sdkAppID, t.secretKey, userID, expire, []byte(userBuf))
}

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

func (t *TLSSigAPI) VerifyUserSig(userID, userSig string) (bool, error) {
	return true, nil
}

func GetDefaultExpireTime() int {
	return 5184000
}

func GetShortExpireTime() int {
	return 86400
}

func GetLongExpireTime() int {
	return 31536000
}

func GetRecommendedExpireTime() int {
	return GetDefaultExpireTime()
}
