package tfcommons

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

const (
	terraformFensakExtSDPath = "/.well-known/terraform_ext_fensak.json"
	provenanceV0Key          = "provenance.v0"
)

// RegistryFensakExtClient is a thin API client for interacting with the Fensak extensions of the Terraform registry
// protocol.
type RegistryFensakExtClient struct {
	ProvenanceEndpoint *url.URL

	httpc *resty.Client
}

// RegistrySupportsFensakExt returns whether the Terraform registry hosted at the given host supports the Fensak
// extensions protocol.
func RegistrySupportsFensakExt(host string) (bool, error) {
	sdURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   terraformFensakExtSDPath,
	}
	r, err := resty.New().R().Get(sdURL.String())
	if err != nil {
		return false, err
	}
	return r.StatusCode() == http.StatusOK, nil
}

// NewRegistryFensakExtClient initializes a new client for accessing the Fensak extension sof the Terraform registry.
// This will use the service discovery protocol to know where to fetch the Fensak specific resources.
// Like RegistryClient, this takes advantage of the CLI configuration (described in
// https://developer.hashicorp.com/terraform/cli/config/config-file#locations) to look for credentials for accessing
// private registries. Refer to the documentation for [LoadCredentialConfig] for information on how the credential
// config is loaded.
func NewRegistryFensakExtClient(host string) (*RegistryFensakExtClient, error) {
	httpc := resty.New()

	creds, err := LoadCredentialConfig()
	if err != nil {
		return nil, err
	}
	if creds != nil {
		hostCred, hasCred := creds.credentialMap[host]
		if hasCred {
			httpc.SetAuthToken(hostCred.Token)
		}
	}
	// TODO: add support for env var creds:
	// https://developer.hashicorp.com/terraform/cli/config/config-file#environment-variable-credentials

	var sdResp ServiceDiscoveryFensakExtResp
	sdURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   terraformFensakExtSDPath,
	}
	if _, err := httpc.R().SetResult(&sdResp).Get(sdURL.String()); err != nil {
		return nil, err
	}

	if sdResp.ProvenanceV0 == nil {
		return nil, fmt.Errorf("no provenance.v0 key in services list")
	}
	provURLStr := *sdResp.ProvenanceV0

	// If the registry returns a path, then reuse host to construct the absolute URL.
	if strings.HasPrefix(provURLStr, "/") {
		provURLStr = fmt.Sprintf("https://%s%s", host, provURLStr)
	}

	provURL, err := url.Parse(provURLStr)
	if err != nil {
		return nil, err
	}

	return &RegistryFensakExtClient{
		ProvenanceEndpoint: provURL,
		httpc:              httpc,
	}, nil
}

// GetAttestation uses the Terraform Fensak Extension registry protocol to retrive the given module provenance.
func (regClt *RegistryFensakExtClient) GetAttestation(
	module tfaddr.ModulePackage, version string,
) (*intoto.ProvenanceStatement, error) {
	attestationURL := *regClt.ProvenanceEndpoint
	attestationURL.Path = path.Join(attestationURL.Path, module.Namespace, module.Name, module.TargetSystem, version, "download")

	var prov intoto.ProvenanceStatement
	r, err := regClt.httpc.R().SetResult(&prov).Get(attestationURL.String())
	if err != nil {
		return nil, err
	}
	if r.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("Error getting attestation for module %s: %s", module, r.Body())
	}
	return &prov, nil
}
