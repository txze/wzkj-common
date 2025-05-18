package ram

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

type Ram struct {
	*sts.Client
	RoleArn         string
	RoleSessionName string
}

func NewRam() *Ram {
	return &Ram{}
}

var defaultRam *Ram

func GetRam() *Ram { return defaultRam.GetRam() }
func (r *Ram) GetRam() *Ram {
	return defaultRam
}

func Init(endpoint, accessKeyId, accessKeySecret string, roleArn, roleSessionName string) error {
	client, err := sts.NewClientWithAccessKey(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	defaultRam = NewRam()
	defaultRam.RoleArn = roleArn
	defaultRam.RoleSessionName = roleSessionName
	defaultRam.Client = client
	return nil
}

func Ticket() (sts.Credentials, error) { return defaultRam.Ticket() }
func (r *Ram) Ticket() (sts.Credentials, error) {
	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。关于参数含义和设置方法，请参见API参考。
	request.RoleArn = r.RoleArn
	request.RoleSessionName = r.RoleSessionName

	//发起请求，并得到响应。
	response, err := r.AssumeRole(request)
	if err != nil {
		return sts.Credentials{}, err
	}
	return response.Credentials, nil
}
