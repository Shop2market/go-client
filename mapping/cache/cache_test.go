package cache_test

import (
	"time"

	. "github.com/Shop2market/go-client/mapping/cache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var mappings map[string][][]string

	BeforeEach(func() {
		mappings = map[string][][]string{}
		mapping := [][]string{}
		mapping = append(mapping, []string{"4386", "9002514"})
		mapping = append(mapping, []string{"4916", "9002514"})
		mappings["test"] = mapping
	})

	Describe("New()", func() {
		It("returns cache", func() {
			Expect(New(nil)).To(BeAssignableToTypeOf(&Cache{}))
		})
	})

	Describe(".Update()", func() {
		It("updates cached data", func() {

		})
	})

	Describe(".IsEmpty()", func() {
		Context("when cache is empty", func() {
			It("returns true", func() {
				cache := New(nil)
				Expect(cache.IsEmpty()).To(BeTrue())
			})
		})

		Context("when cache is not empty", func() {
			It("returns false", func() {
				cache := New(mappings)
				Expect(cache.IsEmpty()).To(BeFalse())
			})
		})
	})

	Describe(".IsOutdated()", func() {
		Context("when cached data is outdated", func() {
			It("returns true", func() {
				now := time.Now().UTC()
				tObj := &now
				oldDate := tObj.Add(-1*60*60*24*time.Second - 1*time.Second)
				tObj = &oldDate
				cache := &Cache{mappings, tObj}
				Expect(cache.IsOutdated()).To(BeTrue())
			})
		})

		Context("when cached data is not outdated", func() {
			It("returns false", func() {
				now := time.Now().UTC()
				tObj := &now
				cache := &Cache{mappings, tObj}
				Expect(cache.IsOutdated()).To(BeFalse())
			})
		})
	})

	Describe(".IsValid()", func() {
		Context("when cache is empty", func() {
			It("returns false", func() {
				cache := New(nil)
				Expect(cache.IsValid()).To(BeFalse())
			})
		})

		Context("when cached data is outdated", func() {
			It("returns false", func() {
				now := time.Now().UTC()
				tObj := &now
				oldDate := tObj.Add(-1*60*60*24*time.Second - 1*time.Second)
				tObj = &oldDate
				cache := &Cache{mappings, tObj}
				Expect(cache.IsValid()).To(BeFalse())
			})
		})

		Context("when cache contains fresh data", func() {
			It("returns true", func() {
				now := time.Now().UTC()
				tObj := &now
				cache := &Cache{mappings, tObj}
				Expect(cache.IsValid()).To(BeTrue())
			})
		})
	})

	Describe(".Get()", func() {
		Context("when cache is valid", func() {
			It("returns cached data", func() {
				cache := New(mappings)
				data, err := cache.Get()
				Expect(err).NotTo(HaveOccurred())
				Expect(data).To(Equal(mappings))
			})
		})

		Context("when cache is not valid", func() {
			It("returns error", func() {
				now := time.Now().UTC()
				tObj := &now
				oldDate := tObj.Add(-1*60*60*24*time.Second - 1*time.Second)
				tObj = &oldDate
				cache := &Cache{mappings, tObj}
				data, err := cache.Get()
				Expect(err).To(HaveOccurred())
				Expect(data).To(BeNil())
			})
		})
	})
})
