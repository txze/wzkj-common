package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Entry struct {
	ctx            context.Context
	mq             *RabbitMQ
	exchangeOption *ExchangeOption
	queueOption    *QueueOption
	consumeOption  *ConsumeOption
	publishOption  *PublishOption
	channel        *amqp.Channel
	channelCloseCh chan *amqp.Error
	headers        amqp.Table
	queueName      string // 队列名称
	exchangeName   string // 交换机
	routingName    string // 路由值
}

// ConsumeOption 消费者参数
type ConsumeOption struct {
	QueueName  string
	Consumer   string
	AutoAck    bool
	Exclusive  bool
	NoLocal    bool
	NoWait     bool
	Args       amqp.Table
	Workers    int
	RetryCount int
}

// PublishOption 生产者参数
type PublishOption struct {
	Priority   int
	Expiration string
	MessageId  string
}

type IEntry interface {
	WithContext(ctx context.Context) IEntry
	Channel(ch *amqp.Channel) IEntry
	ChannelCloseCh() <-chan *amqp.Error
	Close()
	Exchange(data string) IEntry
	Queue(data string) IEntry
	RoutingKey(data string) IEntry
	Publish(message []byte, opts ...EntryOption) error
	PublishDelay(message []byte, delay DelayTime, opts ...EntryOption) error
	Consume(opts ...EntryOption) (<-chan amqp.Delivery, error)
	ConsumeFn(fn Callback, opts ...EntryOption) error
	WithWorkers(worker int) IEntry
	WithRetryCount(retryCount int) IEntry
}

func NewEntry(mq *RabbitMQ) *Entry {
	return &Entry{
		ctx:            context.Background(),
		mq:             mq,
		exchangeOption: NewExchangeOption(),
		queueOption:    NewQueueOption(),
		consumeOption:  NewConsumeOption(),
		publishOption:  NewPublishOption(),
		headers:        make(amqp.Table),
	}
}

func (e *Entry) WithContext(ctx context.Context) IEntry {
	e.ctx = ctx
	return e
}

func (e *Entry) WithWorkers(worker int) IEntry {
	e.consumeOption.Workers = worker
	return e
}

// WithRetryNum 设置重试次数, Default: 3次
func (e *Entry) WithRetryCount(retryCount int) IEntry {
	max := len(delayStage)
	if retryCount > max {
		retryCount = max - 1
	}
	if retryCount < 1 {
		retryCount = 1
	}
	e.consumeOption.RetryCount = retryCount
	return e
}

func (e *Entry) ChannelCloseCh() <-chan *amqp.Error {
	return e.channelCloseCh
}

func (e *Entry) Channel(ch *amqp.Channel) IEntry {
	e.channel = ch
	return e
}
func (e *Entry) Close() {
	e.channel.Close()
}

// WatchChannel 监听通道关闭通知，自动创建新通道
func (e *Entry) WatchChannel() {
	e.mq.WaitReConnect()
	e.channel = e.mq.NewChannel()
	e.channelCloseCh = make(chan *amqp.Error, 1)
	e.channel.NotifyClose(e.channelCloseCh)
}

func (e *Entry) Exchange(data string) IEntry {
	e.exchangeName = data
	return e
}

func (e *Entry) Queue(data string) IEntry {
	e.queueName = data
	return e
}

func (e *Entry) RoutingKey(data string) IEntry {
	e.routingName = data
	return e
}

type EntryOption func(*Entry)

func (o EntryOption) Apply(e *Entry) {
	o(e)
}
func WithQueueOption(data *QueueOption) EntryOption {
	return func(e *Entry) {
		e.queueOption = data
	}
}

// WithConsumeOption 消费参数
func WithConsumeOption(data *ConsumeOption) EntryOption {
	return func(e *Entry) {
		e.consumeOption = data
	}
}

// WithConsumeConsumer 消费者名
func WithConsumeConsumer(data string) EntryOption {
	return func(e *Entry) {
		e.consumeOption.Consumer = data
	}
}

// WithConsumeAutoAck 消费自动应答
func WithConsumeAutoAck(data bool) EntryOption {
	return func(e *Entry) {
		e.consumeOption.AutoAck = data
	}
}

// WithConsumeArgs 消费附加参数
func WithConsumeArgs(data map[string]interface{}) EntryOption {
	return func(e *Entry) {
		e.consumeOption.Args = data
	}
}

// WithConsumeNoWait 消费阻塞等待
func WithConsumeNoWait(data bool) EntryOption {
	return func(e *Entry) {
		e.consumeOption.NoWait = data
	}
}

// WithHeaders 生产headers
func WithHeaders(data map[string]interface{}) EntryOption {
	return func(e *Entry) {
		e.headers = data
	}
}

// WithHeader 生产header
func WithHeader(key string, value interface{}) EntryOption {
	return func(e *Entry) {
		e.headers[key] = value
	}
}

// WithMsgExpire 生产消息有效时间
func WithMsgExpire(duration time.Duration) EntryOption {
	return func(e *Entry) {
		second := duration.Milliseconds()
		if second == 0 {
			return
		}
		e.publishOption.Expiration = strconv.Itoa(int(second))
	}
}

func (e *Entry) NewQueue() *Queue {
	return &Queue{
		ch:     e.mq.ChannelAcquire(),
		Option: e.queueOption,
	}
}

func (e *Entry) NewExchange() *Exchange {
	return &Exchange{
		ch:     e.mq.ChannelAcquire(),
		Option: e.exchangeOption,
	}
}

func (e *Entry) InjectHeaders() {
	e.headers["sender"] = e.mq.serverName
	propagator := otel.GetTextMapPropagator()
	carr := propagation.MapCarrier{}
	propagator.Inject(e.ctx, &carr)
	for k, v := range carr {
		e.headers[k] = v
	}
}

func (e *Entry) ExtractHeaders(data *amqp.Delivery) context.Context {
	propagator := otel.GetTextMapPropagator()
	carr := propagation.MapCarrier{}
	if traceparent, ok := data.Headers["traceparent"].(string); ok {
		carr.Set("traceparent", traceparent)
	}
	if tracestate, ok := data.Headers["tracestate"].(string); ok {
		carr.Set("tracestate", tracestate)
	}
	ctx := propagator.Extract(e.ctx, &carr)
	return ctx
}

func NewPublishOption() *PublishOption {
	return &PublishOption{
		Priority:   0,
		Expiration: "",
		MessageId:  "",
	}
}

func (e *Entry) RetryPublish(fc func() error, number int) error {
	err := fc()
	if err == nil {
		return nil
	}
	if number == 0 {
		return err
	}
	number--
	time.After(50 * time.Millisecond)
	return e.RetryPublish(fc, number)
}

func (e *Entry) Publish(message []byte, opts ...EntryOption) error {
	for _, o := range opts {
		o.Apply(e)
	}
	// 定义队列
	// e.queueOption.QueueName = e.queueName
	// q := e.NewQueue().Declare()
	// if q.Error != nil {
	// 	return fmt.Errorf("declare queue: %w", q.Error)
	// }
	e.InjectHeaders()
	return e.RetryPublish(func() error {
		return e.mq.ChannelAcquire().Publish(
			e.exchangeName, // 默认的Exchange交换机是default,类型是direct直接类型
			e.routingName,  // 要赋值的队列名称
			false,          // 如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
			false,          // 如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
			//消息
			amqp.Publishing{
				//类型
				ContentType: "application/json",
				Headers:     e.headers,
				Body:        message,
				Expiration:  e.publishOption.Expiration,
			})
	}, 3)
}

func (e *Entry) PublishDelay(message []byte, delay DelayTime, opts ...EntryOption) error {
	for _, o := range opts {
		o.Apply(e)
	}
	// 停止发布
	if delay <= -1 {
		return nil
	}
	// 定义定时交换机
	tickExchange := fmt.Sprintf(delayExchangeDefault, delay)
	if q := e.NewExchange().WithExchangeName(tickExchange).Declare(); q.Error != nil {
		return q.Error
	}
	// 定义死信交换机
	if q := e.NewExchange().WithExchangeName(delayDlX).Declare(); q.Error != nil {
		return q.Error
	}
	// 定义目标队列
	if e.queueName == "" {
		e.queueName = e.routingName
	}
	if q := e.NewQueue().WithQueueName(e.queueName).Declare().Bind(delayDlX, e.queueName); q.Error != nil {
		return q.Error
	}

	// 定义延迟队列
	e.queueOption.Args = amqp.Table{
		"x-dead-letter-exchange": delayDlX,
		"x-queue-mode":           "lazy",
		"x-message-ttl":          delay.Int64(),
	}
	delayQueueName := fmt.Sprintf(delayQueue, delay)
	if q := e.NewQueue().WithQueueName(delayQueueName).Declare().Bind(tickExchange, e.queueName); q.Error != nil {
		return q.Error
	}

	e.headers[nextDelayTimeKey] = delay.Next().Int64()
	// 注入trace追踪
	e.InjectHeaders()
	return e.RetryPublish(func() error {
		return e.mq.ChannelAcquire().Publish(
			tickExchange, // 默认的Exchange交换机是default,类型是direct直接类型
			e.queueName,  // 要赋值的队列名称
			false,        // 如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
			false,        // 如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
			// 消息
			amqp.Publishing{
				//类型
				ContentType: "application/json",
				Headers:     e.headers,
				Body:        message,
				Expiration:  e.publishOption.Expiration,
			})
	}, 3)
}

func NewConsumeOption() *ConsumeOption {
	return &ConsumeOption{
		AutoAck:    true,
		Exclusive:  false,
		NoLocal:    false,
		NoWait:     false,
		Workers:    1,
		RetryCount: 3,
	}
}

func (e *Entry) Consume(opts ...EntryOption) (<-chan amqp.Delivery, error) {
	for _, o := range opts {
		o.Apply(e)
	}
	e.WatchChannel()
	// 定义队列
	e.queueOption.QueueName = e.queueName
	if e.routingName == "" {
		e.routingName = e.queueName
	}
	q := e.NewQueue().Declare().Bind(e.exchangeName, e.routingName)
	if q.Error != nil {
		return nil, q.Error
	}

	// 根据worker数量动态设置QoS
	// prefetchCount := e.consumeOption.Workers
	// if prefetchCount == 0 {
	// 	prefetchCount = 1
	// }
	// 设置prefetchCount为worker数量的2倍，确保worker不会空闲
	// err := e.channel.Qos(prefetchCount*2, 0, false)
	// if err != nil {
	// 	log.Printf("failed to set QoS: %v", err)
	// 	return nil, err
	// }

	//接收消息
	return e.channel.Consume(
		e.queueName,
		e.consumeOption.Consumer,  // 用来区分多个消费者
		e.consumeOption.AutoAck,   // 是否自动应答
		e.consumeOption.Exclusive, // 是否具有排他性
		e.consumeOption.NoLocal,   // 如果设置为true,表示不能同一个connection中发送的消息传递给这个connection中的消费者
		e.consumeOption.NoWait,    // 队列是否阻塞
		e.consumeOption.Args,
	)
}

type ConsumeMessage struct {
	*amqp.Delivery
	Ctx     context.Context
	Entry   *Entry
	Message []byte
}
type Callback func(msg *ConsumeMessage) error

func (e *Entry) DeferAck(resp *amqp.Delivery, sign *bool) (err error) {
	// recovery 会导致 mq 的连接损坏, 无法继续使用. 且业务调用方无法感知自己的代码 panic, 导致问题被隐藏, 相关的部分服务还不可用
	// 所以先去掉 recovery, 暴露问题
	// if err := recover(); err != nil {
	// 	log.Errorf("consume fn panic: %v", err)
	// }
	if !*sign {
		err = resp.Reject(false)
		if err != nil {
			log.Printf("failed to reject message: %v", err)
		}
		return
	}
	err = resp.Ack(false)
	if err != nil {
		log.Printf("failed to ack message: %v", err)
	}
	return
}

// ConsumeFn 消费方法 返回error 会进行延迟重试， 返回NotRequeueErr 或 nil 会正常ack
func (e *Entry) ConsumeFn(fn Callback, opts ...EntryOption) error {
	// 强制手动ack
	opts = append(opts, WithConsumeAutoAck(false))
	tracer := otel.Tracer("RabbitMQ")
	// 避免 goroutine 泄漏
	if e.consumeOption.Workers == 0 {
		e.consumeOption.Workers = 1
	}
	workChan := make(chan struct{}, e.consumeOption.Workers)
Start:
	ch, err := e.Consume(opts...)
	if err != nil {
		fmt.Println("[RabbitMQ] ConsumeFn Error: ", err)
		goto Start
	}
	for {
		select {
		case r := <-ch:
			workChan <- struct{}{}
			go func(res amqp.Delivery) {
				defer func() {
					<-workChan
				}()

				success := true

				ctx := e.ExtractHeaders(&res)
				ctx, span := tracer.Start(ctx, e.queueName, trace.WithSpanKind(trace.SpanKindConsumer))
				defer span.End()

				data := ConsumeMessage{
					Delivery: &res,
					Entry:    e,
					Ctx:      ctx,
					Message:  r.Body,
				}
				err := fn(&data)

				if err != nil {
					span.SetStatus(codes.Error, err.Error())
					// 先 ACK 避免重复重试
					ackErr := e.DeferAck(&res, &success)
					if ackErr != nil {
						log.Printf("failed to defer ack message: %v", ackErr)
						return
					}
					// 异常错误直接重试
					err = data.RetryError(err)
					if err != nil {
						log.Printf("ConsumeFn Retry Error: %s", e.queueName)
					}
					return
				}
				span.SetStatus(codes.Ok, "")
				e.DeferAck(&res, &success)
			}(r)
		case <-e.channelCloseCh:
			// 等待所有 worker 完成
			for i := 0; i < e.consumeOption.Workers; i++ {
				// 填满 channel，阻止新的 goroutine 启动
				workChan <- struct{}{}
			}
			for i := 0; i < e.consumeOption.Workers; i++ {
				// 等待所有正在运行的 goroutine 完成
				<-workChan
			}
			goto Start
		case <-e.ctx.Done():
			e.channel.Close()
			return nil
		}
	}
}

// GetNextDelayTime 根据已有消息体获取下次重试的间距时间
func (c *ConsumeMessage) GetNextDelayTime() DelayTime {
	nextDelayTime, ok := c.Headers[nextDelayTimeKey].(int64)
	if !ok {
		return delayStage[0]
	}
	return DelayTime(nextDelayTime)
}

// RetryError 消息错误重试
func (c *ConsumeMessage) RetryError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, NotRequeueErr) {
		log.Printf("ConsumeFn Not Requeue: %s", c.Entry.queueName)
		return nil
	}
	// 判断重试次数
	nextDelayTime := c.GetNextDelayTime()
	if nextDelayTime == -1 || nextDelayTime > delayStage[c.Entry.consumeOption.RetryCount-1] {
		log.Printf("⚠️ ConsumeFn RetryCount Exceeded: %s", c.Entry.queueName)
		return nil
	}
	log.Printf("ConsumeFn Retry: %s, Delay: %d", c.Entry.queueName, nextDelayTime)
	defer c.RetryEndLog(err)
	return c.Retry()
}

// Retry 消息体重试
func (c *ConsumeMessage) Retry() error {
	return c.Entry.mq.WithContext(c.Ctx).RoutingKey(c.Entry.queueName).PublishDelay(c.Message, c.GetNextDelayTime())
}

// RetryEndLog 消息体重试日志记录
func (c *ConsumeMessage) RetryEndLog(err error) {
	switch c.GetNextDelayTime() {
	case -1: // 最后一次重试记录
		log.Printf("ConsumeFn Last Retry: %s", c.Entry.queueName)
	case Delay1s: // 第一次重试记录
		log.Printf("ConsumeFn First Retry: %s", c.Entry.queueName)
	}
}

func (c *ConsumeMessage) UnmarshalJson(target interface{}) error {
	return json.Unmarshal(c.Body, target)
}

var (
	NotRequeueErr = NotRequeueError{}
)

type NotRequeueError struct{}

func (e NotRequeueError) Error() string {
	return "disable requeue"
}
