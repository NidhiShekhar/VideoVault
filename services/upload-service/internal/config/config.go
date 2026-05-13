package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Service ServiceConfig `mapstructure:"service"`
	Server  ServerConfig  `mapstructure:"server"`
	Log     LogConfig     `mapstructure:"log"`
}

type ServiceConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load(configPath string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Allow env variables like VIDEOVAULT_SERVER_PORT to override YAML values.
	v.SetEnvPrefix("VIDEOVAULT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("service.name", "upload-service")
	v.SetDefault("service.env", "dev")
	v.SetDefault("server.port", "8080")
	v.SetDefault("log.level", "info")
}
