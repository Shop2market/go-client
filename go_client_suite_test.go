package go_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoClient Suite")
}
