package statistic_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStatisticStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StatisticStore Suite")
}
