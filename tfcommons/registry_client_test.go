package tfcommons_test

import (
	"fmt"
	"sort"

	tfaddr "github.com/hashicorp/terraform-registry-address"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fensak-io/go-commons/tfcommons"
)

var _ = Describe("RegistryClient", func() {
	Describe("initializing base client", func() {
		Context("against public registry", func() {
			It("should work without error", func() {
				clt, err := tfcommons.NewRegistryClient(publicRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())
				Ω(clt.ModulesEndpoint.String()).Should(Equal(
					fmt.Sprintf("https://%s/v1/modules/", publicRegistryHost),
				))
			})
		})
	})

	Describe("fetching versions", func() {
		Context("for test module", func() {
			It("should fetch expected versions", func() {
				clt, err := tfcommons.NewRegistryClient(publicRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())

				versions, err := clt.GetVersions(tfaddr.ModulePackage{
					Namespace:    "yorinasub17",
					Name:         "terragrunt-registry-test",
					TargetSystem: "null",
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).ShouldNot(BeNil())

				vStrs := make([]string, len(versions.Versions))
				for i, v := range versions.Versions {
					vStrs[i] = v.Version
				}
				sort.Strings(vStrs)
				Ω(vStrs).Should(Equal([]string{"0.0.1", "0.0.2"}))
			})
		})
	})
})

const (
	publicRegistryHost = "registry.terraform.io"
)
