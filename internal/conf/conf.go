package conf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type HelmRepo struct {
	URL   string `yaml:"url"`
	Name  string `yaml:"name"`
	Index map[string]any
	Lock  sync.Mutex
}

type Config struct {
	Repos []HelmRepo `yaml:"repos"`
	Port  string     `yaml:"port"`
}

func LoadConfig(config *Config) error {
	if os.Getenv("HELM_AGGREGATOR_CONFIG") == "" {
		os.Setenv("HELM_AGGREGATOR_CONFIG", "config.yaml")
	}

	viper.SetConfigFile(os.Getenv("HELM_AGGREGATOR_CONFIG"))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("configuration read error: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("configuration parsing error: %w", err)
	}

	return nil
}

func (repo *HelmRepo) LoadIndex() error {
	resp, err := http.Get(repo.URL + "/index.yaml")
	if err != nil {
		return fmt.Errorf("error loading the index for the repository %s: %w", repo.Name, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading the response body for a repository %s: %w", repo.Name, err)
	}

	var index map[string]any

	err = yaml.Unmarshal(body, &index)
	if err != nil {
		return fmt.Errorf("index parsing error for a repository %s: %w", repo.Name, err)
	}

	repo.Lock.Lock()
	defer repo.Lock.Unlock()
	repo.Index = index

	return nil
}
