package tencent_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/txze/wzkj-common/pkg/third/tencent"
	"github.com/txze/wzkj-common/pkg/util"
)

var (
	bucketURL = os.Getenv("bucketURL")
	secretId  = os.Getenv("secretId")
	secretKey = os.Getenv("secretKey")
)

func TestTencentCos(t *testing.T) {
	tencent.Init(bucketURL, secretId, secretKey)

	result, resp, err := tencent.GetClient().Object.Upload(context.Background(), "test_upload.txt", "test_upload.txt", nil)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("result: ", util.S2Json(result))
	fmt.Println("resp: ", util.S2Json(resp))
}
