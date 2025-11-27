package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Services ServicesConfig `mapstructure:"services"`
	Server   ServerConfig   `mapstructure:"server"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// ------------------------------------------------
// APP
// ------------------------------------------------
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// ------------------------------------------------
// DATABASE
// ------------------------------------------------
type DatabaseConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Name        string `yaml:"name"`
	MaxPoolSize int    `yaml:"max_pool_size"`
	MinPoolSize int    `yaml:"min_pool_size"`
	SSLMode     string `mapstructure:"sslmode"`
}

type DBPoolConfig struct {
	MaxConnections    int           `mapstructure:"max_connections"`
	MinConnections    int           `mapstructure:"min_connections"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime   time.Duration `mapstructure:"max_conn_idle_time"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
}

// ------------------------------------------------
// JWT
// ------------------------------------------------
type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expiry time.Duration `mapstructure:"expiry"`
	Issuer string        `mapstructure:"issuer"`
}

// ------------------------------------------------
// SERVICES
// ------------------------------------------------
type ServicesConfig struct {
	Chat  ServiceConfig `mapstructure:"chat"`
	Front ServiceConfig `mapstructure:"front"`
}

type ServiceConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// ------------------------------------------------
// HTTP SERVER
// ------------------------------------------------
type ServerConfig struct {
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// ------------------------------------------------
// REDIS
// ------------------------------------------------
type RedisConfig struct {
	Host     string   `mapstructure:"host"`
	Port     string   `mapstructure:"port"`
	DB       int      `mapstructure:"db"`
	Password string   `mapstructure:"password"`
	TTL      RedisTTL `mapstructure:"ttl"`
}

type RedisTTL struct {
	UploadSession time.Duration `mapstructure:"upload_session"`
	Presign       time.Duration `mapstructure:"presign"`
}

// ------------------------------------------------
// LOAD
// ------------------------------------------------

func Load() (*Config, error) {
	cfgPath := os.Getenv("CONFIG_PATH") // check if CONFIG_PATH is set
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath) // use exact path if provided
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	// ENV override support
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("❌ Config read error: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("❌ Config parse error: %w", err)
	}

	return &cfg, nil
}
