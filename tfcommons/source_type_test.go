package tfcommons_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fensak-io/gostd/tfcommons"
)

var _ = Describe("SourceType", func() {
	Describe("determining source types", func() {
		for _, tc := range sourceTypeTestCases {
			tc := tc
			Context(fmt.Sprintf("of %s", tc.source), func() {
				It(fmt.Sprintf("should be %s", tc.expected), func() {
					out := tfcommons.GetTerraformModSrcType(tc.source)
					Ω(out).To(Equal(tc.expected))
				})
			})
		}
	})

})

var (
	sourceTypeTestCases = []struct {
		source   string
		expected tfcommons.TerraformModSrcType
	}{
		{
			"./consul",
			tfcommons.TerraformModSrcLocal,
		},
		{
			"./multi/levels/deep/consul",
			tfcommons.TerraformModSrcLocal,
		},
		{
			".\\windows\\path",
			tfcommons.TerraformModSrcLocal,
		},
		{
			"..\\windows\\path\\multilevel",
			tfcommons.TerraformModSrcLocal,
		},
		{
			"../multi/level",
			tfcommons.TerraformModSrcLocal,
		},
		{
			"hashicorp/consul/aws",
			tfcommons.TerraformModSrcRegistry,
		},
		{
			"app.terraform.io/example-corp/k8s-cluster/azurerm",
			tfcommons.TerraformModSrcRegistry,
		},
		{
			"github.com/yorinasub17/foo",
			tfcommons.TerraformModSrcGit,
		},
		{
			"github.com/yorinasub17/foo//some/module/dir",
			tfcommons.TerraformModSrcGit,
		},
		{
			"gitlab.com/yorinasub17/foo",
			tfcommons.TerraformModSrcGit,
		},
		{
			"gitlab.com/yorinasub17/foo//some/module/dir",
			tfcommons.TerraformModSrcGit,
		},
		{
			"bucket.s3.amazonaws.com/yorinasub17",
			tfcommons.TerraformModSrcS3,
		},
		{
			"www.googleapis.com/storage/v1/bucket/yorinasub17",
			tfcommons.TerraformModSrcGCS,
		},
		{
			"/Users/yorinasub17/terraform/modules",
			tfcommons.TerraformModSrcLocal,
		},
		{
			"https://some.url.com/yorinasub17/modules",
			tfcommons.TerraformModSrcHTTP,
		},
		{
			"http://some.url.com/yorinasub17/modules",
			tfcommons.TerraformModSrcHTTP,
		},
		{
			"https://some.url.com/yorinasub17/modules//*",
			tfcommons.TerraformModSrcHTTP,
		},
		{
			"git@github.com:yorinasub17/foo.git",
			tfcommons.TerraformModSrcGit,
		},
		{
			"git@github.com:yorinasub17/foo.git?ref=test-branch",
			tfcommons.TerraformModSrcGit,
		},
		{
			"git@github.com:yorinasub17/foo.git//bar",
			tfcommons.TerraformModSrcGit,
		},
		{
			"git@custom.git.com:yorinasub17/foo.git",
			tfcommons.TerraformModSrcGit,
		},
		{
			"git::ssh://git@github.com:2222/yorinasub17/foo.git",
			tfcommons.TerraformModSrcGit,
		},
	}
)
