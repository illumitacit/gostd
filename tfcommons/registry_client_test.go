package tfcommons_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fensak-io/gostd/tfcommons"
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

				versions, err := clt.GetVersions(regTestModule)
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

	Describe("downloading module", func() {
		Context("for test module", func() {
			It("should successfully download", func() {
				clt, err := tfcommons.NewRegistryClient(publicRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())

				tmpDir, err := os.MkdirTemp("", "")
				Ω(err).ShouldNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)
				destDir := filepath.Join(tmpDir, "tf")

				dlErr := clt.DownloadToPath(regTestModule, "0.0.2", destDir)
				Ω(dlErr).ShouldNot(HaveOccurred())

				module, diags := tfconfig.LoadModule(destDir)
				Ω(diags.HasErrors()).To(BeFalse())
				_, hasRootResource := module.ManagedResources["null_resource.root"]
				Ω(hasRootResource).To(BeTrue())
			})
		})
	})
})

const (
	publicRegistryHost = "registry.terraform.io"
)

var (
	regTestModule = tfaddr.ModulePackage{
		Namespace:    "yorinasub17",
		Name:         "terragrunt-registry-test",
		TargetSystem: "null",
	}
)
