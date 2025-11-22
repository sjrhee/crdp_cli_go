package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config는 애플리케이션 설정을 나타냅니다
type Config struct {
	API struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Timeout int    `yaml:"timeout"`
		TLS     bool   `yaml:"tls"`
	} `yaml:"api"`

	Protection struct {
		Policy string `yaml:"policy"`
	} `yaml:"protection"`

	Execution struct {
		Iterations int    `yaml:"iterations"`
		StartData  string `yaml:"start_data"`
	} `yaml:"execution"`

	Batch struct {
		Enabled bool `yaml:"enabled"`
		Size    int  `yaml:"size"`
	} `yaml:"batch"`

	Output struct {
		ShowProgress bool   `yaml:"show_progress"`
		ShowBody     bool   `yaml:"show_body"`
		Verbose      bool   `yaml:"verbose"`
		File         string `yaml:"file"`
	} `yaml:"output"`

	// JWT 인증 설정
	Auth struct {
		JWT      bool   `yaml:"jwt"`
		JWTToken string `yaml:"jwt_token"`
	} `yaml:"auth"`

	// crdp_file_converter 호환성 섹션
	File struct {
		Delimiter  string `yaml:"delimiter"`
		Column     int    `yaml:"column"`
		SkipHeader bool   `yaml:"skip_header"`
	} `yaml:"file"`

	Parallel struct {
		Workers int `yaml:"workers"`
	} `yaml:"parallel"`
}

// LoadConfig는 config.yaml 파일을 읽어 설정을 로드합니다
func LoadConfig(filename string) (*Config, error) {
	// 파일이 없으면 기본값 반환
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// 파일 읽기
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// YAML 파싱
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// DefaultConfig는 기본 설정을 반환합니다
func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.API.Host = "192.168.0.231"
	cfg.API.Port = 32082
	cfg.API.Timeout = 10
	cfg.API.TLS = false
	cfg.Protection.Policy = "P03"
	cfg.Execution.Iterations = 100
	cfg.Execution.StartData = "1234567890123"
	cfg.Batch.Enabled = false
	cfg.Batch.Size = 50
	cfg.Output.ShowProgress = false
	cfg.Output.ShowBody = false
	cfg.Output.Verbose = false
	cfg.Output.File = ""
	// JWT 인증 설정
	cfg.Auth.JWT = false
	cfg.Auth.JWTToken = ""
	// File 설정 (crdp_file_converter 호환)
	cfg.File.Delimiter = ","
	cfg.File.Column = 0
	cfg.File.SkipHeader = false
	// Parallel 설정 (crdp_file_converter 호환)
	cfg.Parallel.Workers = 1
	return cfg
}

// GetConfigPath는 config.yaml의 경로를 반환합니다
// 현재 디렉토리, 실행 파일 디렉토리, 홈 디렉토리 순서로 찾습니다
func GetConfigPath() string {
	paths := []string{
		"config.yaml",
		filepath.Join(os.Getenv("HOME"), ".crdp", "config.yaml"),
	}

	// 실행 파일 디렉토리 확인
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths, filepath.Join(exeDir, "config.yaml"))
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 파일이 없으면 현재 디렉토리의 config.yaml 경로 반환
	return "config.yaml"
}
