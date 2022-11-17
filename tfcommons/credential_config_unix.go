package tfcommons

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func credentialConfigSearchPaths() []string {
	searchPaths := []string{}

	custom := os.Getenv("TF_CLI_CONFIG_FILE")
	if custom != "" {
		searchPaths = append(searchPaths, custom)
	}

	home, err := homedir.Dir()
	if err != nil {
		return searchPaths
	}

	return append(searchPaths, filepath.Join(home, ".terraformrc"))
}
