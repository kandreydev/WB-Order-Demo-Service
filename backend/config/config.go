package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	UserName string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
}

type AppConfig struct {
	Port     string `env:"APP_PORT" env-required:"true"`
	Address  string `env:"APP_ADDRESS" env-required:"true"`
	LogLevel string `env:"APP_LOG_LEVEL" env-required:"true"`
}

type CacheConfig struct {
	TTL             time.Duration `env:"CACHE_TTL" env-required:"true"`
	CleanupInterval time.Duration `env:"CACHE_CLEANUP_INTERVAL" env-required:"true"`
}

type KafkaConfig struct {
	KafkaBroker  string `env:"KAFKA_BROKERS" env-required:"true"`
	KafkaTopic   string `env:"KAFKA_TOPIC" env-required:"true"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID" env-required:"true"`
}

type CorsConfig struct {
	Enabled        bool     `env:"CORS_ENABLED" env-default:"false"`
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" env-separator:","`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS" env-separator:"," env-default:"GET,OPTIONS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" env-separator:"," env-default:"*"`
}

type Config struct {
	Postgres         PostgresConfig
	App              AppConfig
	Cache            CacheConfig
	Kafka            KafkaConfig
	MigratePath      string `env:"MIGRATE_PATH" env-required:"true"`
	EmulatorMessages int    `env:"EMULATOR_MESSAGES" env-default:"50"`
	Cors             CorsConfig
}

func LoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err.Error())
	}
	return cfg
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.UserName, c.Postgres.Password, c.Postgres.DBName)
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Postgres.UserName, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DBName)
}

func (c *Config) GetKafkaBrokers() []string {
	return []string{c.Kafka.KafkaBroker}
}

func (c *Config) GetKafkeTopics() []string {
	return []string{c.Kafka.KafkaTopic}
}
