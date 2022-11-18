package tfcommons

// ServiceDiscoveryResp represents a response from the Terraform service discovery protocol.
type ServiceDiscoveryResp struct {
	LoginV1     *LoginService `json:"login.v1,omitempty"`
	ModulesV1   *string       `json:"modules.v1,omitempty"`
	ProvidersV1 *string       `json:"providers.v1,omitempty"`
}

// ServiceDiscoveryFensakExtResp represents a response from the Terraform Fensak Extensions service discovery protocol.
type ServiceDiscoveryFensakExtResp struct {
	ProvenanceV0 *string `json:"provenance.v0,omitempty"`
}

// LoginService contains the oauth client information that the Terraform CLI can use to interact with the login protocol
// of the registry.
type LoginService struct {
	Client     string   `json:"client"`
	GrantTypes []string `json:"grant_types"`
	Authz      string   `json:"authz"`
	Token      string   `json:"token"`
	Ports      []int    `json:"ports"`
}

// The following three structs represents a response from the module versions endpoint of the registry.
// START

type ModuleVersionsResp struct {
	Modules []ModuleVersionList `json:"modules"`
}

type ModuleVersionList struct {
	Versions []ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Version string `json:"version"`
}

// END
