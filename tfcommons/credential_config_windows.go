package tfcommons

import (
	"os"
	"path/filepath"
)

func credentialConfigSearchPaths() []string {
	searchPaths := []string{}

	custom := os.Getenv("TF_CLI_CONFIG_FILE")
	if custom != "" {
		searchPaths = append(searchPaths, custom)
	}

	appData := os.Getenv("APPDATA")
	return append(searchPaths, filepath.Join(appData, "terraform.rc"))
}
