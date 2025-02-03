package handlers

import (
	"fmt"
	"net/http"

	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"gopkg.in/yaml.v2"
)

// Обработчик для endpoint /index.yaml
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var config conf.Config
	aggregatedIndex, err := aggregateIndexes(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка агрегации индексов: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	yaml.NewEncoder(w).Encode(aggregatedIndex)
}

func aggregateIndexes(config conf.Config) (map[string]interface{}, error) {
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
