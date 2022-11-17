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
	Describe("with public registry", func() {
		Context("initializing base client", func() {
			It("should work without error", func() {
				clt, err := tfcommons.NewRegistryClient(publicRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())
				Ω(clt.ModulesEndpoint.String()).Should(Equal(
					fmt.Sprintf("https://%s/v1/modules/", publicRegistryHost),
				))
			})
		})

		Context("fetching versions", func() {
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

		Context("downloading module", func() {
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

	Describe("with private registry", func() {
		token := os.Getenv("FENSAK_GOSTD_TFC_TEST_TOKEN")

		var originalTFCCfgPath string
		var tmpDir string

		BeforeEach(func() {
			originalTFCCfgPath = os.Getenv("TF_CLI_CONFIG_FILE")

			newTmpDir, err := os.MkdirTemp("", "gostd-test-tfcregclt*")
			Ω(err).ShouldNot(HaveOccurred())
			tmpDir = newTmpDir

			tfrcFPath := filepath.Join(tmpDir, ".terraform.rc")
			tfrcF, err := os.OpenFile(
				tfrcFPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600,
			)
			Ω(err).ShouldNot(HaveOccurred())
			defer tfrcF.Close()

			_, wrErr := tfrcF.WriteString(fmt.Sprintf(tfrcTemplate, token))
			Ω(wrErr).ShouldNot(HaveOccurred())
			synErr := tfrcF.Sync()
			Ω(synErr).ShouldNot(HaveOccurred())

			setErr := os.Setenv("TF_CLI_CONFIG_FILE", tfrcFPath)
			Ω(setErr).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			cleanupErr := os.RemoveAll(tmpDir)
			setEnvErr := os.Setenv("TF_CLI_CONFIG_FILE", originalTFCCfgPath)

			Ω(cleanupErr).ShouldNot(HaveOccurred())
			Ω(setEnvErr).ShouldNot(HaveOccurred())
		})

		Context("initializing base client", func() {
			It("should work without error", func() {
				clt, err := tfcommons.NewRegistryClient(privateRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())
				Ω(clt.ModulesEndpoint.String()).Should(Equal(
					fmt.Sprintf("https://%s/api/registry/v1/modules/", privateRegistryHost),
				))
			})
		})

		Context("fetching versions", func() {
			It("should fetch expected versions", func() {
				clt, err := tfcommons.NewRegistryClient(privateRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())

				versions, err := clt.GetVersions(privateRegTestModule)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).ShouldNot(BeNil())

				vStrs := make([]string, len(versions.Versions))
				for i, v := range versions.Versions {
					vStrs[i] = v.Version
				}
				sort.Strings(vStrs)
				Ω(vStrs).Should(Equal([]string{
					"0.0.1",
					"0.0.2",
					"0.0.3",
					"0.0.4",
					"0.0.5",
					"0.0.6",
					"0.0.6-alpha.1",
				}))
			})
		})

		Context("downloading module", func() {
			It("should successfully download", func() {
				clt, err := tfcommons.NewRegistryClient(privateRegistryHost)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(clt.ModulesEndpoint).ShouldNot(BeNil())

				tmpDir, err := os.MkdirTemp("", "")
				Ω(err).ShouldNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)
				destDir := filepath.Join(tmpDir, "tf")

				dlErr := clt.DownloadToPath(privateRegTestModule, "0.0.5", destDir)
				Ω(dlErr).ShouldNot(HaveOccurred())

				module, diags := tfconfig.LoadModule(filepath.Join(destDir, "moduleone"))
				Ω(diags.HasErrors()).To(BeFalse())
				_, hasTestResource := module.ManagedResources["null_resource.test"]
				Ω(hasTestResource).To(BeTrue())
			})
		})
	})
})

const (
	publicRegistryHost  = "registry.terraform.io"
	privateRegistryHost = "app.terraform.io"
)

var (
	regTestModule = tfaddr.ModulePackage{
		Namespace:    "yorinasub17",
		Name:         "terragrunt-registry-test",
		TargetSystem: "null",
	}
	privateRegTestModule = tfaddr.ModulePackage{
		Host:         privateRegistryHost,
		Namespace:    "fensak-test",
		Name:         "testfensak",
		TargetSystem: "null",
	}
)

const tfrcTemplate = `
credentials "app.terraform.io" {
  token = "%s"
}
`
