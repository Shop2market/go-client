package publisher_max_cpc_range_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPublisherMaxCpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PublisherMaxCpc Suite")
}
