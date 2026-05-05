package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Option func(info *BasicInfo)

type BasicInfo struct {
	Host        string
	Port        uint64
	NamespaceID string
	UserName    string
	Password    string
	LogLevel    string
}

type Client struct {
	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient
}

func WithHost(host string) Option {
	return func(info *BasicInfo) {
		info.Host = host
	}
}

func WithPort(port uint64) Option {
	return func(info *BasicInfo) {
		info.Port = port
	}
}

func WithNamespaceID(namespaceID string) Option {
	return func(info *BasicInfo) {
		info.NamespaceID = namespaceID
	}
}

func WithUserName(userName string) Option {
	return func(info *BasicInfo) {
		info.UserName = userName
	}
}

func WithPassword(password string) Option {
	return func(info *BasicInfo) {
		info.Password = password
	}
}

func WithLogLevel(logLevel string) Option {
	return func(info *BasicInfo) {
		info.LogLevel = logLevel
	}
}

func NewNacosClient(opts ...Option) (*Client, error) {
	basicInfo := &BasicInfo{
		Host:        "localhost",
		Port:        8848,
		NamespaceID: "public",
		UserName:    "nacos",
		Password:    "nacos",
		LogLevel:    "info",
	}
	for _, opt := range opts {
		opt(basicInfo)
	}

	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(basicInfo.Host, basicInfo.Port),
	}

	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(basicInfo.NamespaceID),
		constant.WithUsername(basicInfo.UserName),
		constant.WithPassword(basicInfo.Password),
		constant.WithLogLevel(basicInfo.LogLevel),
	)

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("nacos naming client: %w", err)
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("nacos config client: %w", err)
	}

	return &Client{
		NamingClient: namingClient,
		ConfigClient: configClient,
	}, nil
}
