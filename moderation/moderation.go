package moderation

import (
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
)

type InterfaceModeration[T any] interface {
	Invoke() T
	SetClient(*green20220302.Client)
	SetContent(string)
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

func NewModeration[T any](client *green20220302.Client) *Moderation[T] {
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

func (m *Moderation[T]) Invoke(content string) T {
	m.strategy.SetContent(content)
	return m.strategy.Invoke()
}
