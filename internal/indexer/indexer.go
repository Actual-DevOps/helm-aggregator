package indexer

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"gopkg.in/yaml.v2"
)


func AggregateIndexes(config conf.Config) (map[string]any, error) {
	aggregatedIndex := make(map[string]any)
	entries := make(map[string]any)

	for i := range config.Repos {
		repo := &config.Repos[i]

		repo.Lock.Lock()

		if repo.Index != nil {
			if repoIndexEntries, ok := repo.Index["entries"].(map[any]any); ok {
				for chart, versions := range repoIndexEntries {
					chart, ok := chart.(string)
					if !ok {
						return nil, fmt.Errorf("chart is not a string")
					}

					entries[fmt.Sprintf("%s/%s", repo.Name, chart)] = versions
				}
			}
		}
		repo.Lock.Unlock()
	}

	aggregatedIndex["apiVersion"] = "v1"
	aggregatedIndex["entries"] = entries
	aggregatedIndex["generated"] = time.Now().Format(time.RFC3339)

	return aggregatedIndex, nil
}

func LoadIndex(repo *conf.HelmRepo) error {
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
