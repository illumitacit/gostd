package tfcommons_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fensak-io/go-commons/tfcommons"
)

var _ = Describe("SourceType", func() {
	Describe("determining source types", func() {
		for _, tc := range sourceTypeTestCases {
			tc := tc
			Context(fmt.Sprintf("of %s", tc.source), func() {
				It(fmt.Sprintf("should be %s", tc.expected), func() {
					out := tfcommons.GetTerraformDepType(tc.source)
					Ω(out).To(Equal(tc.expected))
				})
			})
		}
	})

})

var (
	sourceTypeTestCases = []struct {
		source   string
		expected tfcommons.TerraformDepType
	}{
		{
			"./consul",
			tfcommons.TerraformDepLocal,
		},
		{
			"./multi/levels/deep/consul",
			tfcommons.TerraformDepLocal,
		},
		{
			".\\windows\\path",
			tfcommons.TerraformDepLocal,
		},
		{
			"..\\windows\\path\\multilevel",
			tfcommons.TerraformDepLocal,
		},
		{
			"../multi/level",
			tfcommons.TerraformDepLocal,
		},
		{
			"hashicorp/consul/aws",
			tfcommons.TerraformDepRegistry,
		},
		{
			"app.terraform.io/example-corp/k8s-cluster/azurerm",
			tfcommons.TerraformDepRegistry,
		},
		{
			"github.com/yorinasub17/foo",
			tfcommons.TerraformDepGit,
		},
		{
			"github.com/yorinasub17/foo//some/module/dir",
			tfcommons.TerraformDepGit,
		},
		{
			"gitlab.com/yorinasub17/foo",
			tfcommons.TerraformDepGit,
		},
		{
			"gitlab.com/yorinasub17/foo//some/module/dir",
			tfcommons.TerraformDepGit,
		},
		{
			"bucket.s3.amazonaws.com/yorinasub17",
			tfcommons.TerraformDepS3,
		},
		{
			"www.googleapis.com/storage/v1/bucket/yorinasub17",
			tfcommons.TerraformDepGCS,
		},
		{
			"/Users/yorinasub17/terraform/modules",
			tfcommons.TerraformDepLocal,
		},
		{
			"https://some.url.com/yorinasub17/modules",
			tfcommons.TerraformDepHTTP,
		},
		{
			"http://some.url.com/yorinasub17/modules",
			tfcommons.TerraformDepHTTP,
		},
		{
			"https://some.url.com/yorinasub17/modules//*",
			tfcommons.TerraformDepHTTP,
		},
		{
			"git@github.com:yorinasub17/foo.git",
			tfcommons.TerraformDepGit,
		},
		{
			"git@github.com:yorinasub17/foo.git?ref=test-branch",
			tfcommons.TerraformDepGit,
		},
		{
			"git@github.com:yorinasub17/foo.git//bar",
			tfcommons.TerraformDepGit,
		},
		{
			"git@custom.git.com:yorinasub17/foo.git",
			tfcommons.TerraformDepGit,
		},
		{
			"git::ssh://git@github.com:2222/yorinasub17/foo.git",
			tfcommons.TerraformDepGit,
		},
	}
)
