package taxonomy_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/shop/publisher/taxonomy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Shop/Publisher/Taxonomy", func() {
	It("fetches publisher taxonomy and categories", func() {
		content, err := ioutil.ReadFile("fixtures/shop_publisher_taxonomy.json")
		Expect(err).NotTo(HaveOccurred())
		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyBasicAuth("username", "password"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Username = "username"
		Password = "password"
		Endpoint = server.URL()

		taxonomies, err := Find(&Query{ShopId: 1, PublisherId: 2})
		cpc := 13.0
		Expect(err).NotTo(HaveOccurred())
		Expect(taxonomies).To(HaveLen(3))
		Expect(taxonomies[0].Name).To(Equal("Age group"))
		Expect(taxonomies[0].ID).To(Equal(234))
		Expect(taxonomies[0].Categories).To(HaveLen(10))
		Expect(taxonomies[0].Categories[0].Name).To(Equal("1-3 years"))
		Expect(taxonomies[0].Categories[0].CPC).To(Equal(&cpc))
		Expect(taxonomies[0].Categories[0].ID).To(Equal(335990))
		Expect(taxonomies[0].Categories[0].ParentID).To(Equal(0))
	})
	It("builds paths", func() {
		content, err := ioutil.ReadFile("fixtures/shop_publisher_taxonomy_path.json")
		Expect(err).NotTo(HaveOccurred())
		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyBasicAuth("username", "password"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Username = "username"
		Password = "password"
		Endpoint = server.URL()

		taxonomies, err := Find(&Query{ShopId: 1, PublisherId: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(taxonomies).To(HaveLen(1))
		Expect(taxonomies[0].Name).To(Equal("Energy Label"))
		Expect(taxonomies[0].ID).To(Equal(233))
		Expect(taxonomies[0].Categories).To(HaveLen(8))
		Expect(taxonomies[0].Categories[0].Name).To(Equal("Green"))
		Expect(taxonomies[0].Categories[0].ID).To(Equal(10001))
		Expect(taxonomies[0].Categories[0].ParentID).To(Equal(0))
		Expect(taxonomies[0].Categories[0].Path).To(Equal("Green"))
		Expect(taxonomies[0].Categories[4].Path).To(Equal("Green -> A -> A+"))
	})
	It("handles broken paths", func() {
		content, err := ioutil.ReadFile("fixtures/shop_publisher_taxonomy_path.json")
		Expect(err).NotTo(HaveOccurred())
		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyBasicAuth("username", "password"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Username = "username"
		Password = "password"
		Endpoint = server.URL()

		taxonomies, err := Find(&Query{ShopId: 1, PublisherId: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(taxonomies[0].Categories[1].Path).To(Equal("Black"))
		Expect(taxonomies[0].Categories[2].Path).To(Equal("Black -> Grey"))
		Expect(taxonomies[0].Categories[6].Path).To(Equal("Broken -> G"))
		Expect(taxonomies[0].Categories[7].Path).To(Equal("Broken"))
	})
})
