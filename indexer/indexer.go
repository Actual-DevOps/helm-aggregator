package indexer

import (
	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
)

func AggregateIndexes(config conf.Config) (map[string]interface{}, error) {
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
	aggregatedIndex["generated"] = "2023-10-01T00:00:00Z"

	return aggregatedIndex, nil
}
