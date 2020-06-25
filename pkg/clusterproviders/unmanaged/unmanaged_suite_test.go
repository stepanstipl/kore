package unmanaged_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUnmanaged(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unmanaged Suite")
}
