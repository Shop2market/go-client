package product_statistic_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestProductStatistic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProductStatistic Suite")
}
