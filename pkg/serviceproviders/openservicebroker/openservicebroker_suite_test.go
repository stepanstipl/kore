package openservicebroker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOpenservicebroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Openservicebroker Suite")
}
