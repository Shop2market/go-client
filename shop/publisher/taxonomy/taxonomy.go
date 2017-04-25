package taxonomy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

var Endpoint string
var Username string
var Password string

type Query struct {
	ShopId      int
	PublisherId int
}

// Category - Categories inside taxonomy
type Category struct {
	Name       string   `json:"name"`
	ParentID   int      `json:"parent_id"`
	ExternalID string   `json:"external_id"`
	Path       string   `json:"-"`
	ID         int      `json:"id"`
	CPC        *float64 `json:"cpc"`
	Keywords   string   `json:"keywords"`
}

// Taxonomy - Root structure to hold Taxonomy type
type Taxonomy struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	IsCategory       bool       `json:"is_category"`
	MappingMandatory bool       `json:"mapping_mandatory"`
	Categories       []Category `json:"categories"`
}

// Finder - Main find functor, can be overloaded for stubs or assigned with package Find function
type Finder func(query *Query) ([]Taxonomy, error)

// Use for tests to stub calls to API
func DummyFinder(query *Query) ([]Taxonomy, error) {
	return []Taxonomy{}, nil
}

type Categories []Category

type CategoriesByID struct{ Categories }

func (s Categories) Len() int               { return len(s) }
func (s Categories) Swap(i, j int)          { s[i], s[j] = s[j], s[i] }
func (s CategoriesByID) Less(i, j int) bool { return s.Categories[i].ID < s.Categories[j].ID }

func buildPaths(categories []Category) {
	cats := categories
	for i := range cats {
		paths := buildPath(cats[i], cats)
		cats[i].Path = strings.Join(paths, "»")
	}
}

func buildPath(category Category, categories []Category) []string {
	if category.ParentID == 0 {
		return []string{category.Name}
	}
	for i := range categories {
		if categories[i].ID == category.ParentID {
			if category.ParentID > category.ID {
				return []string{category.Name}
			}
			return append(buildPath(categories[i], categories), category.Name)
		}
	}
	// at the moment we return the category name if the parent is missing
	return []string{category.Name}
}

func Find(query *Query) ([]Taxonomy, error) {
	url, err := apiUrl(query)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(Username, Password)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	taxonomy := []Taxonomy{}

	if err := json.NewDecoder(response.Body).Decode(&taxonomy); err != nil {
		return nil, err
	}
	for i := range taxonomy {
		sort.Sort(CategoriesByID{taxonomy[i].Categories})
		buildPaths(taxonomy[i].Categories)
	}
	return taxonomy, nil
}
func FindById(taxonomies []Taxonomy, taxonomyId int) (*Category, *Taxonomy) {
	for _, taxonomyObj := range taxonomies {
		i := sort.Search(len(taxonomyObj.Categories), func(i int) bool {
			return taxonomyObj.Categories[i].ID >= taxonomyId
		})
		if i != len(taxonomyObj.Categories) && taxonomyObj.Categories[i].ID == taxonomyId {
			return &taxonomyObj.Categories[i], &taxonomyObj
		}
	}
	return nil, nil
}

func apiUrl(query *Query) (string, error) {
	url, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/publisher_taxonomies.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
