package statistic_test

import (
	"time"

	. "github.com/Shop2market/go-client/statistic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeId", func() {
	It("Generate time_id for the give date", func() {
		time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		Expect(DailyTimeId(time)).To(Equal("20091110"))
	})
})
