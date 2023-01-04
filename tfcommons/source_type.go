package tfcommons

import (
	"fmt"
	"strings"

	gcsgetter "github.com/hashicorp/go-getter/gcs/v2"
	s3getter "github.com/hashicorp/go-getter/s3/v2"
	getter "github.com/hashicorp/go-getter/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type TerraformModSrcType int

const (
	TerraformModSrcUnknown TerraformModSrcType = iota
	TerraformModSrcLocal
	TerraformModSrcRegistry
	TerraformModSrcGit
	TerraformModSrcHG
	TerraformModSrcS3
	TerraformModSrcGCS
	TerraformModSrcHTTP
)

func (typ TerraformModSrcType) String() string {
	switch typ {
	case TerraformModSrcLocal:
		return "local"
	case TerraformModSrcRegistry:
		return "registry"
	case TerraformModSrcGit:
		return "git"
	case TerraformModSrcS3:
		return "s3"
	case TerraformModSrcGCS:
		return "gcs"
	case TerraformModSrcHTTP:
		return "http"
	default:
		return "unknown"
	}
}

// GetTerraformModSrcType determines the type of the module dependency based on the source string. Terraform determines how
// to fetch dependencies based on the following logic:
// - If it starts with ./ or ../, then it's a local path reference.
// - If it can be parsed as a module registry address, then it's referencing a module registry.
// - Otherwise, use go-getter.
// Therefore, we implement the same logic here to determine the source address type for the Terraform module.
func GetTerraformModSrcType(source string) TerraformModSrcType {
	if isLocalModuleSource(source) {
		return TerraformModSrcLocal
	}

	if isRegistryModuleSource(source) {
		return TerraformModSrcRegistry
	}

	g := getGoGetterSourceType(source)
	switch t := g.(type) {
	case *getter.GitGetter:
		return TerraformModSrcGit
	case *getter.HgGetter:
		return TerraformModSrcHG
	case *s3getter.Getter:
		return TerraformModSrcS3
	case *gcsgetter.Getter:
		return TerraformModSrcGCS
	case *getter.HttpGetter:
		return TerraformModSrcHTTP
	case *getter.FileGetter:
		return TerraformModSrcLocal
	default:
		fmt.Printf("WTF!!! %v", t)
	}

	return TerraformModSrcUnknown
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

func getGoGetterSourceType(source string) getter.Getter {
	req := &getter.Request{
		Src:     source,
		Dst:     "",
		GetMode: getter.ModeAny,
	}
	for _, g := range getAllGetters() {
		detected, err := getter.Detect(req, g)
		if err != nil {
			return nil
		}
		if detected {
			return g
		}
	}
	return nil
}

func getAllGetters() []getter.Getter {
	return append(
		[]getter.Getter{
			new(s3getter.Getter),
			new(gcsgetter.Getter),
		},
		getter.Getters...,
	)
}
