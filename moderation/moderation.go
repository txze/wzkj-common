package moderation

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/txze/wzkj-common/logger"
)

type InterfaceModeration[T any] interface {
	Invoke()
	SetClient(*green20220302.Client)
	SetContent(T)
}

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
	Endpoint        string
}

type Moderation[T any] struct {
	client   *green20220302.Client
	strategy InterfaceModeration[T]
}

func NewModeration[T any](moderationConfig Config) *Moderation[T] {
	config := &openapi.Config{
		AccessKeyId:     tea.String(moderationConfig.AccessKeyId),
		AccessKeySecret: tea.String(moderationConfig.AccessKeySecret),
		RegionId:        tea.String(moderationConfig.RegionId),
		Endpoint:        tea.String(moderationConfig.Endpoint),
		/**
		 * 请设置超时时间。服务端全链路处理超时时间为10秒，请做相应设置。
		 * 如果您设置的ReadTimeout小于服务端处理的时间，程序中会获得一个ReadTimeout异常。
		 */
		ConnectTimeout: tea.Int(3000),
		ReadTimeout:    tea.Int(6000),
	}
	client, err := green20220302.NewClient(config)
	if err != nil {
		logger.Error("内容审核信息初始化失败")
		return nil
	}

	return &Moderation[T]{
		client: client,
	}
}

func (m *Moderation[T]) SetStrategy(strategy InterfaceModeration[T]) {
	if m.client != nil {
		strategy.SetClient(m.client)
	}

	m.strategy = strategy
}

func (m *Moderation[T]) Invoke(content T) {
	m.strategy.SetContent(content)
	m.strategy.Invoke()
}
