package mapping_test

import (
	"fmt"

	. "github.com/Shop2market/go-client/mapping"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var err error
var _ = Describe("Mapping package", func() {
	Describe("Repo", func() {
		var (
			server     *ghttp.Server
			repo       *Repo
			statusCode = 200
			mappings   map[string]map[string]string
			endpoint   = "https://demo.shop2market.com"
			username   = "test"
			password   = "test"
		)
		BeforeEach(func() {
			server = ghttp.NewServer()
			mappings = map[string]map[string]string{}
			mappings["first"] = map[string]string{}
			mappings["first"]["100500"] = "1050"
			mappings["first"]["500100"] = "5010"

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", fmt.Sprintf("%s/%s", endpoint, PATH)),
				ghttp.VerifyBasicAuth(username, password),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &mappings),
			))
		})
		Describe(".FindAllMappings()", func() {
			Context("when good creds", func() {
				It("returns a number of mappings", func() {
					repo = New(endpoint, username, password)
					mappings, err = repo.FindAllMappings()
					expected := map[string]string{}
					expected["100500"] = "1050"
					expected["500100"] = "5010"
					Expect(mappings["first"]).To(Equal(expected))
				})
			})
			Context("when wrong creds", func() {
				It("returns error", func() {
					repo = New(endpoint, "wrong", password)
					mappings, err = repo.FindAllMappings()
					Expect(mappings).To(BeEmpty())
					Expect(err.Error()).To(Equal("Responded with error: 404 Not Found"))
				})
			})
		})
		Describe(".Find()", func() {
			Context("when mapping with given name exists", func() {
				It("returns mapping", func() {
					repo = New(endpoint, username, password)
					mapping, err := repo.Find("first")
					Expect(err).ShouldNot(HaveOccurred())
					expected := map[string]string{}
					expected["100500"] = "1050"
					expected["500100"] = "5010"
					Expect(mapping).To(Equal(expected))
				})
			})

			Context("when mapping with given name doesn't exist", func() {
				It("returns error", func() {
					repo = New(endpoint, username, password)
					mapping, err := repo.Find("first")
					Expect(err).Should(HaveOccurred())
					fmt.Printf("\n\n\nHERE: %#v\n\n\n", mapping)
					Expect(mapping).To(BeNil())
				})
			})
		})
	})
})
