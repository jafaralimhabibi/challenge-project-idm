package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	JWT struct {
		Secret string `mapstructure:"secret"`
		Expire int    `mapstructure:"expire"`
	} `mapstructure:"jwt"`

	MongoDB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"mongodb"`

	Admin struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"admin"`

	Socket struct {
		Namespace string `mapstructure:"namespace"`
	} `mapstructure:"socket"`
}

var cfg *Config

func LoadConfig(path string) error {
	v := viper.New()
	var c Config

	if path == "" {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	} else {
		v.SetConfigFile(path)
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Printf("viper.ReadInConfig error: %v", err)
		log.Printf("viper.AllSettings() = %+v", v.AllSettings())
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := v.Unmarshal(&c); err != nil {
		log.Printf("viper.Unmarshal error: %v", err)
		log.Printf("viper.AllSettings() = %+v", v.AllSettings())
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if c.MongoDB.Host == "" || c.MongoDB.Database == "" {
		log.Printf("config after unmarshal: %+v", c)
		return fmt.Errorf("mongodb.host or mongodb.name is empty (check config.yml keys and mapstructure tags)")
	}

	cfg = &c
	return nil
}

func GetConfig() *Config {
	if cfg == nil {
		log.Panic("config not loaded; call LoadConfig first")
	}
	return cfg
}
