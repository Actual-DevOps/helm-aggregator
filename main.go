package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type HelmRepo struct {
	URL   string `yaml:"url"`
	Name  string `yaml:"name"`
	Index map[string]interface{}
	Lock  sync.Mutex
}

type Config struct {
	Repos []HelmRepo `yaml:"repos"`
	Port  string     `yaml:"port"`
}

var (
	config Config
)

func loadConfig(filename string) error {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Can't load config: %v", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return fmt.Errorf("Can't parse config: %v", err)
	}

	return nil
}

func (repo *HelmRepo) loadIndex() error {
	resp, err := http.Get(repo.URL + "/index.yaml")
	if err != nil {
		return fmt.Errorf("Can't get indexies for repo %s: %v", repo.Name, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading body response for repo %s: %v", repo.Name, err)
	}

	var index map[string]interface{}

	err = yaml.Unmarshal(body, &index)
	if err != nil {
		return fmt.Errorf("Error parsing index for repo %s: %v", repo.Name, err)
	}

	repo.Lock.Lock()
	defer repo.Lock.Unlock()
	repo.Index = index

	return nil
}

func aggregateIndexes() (map[string]interface{}, error) {
	aggregatedIndex := make(map[string]interface{})
	entries := make(map[string]interface{})

	for _, repo := range config.Repos {
		repo.Lock.Lock()
		if repo.Index != nil {
			if repoIndexEntries, ok := repo.Index["entries"].(map[string]interface{}); ok {
				for chart, versions := range repoIndexEntries {
					entries[chart] = versions
				}
			}
		}
		repo.Lock.Unlock()
	}

	aggregatedIndex["apiVersion"] = "v1"
	aggregatedIndex["entries"] = entries
	aggregatedIndex["generated"] = "2023-10-01T00:00:00Z" // Пример даты генерации

	return aggregatedIndex, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	aggregatedIndex, err := aggregateIndexes()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error aggegation index: %v", err), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	yaml.NewEncoder(w).Encode(aggregatedIndex)
}

func main() {
	err := loadConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for i := range config.Repos {

		wg.Add(1)

		go func(repo *HelmRepo) {
			defer wg.Done()

			err := repo.loadIndex()
			if err != nil {
				fmt.Println(err)
			}
		}(&config.Repos[i])
	}

	wg.Wait()

	http.HandleFunc("/index.yaml", indexHandler)

	fmt.Printf("Run on port %s\n", config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Printf("Error run server: %v\n", err)
		os.Exit(1)
	}
}
