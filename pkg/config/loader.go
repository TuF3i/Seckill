package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"seckill/configs"
	n "seckill/infrastructures/nacos"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Loader struct {
	mu  sync.RWMutex
	cfg *configs.Config
}

func NewLoader(nacosClient *n.Client, dataID string, group string) (*Loader, error) {
	content, err := nacosClient.ConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	if err != nil {
		return nil, fmt.Errorf("config loader: get config from nacos: %w", err)
	}

	if content == "" {
		return nil, fmt.Errorf("config loader: empty config content from nacos (dataId=%s, group=%s)", dataID, group)
	}

	var cfg configs.Config
	if err := json.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("config loader: unmarshal config: %w", err)
	}

	l := &Loader{
		cfg: &cfg,
	}

	log.Printf("[ConfigLoader] config loaded successfully from nacos (dataId=%s, group=%s)", dataID, group)
	return l, nil
}

func (l *Loader) GetConfig() *configs.Config {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.cfg
}
