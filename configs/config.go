package configs

type Config struct {
	PostgreSQL PostgreSQLConfig `json:"postgresql" mapstructure:"postgresql"`
	Redis      RedisConfig      `json:"redis" mapstructure:"redis"`
	Kafka      KafkaConfig      `json:"kafka" mapstructure:"kafka"`
	Gateway    GatewayConfig    `json:"gateway" mapstructure:"gateway"`
}

type PostgreSQLConfig struct {
	Host      string `json:"host" mapstructure:"host"`
	Port      int    `json:"port" mapstructure:"port"`
	User      string `json:"user" mapstructure:"user"`
	Password  string `json:"password" mapstructure:"password"`
	DefaultDB string `json:"defaultDB" mapstructure:"defaultDB"`
	SSLMode   string `json:"sslMode" mapstructure:"sslMode"`
}

type RedisConfig struct {
	MasterName       string   `json:"masterName" mapstructure:"masterName"`
	SentinelAddrs    []string `json:"sentinelAddrs" mapstructure:"sentinelAddrs"`
	Password         string   `json:"password" mapstructure:"password"`
	SentinelPassword string   `json:"sentinelPassword" mapstructure:"sentinelPassword"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers" mapstructure:"brokers"`
}

type GatewayConfig struct {
	ListenAddr     string `json:"listenAddr" mapstructure:"listenAddr"`
	ListenPort     string `json:"listenPort" mapstructure:"listenPort"`
	MonitoringPort string `json:"monitoringPort" mapstructure:"monitoringPort"`
}
