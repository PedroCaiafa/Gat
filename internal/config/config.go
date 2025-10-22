// internal/config/config.go
package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Provider string

const (
	ProviderOpenAI Provider = "openai"
	ProviderGemini Provider = "gemini"
)

type ProviderConfig struct {
	Provider     Provider `mapstructure:"provider"`
	APIKey       string   `mapstructure:"api_key"`
	DefaultModel string   `mapstructure:"default_model"`
}

type AppConfig struct {
	ChatProvider         ProviderConfig `mapstructure:"chat_provider"`
	EmbeddingProvider    ProviderConfig `mapstructure:"embedding_provider"`
	EmbeddingDim         uint64         `mapstructure:"embedding_dim"`
	QdrantAddress        string         `mapstructure:"qdrant_address"`
	QdrantPort           int            `mapstructure:"qdrant_port"`
	QdrantAPIKey         string         `mapstructure:"qdrant_api_key"`
	CollectionName       string         `mapstructure:"collection_name"`
	GlobalQdrantAddress  string         `mapstructure:"global_qdrant_address"` // ideia eh usar um banco local no usuario e um global pros times dividerem as coisas
	GlobalQdrantPort     int            `mapstructure:"global_qdrant_port"`
	GlobalQdrantAPIKey   string         `mapstructure:"global_qdrant_api_key"`
	GlobalCollectionName string         `mapstructure:"global_collection_name"`
	UVPath               string         `mapstructure:"uv_path"`
	QdrantMode           string         `mapstructure:"qdrant_mode"`        // "bundled" | "external" | "disabled"
	QdrantBinaryPath     string         `mapstructure:"qdrant_binary_path"` // optional override; if empty we autodetect next to the pyhelp binary
}

func Default() AppConfig {
	var qdrantBinaryPath string

	switch runtime.GOOS {
	case "windows":
		qdrantBinaryPath = "qdrant.exe"
	case "linux":
		qdrantBinaryPath = "qdrant"
	case "darwin":
		qdrantBinaryPath = "qdrant"
	default:
		qdrantBinaryPath = ""
	}

	return AppConfig{
		ChatProvider: ProviderConfig{
			Provider:     ProviderOpenAI,
			APIKey:       "",
			DefaultModel: "gpt-4o-mini",
		},
		EmbeddingProvider: ProviderConfig{
			Provider:     ProviderOpenAI,
			APIKey:       "",
			DefaultModel: "text-embedding-3-small",
		},
		EmbeddingDim:     1536,
		QdrantAddress:    "localhost",
		QdrantPort:       6333,
		QdrantAPIKey:     "",
		CollectionName:   "gat-py-tools",
		QdrantMode:       "bundled",
		QdrantBinaryPath: qdrantBinaryPath,
	}
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gat"), nil
}

func GetConfigFilePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func EnsureConfigFile() error {
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		// caso default
		config := Default()
		if err := SaveConfig(config); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig() (AppConfig, error) {
	path, err := GetConfigFilePath()
	if err != nil {
		return AppConfig{}, err
	}
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return AppConfig{}, err
	}
	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return AppConfig{}, err
	}
	return config, nil

}

func SaveConfig(config AppConfig) error {
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.Set("chat_provider.provider", config.ChatProvider.Provider)
	v.Set("chat_provider.api_key", config.ChatProvider.APIKey)
	v.Set("chat_provider.default_model", config.ChatProvider.DefaultModel)
	v.Set("embedding_provider.provider", config.EmbeddingProvider.Provider)
	v.Set("embedding_provider.api_key", config.EmbeddingProvider.APIKey)
	v.Set("embedding_provider.default_model", config.EmbeddingProvider.DefaultModel)
	v.Set("qdrant_address", config.QdrantAddress)
	v.Set("qdrant_api_key", config.QdrantAPIKey)
	v.Set("collection_name", config.CollectionName)
	v.Set("uv_path", config.UVPath)
	v.Set("qdrant_mode", config.QdrantMode)
	v.Set("qdrant_binary_path", config.QdrantBinaryPath)

	if err := v.WriteConfigAs(path); err != nil {
		return err
	}

	return nil

}
