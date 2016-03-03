package channel_category_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestChannelCategory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ChannelCategory Suite")
}
