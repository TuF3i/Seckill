package env

import (
	"os"
	"seckill/configs"

	"github.com/google/uuid"
)

func GetEnv() *configs.BasicEnv {
	conf := &configs.BasicEnv{
		NacosAddr:     "localhost",
		NacosPort:     "8848",
		NacosUser:     "admin",
		NacosPassword: "admin",
		ConfigID:      "seckill",
		ConfigGroup:   "REDROCK",
		ContainerName: uuid.New().String()[:7],
	}

	if d := os.Getenv("NACOS_ADDR"); d != "" {
		conf.NacosAddr = d
	}

	if d := os.Getenv("NACOS_PORT"); d != "" {
		conf.NacosPort = d
	}

	if d := os.Getenv("NACOS_USER"); d != "" {
		conf.NacosUser = d
	}

	if d := os.Getenv("NACOS_PASSWORD"); d != "" {
		conf.NacosPassword = d
	}

	if d := os.Getenv("CONFIG_ID"); d != "" {
		conf.ConfigID = d
	}

	if d := os.Getenv("CONFIG_GROUP"); d != "" {
		conf.ConfigGroup = d
	}

	if d := os.Getenv("CONTAINER_NAME"); d != "" {
		conf.ContainerName = d
	}

	if d := os.Getenv("BENCHMARK"); d != "" {
		conf.Benchmark = d
	}

	return conf
}
