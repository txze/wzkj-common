package rabbitmq

type Message struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}
