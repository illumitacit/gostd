package chistd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestChistd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chistd Suite")
}
