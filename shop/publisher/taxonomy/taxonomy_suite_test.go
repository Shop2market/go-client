package taxonomy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTaxonomy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Taxonomy Suite")
}
