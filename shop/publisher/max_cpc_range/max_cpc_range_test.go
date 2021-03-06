package max_cpc_range_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/Shop2market/go-client/shop/publisher/max_cpc_range"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaxCpcRange", func() {
	It("unmarshals publisher max cpc ranges", func() {
		maxCpcRanges := []*MaxCpcRange{}
		body, _ := ioutil.ReadFile("fixtures/max_cpc_ranges_response.json")
		json.Unmarshal(body, &maxCpcRanges)
		Expect(maxCpcRanges[0]).To(Equal(&MaxCpcRange{
			Name:              "Group A",
			ChannelCategoryID: 352261,
			MaxCpcMin:         3.0,
			MaxCpcMax:         7.0,
		},
		))
	})
})
