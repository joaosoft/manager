package mgr

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShippingInteractors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interactor Test Suite")
}
