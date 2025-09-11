package moderation

import (
	"testing"

	"github.com/txze/wzkj-common/moderation/model"
)

func TestModerationContext_QueryModeration(t *testing.T) {

	t.Run("TestModerationContext_QueryModerationByNumber", func(t *testing.T) {
		moderation := NewModeration[string](Config{
			AccessKeyId:     "xx",
			AccessKeySecret: "xxx",
			RegionId:        "cn-shanghai",
			Endpoint:        "green-cip.cn-shanghai.aliyuncs.com",
		})
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
