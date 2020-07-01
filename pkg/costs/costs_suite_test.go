package costs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCosts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Costs Suite")
}
