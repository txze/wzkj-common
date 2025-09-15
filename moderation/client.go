package moderation

import (
	"context"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/txze/wzkj-common/logger"
)

var ModerationClient *green20220302.Client
var syncOnceModerationClient sync.Once

func NewModerationClient(ctx context.Context, moderationConfig Config) {
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

	syncOnceModerationClient.Do(func() {
		client, err := green20220302.NewClient(config)
		if err != nil {
			logger.FromContext(ctx).Error("内容审核信息初始化失败")
			return
		}
		ModerationClient = client
	})
}
