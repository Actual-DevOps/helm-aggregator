package conf

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

const (
	HelmAggregatorConfig string = "HELM_AGGREGATOR_CONFIG"
)

type HelmRepo struct {
	URL    string `yaml:"url"`
	Name   string `yaml:"name"`
	Charts []any  `yaml:"charts"`
	Index  map[string]any
	Lock   sync.Mutex
}

type Config struct {
	Repos []HelmRepo `yaml:"repos"`
	Port  string     `yaml:"port"`
}

func setDefaultConfigPath() {
	if os.Getenv(HelmAggregatorConfig) == "" {
		os.Setenv(HelmAggregatorConfig, "config.yaml")
	}
}

func LoadConfig(config *Config) error {
	setDefaultConfigPath()

	viper.SetConfigFile(os.Getenv(HelmAggregatorConfig))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("configuration read error: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("configuration parsing error: %w", err)
	}

	return nil
}
