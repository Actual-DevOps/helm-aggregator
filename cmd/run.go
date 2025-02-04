package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
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

		http.HandleFunc("/index.yaml", handlers.IndexHandler(config))

		log.Printf("Server run on port %s\n", config.Port)
		if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
			log.Fatalf("Run error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
