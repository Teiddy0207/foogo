package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Log      LogConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Redis    RedisConfig
	SMTP     SMTPConfig
	Kafka    KafkaConfig
	Minio    MinIOConfig
}

type AppConfig struct {
	Name           string
	Port           string
	Env            string
	DetectGRPCAddr string
}

type AuthConfig struct {
	DevTokenPrefix string
}

type LogConfig struct {
	Level         string
	JSON          bool
	DailyRotation bool
	EnableFile    bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type KafkaConfig struct {
	Brokers string
	Topic   string
	GroupID string
}

type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	BucketName      string `mapstructure:"bucket_name"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	Enabled         bool   `mapstructure:"enabled"`
}

var (
	instance *Config
	once     sync.Once
	initErr  error
)

func Init() error {
	once.Do(func() {
		if _, err := os.Stat(".env"); err == nil {
			_ = godotenv.Load(".env")
		}

		v := viper.New()
		v.AutomaticEnv()

		cfg := &Config{
			App: AppConfig{
				Name:           strings.TrimSpace(v.GetString("APP_NAME")),
				Port:           strings.TrimSpace(v.GetString("APP_PORT")),
				Env:            strings.TrimSpace(v.GetString("APP_ENV")),
				DetectGRPCAddr: strings.TrimSpace(v.GetString("APP_DETECT_GRPC_ADDR")),
			},
			Log: LogConfig{
				Level:         strings.TrimSpace(v.GetString("APP_LOG_LEVEL")),
				JSON:          v.GetBool("APP_LOG_JSON"),
				DailyRotation: v.GetBool("APP_LOG_DAILY_ROTATION"),
				EnableFile:    v.GetBool("APP_LOG_ENABLE_FILE"),
			},
			Auth: AuthConfig{
				DevTokenPrefix: strings.TrimSpace(v.GetString("APP_AUTH_DEV_TOKEN_PREFIX")),
			},
			Database: DatabaseConfig{
				Host:     strings.TrimSpace(v.GetString("DB_HOST")),
				Port:     strings.TrimSpace(v.GetString("DB_PORT")),
				Name:     strings.TrimSpace(v.GetString("DB_NAME")),
				User:     strings.TrimSpace(v.GetString("DB_USER")),
				Password: strings.TrimSpace(v.GetString("DB_PASSWORD")),
			},
			Redis: RedisConfig{
				Host:     strings.TrimSpace(v.GetString("REDIS_HOST")),
				Port:     strings.TrimSpace(v.GetString("REDIS_PORT")),
				Password: strings.TrimSpace(v.GetString("REDIS_PASSWORD")),
				DB:       strings.TrimSpace(v.GetString("REDIS_DB")),
			},
			SMTP: SMTPConfig{
				Host:     strings.TrimSpace(v.GetString("SMTP_HOST")),
				Port:     strings.TrimSpace(v.GetString("SMTP_PORT")),
				User:     strings.TrimSpace(v.GetString("SMTP_USER")),
				Password: strings.TrimSpace(v.GetString("SMTP_PASSWORD")),
				From:     strings.TrimSpace(v.GetString("SMTP_FROM")),
			},
			Kafka: KafkaConfig{
				Brokers: strings.TrimSpace(v.GetString("KAFKA_BROKERS")),
				Topic: strings.TrimSpace(firstNonEmpty(
					v.GetString("KAFKA_TOPIC_ORDERS"),
					v.GetString("KAFKA_TOPIC"),
				)),
				GroupID: strings.TrimSpace(v.GetString("KAFKA_GROUP_ID")),
			},
			Minio: MinIOConfig{
				Endpoint:        strings.TrimSpace(v.GetString("APP_MINIO_ENDPOINT")),
				AccessKeyID:     strings.TrimSpace(v.GetString("APP_MINIO_ACCESS_KEY_ID")),
				SecretAccessKey: strings.TrimSpace(v.GetString("APP_MINIO_SECRET_ACCESS_KEY")),
				BucketName:      strings.TrimSpace(v.GetString("APP_MINIO_BUCKET_NAME")),
				UseSSL:          v.GetBool("APP_MINIO_USE_SSL"),
				Enabled:         v.GetBool("APP_MINIO_ENABLED"),
			},
		}

		if cfg.App.Port == "" {
			cfg.App.Port = "3000"
		}
		if cfg.App.DetectGRPCAddr == "" {
			cfg.App.DetectGRPCAddr = "localhost:50051"
		}
		if cfg.Log.Level == "" {
			cfg.Log.Level = "INFO"
		}
		if strings.TrimSpace(v.GetString("APP_LOG_ENABLE_FILE")) == "" {
			cfg.Log.EnableFile = true
		}
		if strings.TrimSpace(v.GetString("APP_LOG_DAILY_ROTATION")) == "" {
			cfg.Log.DailyRotation = true
		}
		if cfg.Auth.DevTokenPrefix == "" {
			cfg.Auth.DevTokenPrefix = "dev-token-"
		}

		instance = cfg
	})

	return initErr
}

func Load() (*Config, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	cfg, ok := GetSafe()
	if !ok {
		return nil, fmt.Errorf("config is not initialized")
	}
	return cfg, nil
}

func Get() *Config {
	if cfg, ok := GetSafe(); ok {
		return cfg
	}
	return &Config{}
}

func GetSafe() (*Config, bool) {
	return instance, instance != nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		v := strings.TrimSpace(value)
		if v != "" {
			return v
		}
	}
	return ""
}
