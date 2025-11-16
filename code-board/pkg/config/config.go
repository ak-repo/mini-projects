package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all the configuration for the service
type Config struct {
	// Server struct {
	// 	Name string `mapstructure:"name"`
	// 	Host string `mapstructure:"host"`
	// 	Port int    `mapstructure:"port"`
	// 	Mode string `mapstructure:"mode"`
	// } `mapstructure:"server"`

	Database struct {
		Host               string `mapstructure:"host"`
		Port               int    `mapstructure:"port"`
		User               string `mapstructure:"user"`
		Password           string `mapstructure:"password"`
		Name               string `mapstructure:"name"`
		SSLMode            string `mapstructure:"sslmode"`
		MaxConnections     int    `mapstructure:"max_connections"`
		MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	} `mapstructure:"database"`

	Redis struct {
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		DB         int    `mapstructure:"db"`
		Password   string `mapstructure:"password"`
		Prefix     string `mapstructure:"prefix"`
		TTLMinutes int    `mapstructure:"ttl_minutes"`
	} `mapstructure:"redis"`

	JWT struct {
		Secret              string `mapstructure:"secret"`
		Issuer              string `mapstructure:"issuer"`
		AccessExpireMinutes int    `mapstructure:"access_expire_minutes"`
		RefreshExpireHours  int    `mapstructure:"refresh_expire_hours"`
	} `mapstructure:"jwt"`

	OAuth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURL  string `mapstructure:"redirect_url"`
		} `mapstructure:"google"`
		GitHub struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURL  string `mapstructure:"redirect_url"`
		} `mapstructure:"github"`
	} `mapstructure:"oauth"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
		GroupID string   `mapstructure:"group_id"`
	} `mapstructure:"kafka"`

	Security struct {
		PasswordHashCost int      `mapstructure:"password_hash_cost"`
		AllowSignup      bool     `mapstructure:"allow_signup"`
		AllowedOrigins   []string `mapstructure:"allowed_origins"`
	} `mapstructure:"security"`

	Logging struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	} `mapstructure:"logging"`
}

// LoadConfig reads configuration from a YAML file and environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Try env override first
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Add multiple search paths
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("../../pkg/config")
		v.AddConfigPath("../../../pkg/config")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// try one more fallback manually
		if _, err2 := os.Stat("../../pkg/config/config.yaml"); err2 == nil {
			v.SetConfigFile("../../pkg/config/config.yaml")
			if err3 := v.ReadInConfig(); err3 != nil {
				return nil, fmt.Errorf("failed to read config: %w", err3)
			}
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Println("✅ Loaded config from:", v.ConfigFileUsed())
	return &cfg, nil
}
