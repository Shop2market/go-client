package channel_product_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestChannelProductStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ChannelProductStore Suite")
}
