package shop_product_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShopProduct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ShopProduct Suite")
}
