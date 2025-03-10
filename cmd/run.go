package cmd

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Actual-DevOps/helm-aggregator/internal/conf"
	"github.com/Actual-DevOps/helm-aggregator/internal/handlers"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Run: func(_ *cobra.Command, _ []string) {
		var config conf.Config
		if err := conf.LoadConfig(&config); err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}

		var wg sync.WaitGroup
		for i := range config.Repos {
			wg.Add(1)
			go func(repo *conf.HelmRepo) {
				defer wg.Done()
				err := repo.LoadIndex()
				if err != nil {
					log.Println(err)
				}
			}(&config.Repos[i])
		}

		wg.Wait()

		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			if _, err := w.Write([]byte("<a href=config>config</a></br><a href=index.yaml>index.yaml</a>")); err != nil {
				http.Error(w, fmt.Sprintf("Can't write healthcheck response: %v", err), http.StatusInternalServerError)
			}
		})
		http.HandleFunc("/healthcheck", handlers.Healthcheck())
		http.HandleFunc("/index.yaml", handlers.IndexHandler(config))
		http.HandleFunc("/config", handlers.GetConfigHandler())

		log.Printf("Server run on port %s\n", config.Port)
		if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
			log.Fatalf("Run error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
