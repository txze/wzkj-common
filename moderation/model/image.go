package model

import (
	"encoding/json"
	"fmt"
	"net/http"

	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/txze/wzkj-common/logger"
)

type Image struct {
	client  *green20220302.Client
	content string
	Service string
}

func (i *Image) Invoke() *green20220302.ImageBatchModerationResponseBodyData {
	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			//待检测图片链接，公网可访问的URL。
			"imageUrl": i.content,
			//待检测数据的ID。
			"dataId": uuid.New(),
		},
	)

	request := green20220302.ImageBatchModerationRequest{
		Service:           tea.String(i.Service),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, _err := i.client.ImageBatchModeration(&request)
	if _err != nil {
		logger.Error("识别失败", zap.Error(_err))
		return nil
	}

	if *result.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("response not success. status:%d\n", *result.StatusCode))
		return nil
	}
	body := result.Body
	logger.Info(
		"response success.",
		logger.String("requestId", *body.RequestId),
		logger.Any("code", *body.Code),
	)
	if *body.Code != http.StatusOK {
		logger.Error(fmt.Sprintf("Image moderation not success. code:%d\n", *body.Code))
		return nil
	}

	data := body.Data
	logger.Info("Image moderation data:", logger.Any("data", *data))
	return data
}

func (t *Image) SetClient(client *green20220302.Client) {
	t.client = client
}

func (t *Image) SetContent(content string) {
	t.content = content
}
