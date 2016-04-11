package publisher_connection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPublisherConnection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PublisherConnection Suite")
}
