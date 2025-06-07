package tencent

import (
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentClient struct {
	*cos.Client
}

var defaultClient *TencentClient

func Init(bucketURL, secretId, secretKey string) {
	defaultClient = NewTencentClient(bucketURL, secretId, secretKey)
}

func GetClient() *TencentClient {
	return defaultClient
}

func NewTencentClient(bucketURL, secretId, secretKey string) *TencentClient {
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	return &TencentClient{client}
}
