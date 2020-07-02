package jsonutils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJsonutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jsonutils Suite")
}
