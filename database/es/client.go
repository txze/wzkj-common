package es

import (
	"net/http"
	"sync"

	"github.com/elastic/go-elasticsearch/v9"
)

type EsTypedClient struct {
	es   *elasticsearch.TypedClient
	once sync.Once
}

var esClient *EsTypedClient = &EsTypedClient{}

// New 创建EsTypedClient实例
// addresses: Elasticsearch节点地址列表
// username: Elasticsearch用户名
// password: Elasticsearch密码
// transport: HTTP传输配置
func New(addresses []string, username, password string, transport *http.Transport) *EsTypedClient {
	if len(addresses) == 0 {
		panic("invalid elasticsearch config: addresses is empty")
	}
	if username == "" {
		panic("invalid elasticsearch config: username is empty")
	}
	if password == "" {
		panic("invalid elasticsearch config: password is empty")
	}
	// 使用单例模式创建EsTypedClient实例
	if transport == nil {
		transport = &http.Transport{}
	}
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
		Transport: transport,
	}
	var err error
	esClient.once.Do(func() {
		esClient.es, err = elasticsearch.NewTypedClient(cfg)
		if err != nil {
			panic(err)
		}
	})
	return esClient
}

func (c *EsTypedClient) Client() *elasticsearch.TypedClient {
	return c.es
}

func Client() *elasticsearch.TypedClient {
	return esClient.es
}
