package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Server struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type Config struct {
	// 服务配置
	Server Server `mapstructure:"server" json:"server" yaml:"server"`
}

var cfg *Config

// InitConfig 初始化配置
func InitConfig() error {
	// 设置默认配置
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 8080)

	// 读取配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	if cfg == nil {
		panic("配置未初始化")
		return nil
	}
	return cfg
}
