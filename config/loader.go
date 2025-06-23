package config

import (
    "fmt"
    "os"
    "sync"
    "github.com/spf13/viper"
)

type Config struct {
    Server ServerConfig `mapstructure:"server"`
    Logger LoggerConfig `mapstructure:"logger"`
    CORS CORSConfig `mapstructure:"cors"`
}

type ServerConfig struct {
    Version string `mapstructure:"version"`
	Port string `mapstructure:"port"`
}

type LoggerConfig struct {
    Level string `mapstructure:"level"`
}

type CORSConfig struct {
    AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

var (
    cfg     *Config
    cfgOnce sync.Once
    loadErr error
)

func GetConfig() *Config {
    return cfg
}

func Load() (*Config, error) {
    cfgOnce.Do(func() {

        configPath := os.Getenv("CONFIG_PATH")
        if configPath == "" {
            loadErr = fmt.Errorf("CONFIG_PATH environment variable not set")
            return
        }

        viper.SetConfigFile(configPath)
        viper.SetConfigType("yaml")

        viper.SetDefault("server.port", 8080)
        viper.SetDefault("database.port", 5432)
        viper.SetDefault("cors.allow_origins", []string{"*"})
        viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE"})


        if err := viper.ReadInConfig(); err != nil {
            loadErr = fmt.Errorf("failed to read config file: %w", err)
            return
        }

        var c Config
        if err := viper.Unmarshal(&c); err != nil {
            loadErr = fmt.Errorf("failed to unmarshal config: %w", err)
            return
        }

        // if err := validation.Config(&c); err != nil {
        //     loadErr = fmt.Errorf("config validation failed: %w", err)
        //     return
        // }

        cfg = &c
    })

    return cfg, loadErr
}
