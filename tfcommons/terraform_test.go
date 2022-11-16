package tfcommons_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fensak-io/gostd/tfcommons"
)

var _ = Describe("Terraform", func() {
	Describe("find terraform modules", func() {
		Context("with this test path", func() {
			// NOTE: the current dir is where the test file is, so this should find all the modules in the fixtures
			// folder.
			testFilePath, err := filepath.Abs(".")
			Ω(err).NotTo(HaveOccurred())

			Context("without exclude paths", func() {
				It("should find all modules", func() {
					modulePaths, err := tfcommons.FindTerraformModules(testFilePath, nil)
					Ω(err).NotTo(HaveOccurred())
					Ω(modulePaths).To(Equal([]string{
						"fixtures/flat",
						"fixtures/flat-non-main",
						"fixtures/nested",
						"fixtures/nested/multiple/levels",
					}))
				})
			})

			Context("with exclude paths", func() {
				It("should not find any modules", func() {
					modulePaths, err := tfcommons.FindTerraformModules(testFilePath, []string{"fixtures"})
					Ω(err).NotTo(HaveOccurred())
					Ω(modulePaths).To(Equal([]string{}))
				})
			})

		})
	})
})
