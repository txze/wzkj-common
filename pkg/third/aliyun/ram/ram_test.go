package ram_test

import (
	"testing"

	"wzkj-common/pkg/third/aliyun/ram"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

func TestRam(t *testing.T) {
	var accessKeyId = ""
	var accessKeySecret = ""
	var endpoint = ""
	var roleArn = ""
	var roleSessionName = ""

	var credentials sts.Credentials
	var err error

	err = ram.Init(endpoint, accessKeyId, accessKeySecret, roleArn, roleSessionName)
	if err != nil {
		t.Fatalf("init err: %v\n", err)
	}
	credentials, err = ram.Ticket()
	if err != nil {
		t.Fatalf("ticket err: %v\n", err)
	}

	/**
	{
		"AccessKeySecret":"2KD2hocsWhH4WqCjq4A8QkawUefAXnyjC1SVfNi7hz1d",
		"Expiration":"2021-11-26T08:34:46Z",
		"AccessKeyId":"STS.NUNsiaCskDUTRXoDw24T5qNRR",
		"SecurityToken":"CAIS7wF1q6Ft5B2yfSjIr5b7ONPVrqxK85e/UH7eoHdnONsZnov5sDz2IHBFf3RvCOobsvk+nm5V7/wclqFsWZxHAEvfdpP8bHvkTETzDbDasumZsJYm6vT8a0XxZjf/2MjNGZabKPrWZvaqbX3diyZ32sGUXD6+XlujQ/br4NwdGbZxZASjaidcD9p7PxZrrNRgVUHcLvGwKBXn8AGyZQhKwlMj0D8vuPrimJHAu0WE0QfAp7VL99irEP+NdNJxOZpzadCx0dFte7DJuCwqsEgbr/op3fYZpWqW5I/AWgMP+XWZIenWt9p0Nx/1rG3j9DCpxxqAAURSPlShRSG+JP0Oty/mrW7c1bSDVF1wuoqTGipuewkZNnnAjLPBzcibNXfbZo/uQWYxmNRtnpIjJoV4ZFK62hpUqmPenDfo8b30ShTlHRzYJA8z3Qjoc7gIEg4mGKmNlg+rt3PvQWE2oA5OX9vswHLv789DUy8bukl7NsIEF0A3"
	}
	*/
	t.Logf("credentials res: %+v\n", credentials)
}
