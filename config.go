package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the main configuration struct
type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
}

// ServerConfig is the configuration for the server
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// AppConfig is the configuration for the app
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// read config from file
func ReadConfig(filename string) (*Config, error) {
	//  1. 读取文件内容
	configData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//  2. 创建配置对象
	var config Config
	//  3. 解析配置对象
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	//  4. 返回配置对象
	return &config, nil
}
