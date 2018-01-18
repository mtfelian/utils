package gost_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGOST(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GOST Suite")
}
