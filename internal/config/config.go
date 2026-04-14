package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
	Version  string `mapstructure:"version"`
	AdminKey string `mapstructure:"admin_key"`

	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Webhook WebhookConfig `mapstructure:"webhook"`

	WebSocket WebSocketConfig `mapstructure:"websocket"`
	Worker    WorkerConfig    `mapstructure:"worker"`
	Auth      AuthConfig      `mapstructure:"auth"`
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	DSN         string `mapstructure:"dsn"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Database    string `mapstructure:"database"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxLifetime time.Duration `mapstructure:"max_lifetime"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConn  int           `mapstructure:"min_idle_conn"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	Timeout     time.Duration `mapstructure:"timeout"`
	MaxRetries  int           `mapstructure:"max_retries"`
	QueueSize   int           `mapstructure:"queue_size"`
	BatchSize   int           `mapstructure:"batch_size"`
	DirectStream bool         `mapstructure:"direct_stream"`
}

// WebSocketConfig WebSocket 配置
type WebSocketConfig struct {
	HandshakeTimeout time.Duration `mapstructure:"handshake_timeout"`
	ReadBufferSize   int           `mapstructure:"read_buffer_size"`
	WriteBufferSize  int           `mapstructure:"write_buffer_size"`
	ReadDeadline     time.Duration `mapstructure:"read_deadline"`
	WriteDeadline    time.Duration `mapstructure:"write_deadline"`
	PingInterval     time.Duration `mapstructure:"ping_interval"`
	CheckInterval    time.Duration `mapstructure:"check_interval"`
	MaxMessageSize   int64         `mapstructure:"max_message_size"`
}

// WorkerConfig 工作池配置
type WorkerConfig struct {
	PoolSize    int `mapstructure:"pool_size"`
	MaxTaskLen  int `mapstructure:"max_task_len"`
	WaitTime    int `mapstructure:"wait_time_ms"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	TokenExpire    time.Duration `mapstructure:"token_expire"`
	RefreshExpire  time.Duration `mapstructure:"refresh_expire"`
	AutoAuthInterval time.Duration `mapstructure:"auto_auth_interval"`
}

// Load 加载配置
func Load() (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 配置文件路径
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// 支持环境变量
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 解析到结构体
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖敏感配置
	if cfg.AdminKey == "" {
		cfg.AdminKey = os.Getenv("ADMIN_KEY")
	}
	if cfg.MySQL.DSN == "" {
		cfg.MySQL.DSN = os.Getenv("MYSQL_CONNECT_STR")
	}

	return cfg, nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// 服务器配置
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("port", 1238)
	v.SetDefault("debug", true)
	v.SetDefault("version", "1.0.0")
	v.SetDefault("admin_key", "12345")

	// MySQL 配置
	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.database", "weixin")
	v.SetDefault("mysql.username", "weixin")
	v.SetDefault("mysql.password", "123456")
	v.SetDefault("mysql.max_open_conn", 100)
	v.SetDefault("mysql.max_idle_conn", 10)
	v.SetDefault("mysql.max_lifetime", "1h")

	// Redis 配置
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 1)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conn", 5)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	// Webhook 配置
	v.SetDefault("webhook.enabled", true)
	v.SetDefault("webhook.timeout", "10s")
	v.SetDefault("webhook.max_retries", 3)
	v.SetDefault("webhook.queue_size", 1000)
	v.SetDefault("webhook.batch_size", 20)
	v.SetDefault("webhook.direct_stream", true)

	// WebSocket 配置
	v.SetDefault("websocket.handshake_timeout", "10s")
	v.SetDefault("websocket.read_buffer_size", 4096)
	v.SetDefault("websocket.write_buffer_size", 4096)
	v.SetDefault("websocket.read_deadline", "120s")
	v.SetDefault("websocket.write_deadline", "60s")
	v.SetDefault("websocket.ping_interval", "25s")
	v.SetDefault("websocket.check_interval", "45s")
	v.SetDefault("websocket.max_message_size", 8192)

	// Worker 配置
	v.SetDefault("worker.pool_size", 500)
	v.SetDefault("worker.max_task_len", 1000)
	v.SetDefault("worker.wait_time", 500)

	// Auth 配置
	v.SetDefault("auth.token_expire", "24h")
	v.SetDefault("auth.refresh_expire", "168h")
	v.SetDefault("auth.auto_auth_interval", "30m")
}