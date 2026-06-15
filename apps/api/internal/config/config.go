package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Storage  StorageConfig  `yaml:"storage"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
	Swagger  SwaggerConfig  `yaml:"swagger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Database)
}

// StorageConfig 存储配置
type StorageConfig struct {
	BasePath        string   `yaml:"base_path"`
	ChunkPath       string   `yaml:"chunk_path"`
	DocumentPath    string   `yaml:"document_path"`
	MaxFileSize     int64    `yaml:"max_file_size"`
	ChunkSize       int      `yaml:"chunk_size"`
	AllowedExtensions []string `yaml:"allowed_extensions"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	ChunkExpireHours int `yaml:"chunk_expire_hours"`
}

// SwaggerConfig Swagger 配置
type SwaggerConfig struct {
	Enabled bool   `yaml:"enabled"` // 是否启用 Swagger
	Host    string `yaml:"host"`    // Swagger 主机地址
}

var GlobalConfig *Config

// Load 加载配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}
