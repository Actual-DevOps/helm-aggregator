package helm

import (
	"fmt"
	"helm.sh/helm/v3/pkg/registry"
)

func aggregateOCIAndHelmRepos() (map[string]interface{}, error) {
	aggregatedIndex := make(map[string]interface{})
	entries := make(map[string]interface{})

	// Агрегация классических Helm-репозиториев
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

	// Агрегация OCI-репозиториев
	ociClient, err := registry.NewClient()
	if err != nil {
		return nil, fmt.Errorf("ошибка создания OCI-клиента: %v", err)
	}

	ociRepos := []string{
		"oci://registry-1.docker.io/bitnamicharts/valkey",
		// Добавьте другие OCI-репозитории
	}

	for _, ociURL := range ociRepos {
		ref, err := registry.ParseReference(ociURL)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга OCI-ссылки: %v", err)
		}

		versions, err := ociClient.Tags(ref)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения версий для %s: %v", ociURL, err)
		}

		chartName := ref.Repository
		entries[chartName] = versions
	}

	aggregatedIndex["apiVersion"] = "v1"
	aggregatedIndex["entries"] = entries
	aggregatedIndex["generated"] = "2023-10-01T00:00:00Z"

	return aggregatedIndex, nil
}
