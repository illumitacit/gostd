package tfstd_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fensak-io/gostd/tfstd"
)

var _ = Describe("SourceType", func() {
	Describe("determining source types", func() {
		for _, tc := range sourceTypeTestCases {
			tc := tc
			Context(fmt.Sprintf("of %s", tc.source), func() {
				It(fmt.Sprintf("should be %s", tc.expected), func() {
					out := tfstd.GetTerraformModSrcType(tc.source)
					Ω(out).To(Equal(tc.expected))
				})
			})
		}
	})

})

var (
	sourceTypeTestCases = []struct {
		source   string
		expected tfstd.TerraformModSrcType
	}{
		{
			"./consul",
			tfstd.TerraformModSrcLocal,
		},
		{
			"./multi/levels/deep/consul",
			tfstd.TerraformModSrcLocal,
		},
		{
			".\\windows\\path",
			tfstd.TerraformModSrcLocal,
		},
		{
			"..\\windows\\path\\multilevel",
			tfstd.TerraformModSrcLocal,
		},
		{
			"../multi/level",
			tfstd.TerraformModSrcLocal,
		},
		{
			"hashicorp/consul/aws",
			tfstd.TerraformModSrcRegistry,
		},
		{
			"app.terraform.io/example-corp/k8s-cluster/azurerm",
			tfstd.TerraformModSrcRegistry,
		},
		{
			"github.com/yorinasub17/foo",
			tfstd.TerraformModSrcGit,
		},
		{
			"github.com/yorinasub17/foo//some/module/dir",
			tfstd.TerraformModSrcGit,
		},
		{
			"gitlab.com/yorinasub17/foo",
			tfstd.TerraformModSrcGit,
		},
		{
			"gitlab.com/yorinasub17/foo//some/module/dir",
			tfstd.TerraformModSrcGit,
		},
		{
			"bucket.s3.amazonaws.com/yorinasub17",
			tfstd.TerraformModSrcS3,
		},
		{
			"www.googleapis.com/storage/v1/bucket/yorinasub17",
			tfstd.TerraformModSrcGCS,
		},
		{
			"/Users/yorinasub17/terraform/modules",
			tfstd.TerraformModSrcLocal,
		},
		{
			"https://some.url.com/yorinasub17/modules",
			tfstd.TerraformModSrcHTTP,
		},
		{
			"http://some.url.com/yorinasub17/modules",
			tfstd.TerraformModSrcHTTP,
		},
		{
			"https://some.url.com/yorinasub17/modules//*",
			tfstd.TerraformModSrcHTTP,
		},
		{
			"git@github.com:yorinasub17/foo.git",
			tfstd.TerraformModSrcGit,
		},
		{
			"git@github.com:yorinasub17/foo.git?ref=test-branch",
			tfstd.TerraformModSrcGit,
		},
		{
			"git@github.com:yorinasub17/foo.git//bar",
			tfstd.TerraformModSrcGit,
		},
		{
			"git@custom.git.com:yorinasub17/foo.git",
			tfstd.TerraformModSrcGit,
		},
		{
			"git::ssh://git@github.com:2222/yorinasub17/foo.git",
			tfstd.TerraformModSrcGit,
		},
	}
)
