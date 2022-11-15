package tfcommons_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTfcommons(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tfcommons Suite")
}
