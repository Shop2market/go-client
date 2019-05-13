package cache_test

import (
	. "github.com/Shop2market/go-client/mapping/cache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	Describe("New()", func() {
		It("returns cache", func() {
			Expect(New(nil)).To(BeAssignableToTypeOf(&Cache{}))
		})
	})

	Describe(".Update()", func() {

	})

	Describe(".IsEmpty()", func() {

	})

	Describe(".IsOutdated()", func() {

	})

	Describe(".IsValid()", func() {

	})
})
