package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	*oss.Client
	BaseDir    string
	BucketName string
}

var defaultOSS *OSS

func NewOss() *OSS {
	return &OSS{}
}

func GetOSS() *OSS { return defaultOSS.GetOSS() }
func (o *OSS) GetOSS() *OSS {
	return defaultOSS
}

func Init(endpoint, accessKeyId, accessKeySecret string, bucket, baseDir string) error {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	defaultOSS = NewOss()
	defaultOSS.BaseDir = baseDir
	defaultOSS.BucketName = bucket
	defaultOSS.Client = client

	return nil
}

func UploadFile(objectKey string, localFile string) error {
	return defaultOSS.UploadFile(objectKey, localFile)
}
func (o *OSS) UploadFile(objectKey string, localFile string) error {
	bucket, err := defaultOSS.Client.Bucket(o.BucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(objectKey, localFile)
	if err != nil {
		return err
	}

	return nil
}
