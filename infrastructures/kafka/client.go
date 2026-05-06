package kafka

import (
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Option func(info *BasicInfo)

type BasicInfo struct {
	Brokers  []string
	ClientID string
	Timeout  time.Duration
	GroupID  string
	Topic    string
	Balancer kafka.Balancer
}

func WithBalancer(balancer kafka.Balancer) Option {
	return func(info *BasicInfo) {
		info.Balancer = balancer
	}
}

func WithBrokers(brokers []string) Option {
	return func(info *BasicInfo) {
		info.Brokers = brokers
	}
}

func WithClientID(clientID string) Option {
	return func(info *BasicInfo) {
		info.ClientID = clientID
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(info *BasicInfo) {
		info.Timeout = timeout
	}
}

func WithGroupID(groupID string) Option {
	return func(info *BasicInfo) {
		info.GroupID = groupID
	}
}

func WithTopic(topic string) Option {
	return func(info *BasicInfo) {
		info.Topic = topic
	}
}

func NewKafkaProducerClient(opts ...Option) *kafka.Writer {
	// 生成基础数据
	basicInfo := &BasicInfo{
		Brokers:  []string{},
		ClientID: uuid.New().String(),
		Timeout:  10 * time.Second,
		GroupID:  uuid.New().String()[:7],
		Topic:    "default",
		Balancer: &kafka.RoundRobin{},
	}
	// 编译应用选项
	for _, opt := range opts {
		opt(basicInfo)
	}
	// 连接拨号器
	transport := &kafka.Transport{
		ClientID:    basicInfo.ClientID,
		DialTimeout: basicInfo.Timeout,
	}
	// 构造Client
	client := &kafka.Writer{
		Addr:                   kafka.TCP(basicInfo.Brokers...),
		Topic:                  basicInfo.Topic,
		MaxAttempts:            3,
		BatchSize:              1,
		BatchTimeout:           5 * time.Millisecond,
		RequiredAcks:           kafka.RequireAll,
		Async:                  false,
		AllowAutoTopicCreation: true,
		Transport:              transport,
		WriteTimeout:           5 * time.Second, // 新增写超时，避免快速失败
		ReadTimeout:            5 * time.Second, // 新增读超时
		Balancer:               basicInfo.Balancer,
	}

	return client
}

func NewKafkaConsumerGroup(opts ...Option) *kafka.Reader {
	// 生成基础数据
	basicInfo := &BasicInfo{
		Brokers:  []string{},
		ClientID: uuid.New().String(),
		Timeout:  10 * time.Second,
		GroupID:  uuid.New().String()[:7],
		Topic:    "default",
		Balancer: &kafka.RoundRobin{},
	}
	// 编译应用选项
	for _, opt := range opts {
		opt(basicInfo)
	}
	// 连接拨号器
	dialer := &kafka.Dialer{
		ClientID: basicInfo.ClientID,
		Timeout:  basicInfo.Timeout,
	}
	// 构造Client
	client := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     basicInfo.Brokers,
		GroupID:     basicInfo.GroupID,
		Dialer:      dialer,
		Topic:       basicInfo.Topic,
		StartOffset: kafka.LastOffset,
		// MinBytes:    1e3,  // 1KB
		// MaxBytes:    10e6, // 10MB
	})

	return client
}
