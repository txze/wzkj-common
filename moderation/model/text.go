package model

import (
	"encoding/json"
	"fmt"
	"net/http"

	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"

	"github.com/txze/wzkj-common/logger"
)

type Text struct {
	client  *green20220302.Client
	content string
}

func (t *Text) Invoke() {
	if t.content == "" {
		return
	}
	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": t.content,
		},
	)

	request := green20220302.TextModerationPlusRequest{
		Service:           tea.String("aigc_moderation_byllm"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, _err := t.client.TextModerationPlus(&request)
	if _err != nil {
		logger.Error("识别失败", zap.Error(_err))
		return
	}

	if *result.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("response not success. status:%d\n", *result.StatusCode))
		return
	}
	body := result.Body
	logger.Info(
		"response success.",
		logger.String("requestId", *body.RequestId),
		logger.Any("code", *body.Code),
		logger.String("msg", *body.Message),
	)
	if *body.Code != http.StatusOK {
		logger.Error(fmt.Sprintf("text moderation not success. code:%d\n", *body.Code))
		return
	}

	data := body.Data
	logger.Info("text moderation data:", logger.Any("data", *data))
}

func (t *Text) SetClient(client *green20220302.Client) {
	t.client = client
}

func (t *Text) SetContent(content string) {
	t.content = content
}
