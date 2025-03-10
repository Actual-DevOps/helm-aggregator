package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"github.com/Actual-DevOps/helm-aggregator/internal/indexer"
	"gopkg.in/yaml.v2"
)

func IndexHandler(config conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
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
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		configFile, err := os.ReadFile(os.Getenv(conf.HelmAggregatorConfig))
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
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("OK")); err != nil {
			http.Error(w, fmt.Sprintf("Can't write healthcheck response: %v", err), http.StatusInternalServerError)
		}
	}
}
