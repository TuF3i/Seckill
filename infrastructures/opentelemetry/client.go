package opentelemetry

import (
	kitex_provider "github.com/kitex-contrib/obs-opentelemetry/provider"
)

type Option func(info *BasicInfo)

type BasicInfo struct {
	Endpoint    string
	ServiceName string
}

func WithEndpoint(endpoint string) Option {
	return func(info *BasicInfo) {
		info.Endpoint = endpoint
	}
}

func WithServiceName(serviceName string) Option {
	return func(info *BasicInfo) {
		info.ServiceName = serviceName
	}
}

func NewProvider(opts ...Option) kitex_provider.OtelProvider {
	basicInfo := &BasicInfo{
		Endpoint:    "localhost:4317",
		ServiceName: "defaultName",
	}

	for _, opt := range opts {
		opt(basicInfo)
	}

	p := kitex_provider.NewOpenTelemetryProvider(
		kitex_provider.WithServiceName(basicInfo.ServiceName),
		kitex_provider.WithExportEndpoint(basicInfo.Endpoint),
		kitex_provider.WithInsecure(),
	)

	return p
}
