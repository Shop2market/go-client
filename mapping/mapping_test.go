package mapping_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/mapping"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var err error
var _ = Describe("Mapping package", func() {
	Describe("Repo", func() {
		var (
			server   *ghttp.Server
			repo     *Repo
			mappings map[string][][]string
			endpoint string
			username = "test"
			password = "test"
		)

		BeforeEach(func() {
			content, _ := ioutil.ReadFile("fixtures/mappings.json")
			server = ghttp.NewServer()
			endpoint = fmt.Sprintf("%s%s", server.URL(), PATH)
			mappings = map[string][][]string{}
			mapping := [][]string{}
			mapping = append(mapping, []string{"4386", "9002514"})
			mapping = append(mapping, []string{"4916", "9002514"})
			mappings["test"] = mapping

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", PATH),
				ghttp.VerifyBasicAuth(username, password),
				ghttp.RespondWith(http.StatusOK, string(content)),
			))
		})

		Describe("New()", func() {
			It("returns repo", func() {
				username := "TestUser"
				pwd := "TestPwd"
				repo, err := New(endpoint, username, pwd)
				Expect(err).NotTo(HaveOccurred())
				Expect(repo).To(BeAssignableToTypeOf(&Repo{}))
			})
		})

		Describe(".PrepareRequest()", func() {
			It("returns request", func() {
				username := "TestUser"
				pwd := "TestPwd"
				repo, err := New(endpoint, username, pwd)
				Expect(err).NotTo(HaveOccurred())
				request, err := repo.PrepareRequest()
				Expect(err).NotTo(HaveOccurred())
				Expect(request).NotTo(BeNil())
			})
		})

		Describe(".FindAllMappings()", func() {
			It("returns a number of mappings", func() {
				repo, err = New(endpoint, username, password)
				Expect(err).ShouldNot(HaveOccurred())
				gotMappings, err := repo.FindAllMappings()
				Expect(err).NotTo(HaveOccurred())
				expected := [][]string{}
				expected = append(expected, []string{"4386", "9002514"})
				expected = append(expected, []string{"4916", "9002514"})
				Expect(gotMappings["test"]).To(Equal(expected))
			})
		})
		Describe(".Find()", func() {
			Context("when mapping with given name exists", func() {
				It("returns mapping", func() {
					repo, err = New(endpoint, username, password)
					Expect(err).ShouldNot(HaveOccurred())
					mapping, err := repo.Find("test")
					Expect(err).ShouldNot(HaveOccurred())
					expected := [][]string{}
					expected = append(expected, []string{"4386", "9002514"})
					expected = append(expected, []string{"4916", "9002514"})
					Expect(mapping).To(Equal(expected))
				})
			})
			Context("when mapping with given name doesn't exist", func() {
				It("returns error", func() {
					repo, err = New(endpoint, username, password)
					Expect(err).ShouldNot(HaveOccurred())
					mapping, err := repo.Find("wrong")
					Expect(err).Should(HaveOccurred())
					Expect(mapping).To(BeNil())
				})
			})
		})
	})
})
