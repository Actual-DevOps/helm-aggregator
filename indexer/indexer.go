package indexer

import (
	"fmt"
	"time"

	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
)

func AggregateIndexes(config conf.Config) (map[string]interface{}, error) {
	aggregatedIndex := make(map[string]interface{})
	entries := make(map[string]interface{})

	for _, repo := range config.Repos {
		repo.Lock.Lock()

		if repo.Index != nil {
			if repoIndexEntries, ok := repo.Index["entries"].(map[interface{}]interface{}); ok {
				for chart, versions := range repoIndexEntries {
					entries[fmt.Sprintf("%s/%s", repo.Name, chart.(string))] = versions
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
