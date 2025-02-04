package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Actual-DevOps/helm-aggregator/indexer"
	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"gopkg.in/yaml.v2"
)

func IndexHandler(config conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aggregatedIndex, err := indexer.AggregateIndexes(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка агрегации индексов: %v", err), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		yaml.NewEncoder(w).Encode(aggregatedIndex)
	}
}

func GetConfigHandler(config conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		configFile, err := os.ReadFile(os.Getenv("HELM_AGGREGATOR_CONFIG"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Can't get config from filesystem: %v", err), http.StatusInternalServerError)

			return
		}

		w.Write([]byte(configFile))
	}
}
