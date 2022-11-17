package tfcommons

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// FindTerraformModules walks the local folder looking for Terraform modules. A Terraform module is any folder
// that has a .tf file.
//
// excludePaths is a slice of path regexes relative to the root that should be excluded from the search. For example, if
// you wish to ignore the folder `examples` in a repository root, add the string `^examples` in the `excludPaths` slice.
// Or if you wish to ignore any .terraform paths: ".terraform".
//
// This will return a string slice of relative paths from the root that contains a Terraform module.
func FindTerraformModules(root string, excludePaths []string) ([]string, error) {
	// Use a map for tracking module paths so that it can be deduped
	modulePathTracker := make(map[string]bool)

	excludePathRegexes := make([]*regexp.Regexp, len(excludePaths))
	for i, path := range excludePaths {
		re, err := regexp.Compile(path)
		if err != nil {
			return nil, err
		}
		excludePathRegexes[i] = re
	}

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
			for _, re := range excludePathRegexes {
				if re.MatchString(relPath) {
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
