package indexer

import (
	"fmt"
	"time"

	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
)

func AggregateIndexes(config conf.Config) (map[string]interface{}, error) {
	aggregatedIndex := make(map[string]interface{})
	entries := make(map[string]interface{})

	for i := range config.Repos {
		repo := &config.Repos[i]

		repo.Lock.Lock()

		if repo.Index != nil {
			if repoIndexEntries, ok := repo.Index["entries"].(map[interface{}]interface{}); ok {
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
