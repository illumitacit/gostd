package tfcommons

import (
	"fmt"
	"strings"

	getter "github.com/hashicorp/go-getter"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type TerraformDepType int

const (
	TerraformDepUnknown TerraformDepType = iota
	TerraformDepLocal
	TerraformDepRegistry
	TerraformDepGit
	TerraformDepHG
	TerraformDepS3
	TerraformDepGCS
	TerraformDepHTTP
)

func (typ TerraformDepType) String() string {
	switch typ {
	case TerraformDepLocal:
		return "local"
	case TerraformDepRegistry:
		return "registry"
	case TerraformDepGit:
		return "git"
	case TerraformDepS3:
		return "s3"
	case TerraformDepGCS:
		return "gcs"
	case TerraformDepHTTP:
		return "http"
	default:
		return "unknown"
	}
}

// GetTerraformDepType determines the type of the module dependency based on the source string. Terraform determines how
// to fetch dependencies based on the following logic:
// - If it starts with ./ or ../, then it's a local path reference.
// - If it can be parsed as a module registry address, then it's referencing a module registry.
// - Otherwise, use go-getter.
// Therefore, we implement the same logic here to determine the source address type for the Terraform module.
func GetTerraformDepType(source string) TerraformDepType {
	if isLocalModuleSource(source) {
		return TerraformDepLocal
	}

	if isRegistryModuleSource(source) {
		return TerraformDepRegistry
	}

	switch getGoGetterSourceType(source) {
	case "git":
		return TerraformDepGit
	case "hg":
		return TerraformDepHG
	case "s3":
		return TerraformDepS3
	case "gcs":
		return TerraformDepGCS
	case "http", "https":
		return TerraformDepHTTP
	case "file":
		return TerraformDepLocal
	}

	return TerraformDepUnknown
}

// isLocalModuleSource returns true if the module source is a local path, which always starts with ./ or ../ (or .\ or
// ..\ on Windows).
func isLocalModuleSource(source string) bool {
	sourceSlashed := strings.Replace(source, "\\", "/", -1)
	return strings.HasPrefix(sourceSlashed, "./") || strings.HasPrefix(sourceSlashed, "../")
}

// isRegistryModuleSource returns true if the module source is a Terraform registry based source. We determine this
// based on whether or not the address can be parsed as a registry address.
func isRegistryModuleSource(source string) bool {
	_, err := tfaddr.ParseModuleSource(source)
	return err == nil
}

func getGoGetterSourceType(source string) string {
	detected, err := getter.Detect(source, "", getter.Detectors)
	if err != nil {
		return ""
	}
	for key := range getter.Getters {
		if strings.HasPrefix(detected, fmt.Sprintf("%s::", key)) || strings.HasPrefix(detected, fmt.Sprintf("%s://", key)) {
			return key
		}
	}
	return ""
}
