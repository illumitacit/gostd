package tfstd

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// CredentialConfig represents the credentials subcomponent of the Terraform CLI config (.terraformrc file).
type CredentialConfig struct {
	Credentials []Credential `hcl:"credentials,block"`
	// TODO: add support for credential helpers

	Remain hcl.Body `hcl:",remain"`

	// Internal attributes that are computed at load time to optimize downstream usage
	credentialMap map[string]*Credential
}

// Credential represents a credential block in the terraform CLI config, which maps the registry host to the token to
// use for authentication.
type Credential struct {
	RegistryHost string `hcl:",label"`
	Token        string `hcl:"token,attr"`
}

// LoadCredentialConfig will load the credential information from the Terraform CLI config. The config is loaded from
// the location specified in the Terraform documentation as of version 1.3.0:
//
// https://developer.hashicorp.com/terraform/cli/config/config-file#locations
//
// Practically this is the following:
// - TF_CLI_CONFIG_FILE env var.
// - On Windows, search %APPDATA%/terraform.rc
// - On Unix, search $HOME/.terraformrc
func LoadCredentialConfig() (*CredentialConfig, error) {
	path := findCredentialConfigFile()
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(data, path)
	if diags.HasErrors() {
		return nil, diags
	}

	var cfg CredentialConfig
	if diags := gohcl.DecodeBody(f.Body, nil, &cfg); diags.HasErrors() {
		return nil, diags
	}

	cfg.credentialMap = map[string]*Credential{}
	for _, cred := range cfg.Credentials {
		cfg.credentialMap[cred.RegistryHost] = &cred
	}

	return &cfg, nil
}

func findCredentialConfigFile() string {
	searchPaths := credentialConfigSearchPaths()
	for _, p := range searchPaths {
		finfo, err := os.Stat(p)
		if err == nil && !finfo.IsDir() {
			return p
		}
	}
	return ""
}
