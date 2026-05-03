package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Option func(info *BasicInfo)

type BasicInfo struct {
	MasterName       string
	SentinelAddrs    []string
	Password         string
	SentinelPassword string
	DB               int
}

func WithMasterName(masterName string) Option {
	return func(info *BasicInfo) {
		info.MasterName = masterName
	}
}

func WithSentinelAddrs(sentinelAddrs []string) Option {
	return func(info *BasicInfo) {
		info.SentinelAddrs = sentinelAddrs
	}
}

func WithPassword(password string) Option {
	return func(info *BasicInfo) {
		info.Password = password
	}
}

func WithSentinelPassword(sentinelPassword string) Option {
	return func(info *BasicInfo) {
		info.SentinelPassword = sentinelPassword
	}
}

func WithDB(db int) Option {
	return func(info *BasicInfo) {
		info.DB = db
	}
}

func NewRedisSentinelClient(opts ...Option) (*redis.Client, error) {
	// 构建基础选项
	basicInfo := &BasicInfo{
		MasterName:       "mymaster",
		SentinelAddrs:    []string{"localhost:26379"},
		Password:         "root",
		SentinelPassword: "root",
		DB:               0,
	}
	// 遍历应用
	for _, opt := range opts {
		opt(basicInfo)
	}
	// 创建NewFailoverClient
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       basicInfo.MasterName,
		SentinelAddrs:    basicInfo.SentinelAddrs,
		Password:         basicInfo.Password,
		SentinelPassword: basicInfo.SentinelPassword,
		DB:               basicInfo.DB,
	})
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
