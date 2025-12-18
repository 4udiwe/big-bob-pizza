package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		HTTP     HTTP     `yaml:"http"`
		Postgres Postgres `yaml:"postgres"`
		Log      Log      `yaml:"logger"`
		Kafka    Kafka    `yaml:"kafka"`
		Outbox   Outbox   `yaml:"outbox"`
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
			OrderEvents   string `env-required:"true" yaml:"order_events" env:"KAFKA_TOPIC_ORDER_EVENTS"`
			PaymentEvents string `env-required:"true" yaml:"payment_events" env:"KAFKA_TOPIC_PAYMENT_EVENTS"`
		} `env-required:"true" yaml:"topics" env:"KAFKA_TOPICS"`
		Producer KafkaProducer `yaml:"producer"`
		Consumer KafkaConsumer `yaml:"consumer"`
	}

	KafkaProducer struct {
		RequiredAcks int           `yaml:"required_acks" env:"KAFKA_PRODUCER_REQUIRED_ACKS"`
		BatchSize    int           `yaml:"batch_size" env:"KAFKA_PRODUCER_BATCH_SIZE"`
		BatchTimeout time.Duration `yaml:"batch_timeout" env:"KAFKA_PRODUCER_BATCH_TIMEOUT"`
		Compression  string        `yaml:"compression" env:"KAFKA_PRODUCER_COMPRESSION"`
	}

	KafkaConsumer struct {
		GroupID           string        `yaml:"group_id" env:"KAFKA_CONSUMER_GROUP_ID"`
		MaxWait           time.Duration `yaml:"max_wait" env:"KAFKA_CONSUMER_MAX_WAIT"`
		SessionTimeout    time.Duration `yaml:"session_timeout" env:"KAFKA_CONSUMER_SESSION_TIMEOUT"`
		HeartbeatInterval time.Duration `yaml:"heartbeat_interval" env:"KAFKA_CONSUMER_HEARTBEAT_INTERVAL"`
		CommitInterval    time.Duration `yaml:"commit_interval" env:"KAFKA_CONSUMER_COMMIT_INTERVAL"`
	}

	Outbox struct {
		Topic           string        `env-required:"true" yaml:"topic" env:"OUTBOX_PUB_TOPIC"`
		BatchLimit      int           `env-required:"true" yaml:"batch_limit" env:"OUTBOX_BATCH_LIMIT"`
		Interval        time.Duration `env-required:"true" yaml:"interval" env:"OUTBOX_INTERVAL"`
		RequeBatchLimit int           `env-required:"true" yaml:"reque_batch_limit" env:"OUTBOX_REQUE_BATCH_LIMIT"`
		RequeInterval   time.Duration `env-required:"true" yaml:"reque_interval" env:"OUTBOX_REQUE_INTERVAL"`
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

