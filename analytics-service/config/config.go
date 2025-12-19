package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App       App       `yaml:"app"`
		HTTP      HTTP      `yaml:"http"`
		Postgres  Postgres  `yaml:"postgres"`
		Log       Log       `yaml:"logger"`
		Kafka     Kafka     `yaml:"kafka"`
		Prometheus Prometheus `yaml:"prometheus"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Postgres struct {
		URL            string        `env-required:"true" yaml:"url" env:"POSTGRES_URL"`
		ConnectTimeout time.Duration `env-required:"true" yaml:"connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT"`
	}

	Kafka struct {
		Brokers []string `env-required:"true" yaml:"brokers" env:"KAFKA_BROKERS"`
		Topics  struct {
			OrderEvents string `env-required:"true" yaml:"order_events" env:"KAFKA_TOPIC_ORDER_EVENTS"`
		} `env-required:"true" yaml:"topics" env:"KAFKA_TOPICS"`
		Consumer KafkaConsumer `yaml:"consumer"`
	}

	KafkaConsumer struct {
		GroupID           string        `yaml:"group_id" env:"KAFKA_CONSUMER_GROUP_ID"`
		MaxWait           time.Duration `yaml:"max_wait" env:"KAFKA_CONSUMER_MAX_WAIT"`
		SessionTimeout    time.Duration `yaml:"session_timeout" env:"KAFKA_CONSUMER_SESSION_TIMEOUT"`
		HeartbeatInterval time.Duration `yaml:"heartbeat_interval" env:"KAFKA_CONSUMER_HEARTBEAT_INTERVAL"`
		CommitInterval    time.Duration `yaml:"commit_interval" env:"KAFKA_CONSUMER_COMMIT_INTERVAL"`
	}

	Prometheus struct {
		Enabled bool   `yaml:"enabled" env:"PROMETHEUS_ENABLED"`
		Path    string `yaml:"path" env:"PROMETHEUS_PATH"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}

