package moderation

import (
	"context"
	"testing"

	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/moderation/model"
)

func TestModerationContext_QueryModeration(t *testing.T) {

	t.Run("TestModerationContext_QueryModerationByNumber", func(t *testing.T) {
		logger.Init("./logs")
		NewModerationClient(context.Background(), Config{
			AccessKeyId:     "xxxxxx",
			AccessKeySecret: "xxxxxx",
			RegionId:        "cn-shanghai",
			Endpoint:        "green-cip.cn-shanghai.aliyuncs.com",
		})
		if ModerationClient == nil {
			t.Errorf("TestModerationContext_QueryModerationByNumber() error")
			return
		}

		moderation := NewModeration[*green20220302.ImageBatchModerationResponseBodyData](ModerationClient)
		moderation.SetStrategy(&model.Image{})
		moderation.Invoke(
			"https://gunpladata.oss-cn-guangzhou.aliyuncs.com/yourBasePath/uploads/2025-08-01/1757513458699.jpg",
		)
		//if err != nil {
		//	t.Errorf("QueryModerationByNumber() error = %v", err)
		//	return
		//}
		//t.Log(got)
	})
}
