package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Actual-DevOps/helm-aggregator/internal/indexer"
	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"gopkg.in/yaml.v2"
)

func IndexHandler(config conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aggregatedIndex, err := indexer.AggregateIndexes(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Can't get aggregatedIndex: %v", err), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		err = yaml.NewEncoder(w).Encode(aggregatedIndex)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encode aggregatedIndex: %v", err), http.StatusInternalServerError)

			return
		}
	}
}

func GetConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		configFile, err := os.ReadFile(os.Getenv("HELM_AGGREGATOR_CONFIG"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Can't get config from filesystem: %v", err), http.StatusInternalServerError)

			return
		}

		_, err = w.Write([]byte(configFile))
		if err != nil {
			http.Error(w, fmt.Sprintf("Can't get config file: %v", err), http.StatusInternalServerError)

			return
		}
	}
}

func Healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
