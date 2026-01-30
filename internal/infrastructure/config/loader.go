package config

import (
	"os"
	"path/filepath"

	domainConfig "pact/internal/domain/config"
	"pact/internal/domain/errors"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".pactconfig"

// Loader は設定ファイルの読み書きを行う
type Loader struct{}

// NewLoader は新しいLoaderを作成する
func NewLoader() *Loader {
	return &Loader{}
}

// Load は設定ファイルを読み込む
func (l *Loader) Load(path string) (*domainConfig.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// ファイルが存在しない場合はデフォルト設定を返す
			return domainConfig.Default(), nil
		}
		return nil, &errors.ConfigError{Path: path, Message: err.Error()}
	}

	cfg := domainConfig.Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, &errors.ConfigError{Path: path, Message: "invalid YAML: " + err.Error()}
	}

	return cfg, nil
}

// Save は設定をファイルに保存する
func (l *Loader) Save(path string, cfg *domainConfig.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return &errors.ConfigError{Path: path, Message: err.Error()}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &errors.ConfigError{Path: path, Message: err.Error()}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return &errors.ConfigError{Path: path, Message: err.Error()}
	}

	return nil
}

// FindProjectRoot はプロジェクトルートを探す
func (l *Loader) FindProjectRoot(startPath string) (string, error) {
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", &errors.ConfigError{Message: err.Error()}
	}

	for {
		configPath := filepath.Join(absPath, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return absPath, nil
		}

		parent := filepath.Dir(absPath)
		if parent == absPath {
			// ルートに到達
			return "", &errors.ConfigError{Message: "project root not found"}
		}
		absPath = parent
	}
}
