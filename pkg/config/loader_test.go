package config

import (
	"encoding/json"
	"seckill/configs"
	"sync"
	"testing"
)

var testJSON = `{
  "postgresql": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "postgres",
    "defaultDB": "adminer",
    "sslMode": "disable"
  },
  "redis": {
    "masterName": "mymaster",
    "sentinelAddrs": ["localhost:26379"],
    "password": "root",
    "sentinelPassword": "root"
  },
  "kafka": {
    "brokers": ["localhost:9092"]
  },
  "gateway": {
    "listenAddr": "0.0.0.0",
    "listenPort": "8888",
    "monitoringPort": "8889"
  }
}`

func TestConfigUnmarshal(t *testing.T) {
	var cfg configs.Config
	if err := json.Unmarshal([]byte(testJSON), &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if cfg.PostgreSQL.Host != "localhost" {
		t.Errorf("expected PostgreSQL.Host 'localhost', got '%s'", cfg.PostgreSQL.Host)
	}
	if cfg.PostgreSQL.Port != 5432 {
		t.Errorf("expected PostgreSQL.Port 5432, got %d", cfg.PostgreSQL.Port)
	}
	if cfg.PostgreSQL.User != "postgres" {
		t.Errorf("expected PostgreSQL.User 'postgres', got '%s'", cfg.PostgreSQL.User)
	}
	if cfg.PostgreSQL.DefaultDB != "adminer" {
		t.Errorf("expected PostgreSQL.DefaultDB 'adminer', got '%s'", cfg.PostgreSQL.DefaultDB)
	}
	if cfg.PostgreSQL.SSLMode != "disable" {
		t.Errorf("expected PostgreSQL.SSLMode 'disable', got '%s'", cfg.PostgreSQL.SSLMode)
	}

	if cfg.Redis.MasterName != "mymaster" {
		t.Errorf("expected Redis.MasterName 'mymaster', got '%s'", cfg.Redis.MasterName)
	}
	if len(cfg.Redis.SentinelAddrs) != 1 || cfg.Redis.SentinelAddrs[0] != "localhost:26379" {
		t.Errorf("expected Redis.SentinelAddrs ['localhost:26379'], got %v", cfg.Redis.SentinelAddrs)
	}
	if cfg.Redis.Password != "root" {
		t.Errorf("expected Redis.Password 'root', got '%s'", cfg.Redis.Password)
	}
	if cfg.Redis.SentinelPassword != "root" {
		t.Errorf("expected Redis.SentinelPassword 'root', got '%s'", cfg.Redis.SentinelPassword)
	}

	if len(cfg.Kafka.Brokers) != 1 || cfg.Kafka.Brokers[0] != "localhost:9092" {
		t.Errorf("expected Kafka.Brokers ['localhost:9092'], got %v", cfg.Kafka.Brokers)
	}

	if cfg.Gateway.ListenAddr != "0.0.0.0" {
		t.Errorf("expected Gateway.ListenAddr '0.0.0.0', got '%s'", cfg.Gateway.ListenAddr)
	}
	if cfg.Gateway.ListenPort != "8888" {
		t.Errorf("expected Gateway.ListenPort '8888', got '%s'", cfg.Gateway.ListenPort)
	}
	if cfg.Gateway.MonitoringPort != "8889" {
		t.Errorf("expected Gateway.MonitoringPort '8889', got '%s'", cfg.Gateway.MonitoringPort)
	}
}

func TestConfigUnmarshalInvalidJSON(t *testing.T) {
	var cfg configs.Config
	err := json.Unmarshal([]byte("{invalid}"), &cfg)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestConfigUnmarshalPartial(t *testing.T) {
	partialJSON := `{
		"postgresql": { "host": "test-host", "port": 1234 },
		"redis": { "masterName": "test-master" }
	}`

	var cfg configs.Config
	if err := json.Unmarshal([]byte(partialJSON), &cfg); err != nil {
		t.Fatalf("failed to unmarshal partial config: %v", err)
	}

	if cfg.PostgreSQL.Host != "test-host" {
		t.Errorf("expected PostgreSQL.Host 'test-host', got '%s'", cfg.PostgreSQL.Host)
	}
	if cfg.PostgreSQL.Port != 1234 {
		t.Errorf("expected PostgreSQL.Port 1234, got %d", cfg.PostgreSQL.Port)
	}
	if cfg.Redis.MasterName != "test-master" {
		t.Errorf("expected Redis.MasterName 'test-master', got '%s'", cfg.Redis.MasterName)
	}

	if len(cfg.Kafka.Brokers) != 0 {
		t.Errorf("expected Kafka.Brokers empty, got %v", cfg.Kafka.Brokers)
	}
}

func TestGetConfigThreadSafe(t *testing.T) {
	var cfg configs.Config
	if err := json.Unmarshal([]byte(testJSON), &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	loader := &Loader{
		cfg: &cfg,
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := loader.GetConfig()
			if c == nil {
				t.Error("GetConfig returned nil")
			}
		}()
	}
	wg.Wait()
}
