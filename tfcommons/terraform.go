package tfcommons

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FindTerraformModules walks the local folder looking for Terraform modules. A Terraform module is any folder
// that has a .tf file.
//
// excludePaths is a slice of paths relative to the root that should be excluded from the search. For example, if you
// wish to ignore the folder `examples` in a repository root, add the string `examples` in the `excludPaths` slice.
//
// This will return a string slice of relative paths from the root that contains a Terraform module.
func FindTerraformModules(root string, excludePaths []string) ([]string, error) {
	// Use a map for tracking module paths so that it can be deduped
	modulePathTracker := make(map[string]bool)

	// Walk the root path, looking for any folders that have a .tf file.
	err := filepath.Walk(
		root,
		func(name string, info os.FileInfo, inErr error) error {
			// Ignore directories, as we can't tell it's a module unless we see a .tf file.
			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(root, name)
			if err != nil {
				return err
			}
			for _, ignorePath := range excludePaths {
				if strings.HasPrefix(relPath, ignorePath) {
					return nil
				}
			}

			if filepath.Ext(relPath) == ".tf" {
				modulePathTracker[filepath.Dir(relPath)] = true
			}
			return nil
		},
	)

	modulePaths := make([]string, len(modulePathTracker))
	i := 0
	for path := range modulePathTracker {
		modulePaths[i] = path
		i++
	}
	// Sort the paths so that the return value is consistent.
	sort.Strings(modulePaths)
	return modulePaths, err
}
