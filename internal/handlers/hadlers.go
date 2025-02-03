package handlers

import (
	"fmt"
	"net/http"

	"github.com/Actual-DevOps/helm-aggregator/indexer"
	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"gopkg.in/yaml.v2"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var config conf.Config

	aggregatedIndex, err := indexer.AggregateIndexes(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка агрегации индексов: %v", err), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	yaml.NewEncoder(w).Encode(aggregatedIndex)
}
