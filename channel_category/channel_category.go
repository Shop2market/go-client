package channel_category

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

type Category struct {
	Name       string `json:"name"`
	ParentID   int    `json:"parent_id"`
	ExternalID string `json:"external_id"`
	Path       string `json:"-"`
	ID         int    `json:"id"`
}

func buildPaths(categories *[]*Category) {
	cats := *categories
	for i := range cats {
		paths := buildPath(cats[i], cats)
		cats[i].Path = strings.Join(paths, " -> ")
	}
}

func buildPath(category *Category, categories []*Category) []string {
	if category.ParentID == 0 {
		return []string{category.Name}
	}
	for i := range categories {
		if categories[i].ID == category.ParentID {
			return append(buildPath(categories[i], categories), category.Name)
		}
	}
	return []string{}
}

type Finder func(*Query) ([]*Category, error)

// Use for tests to stub calls to API
func DummyFinder(query *Query) ([]*Category, error) {
	return []*Category{}, nil
}

type Categories []*Category

type CategoriesByID struct{ Categories }

func (s Categories) Len() int               { return len(s) }
func (s Categories) Swap(i, j int)          { s[i], s[j] = s[j], s[i] }
func (s CategoriesByID) Less(i, j int) bool { return s.Categories[i].ID < s.Categories[j].ID }

func Find(query *Query) ([]*Category, error) {
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
	categories := []*Category{}

	if err := json.NewDecoder(response.Body).Decode(&categories); err != nil {
		return nil, err
	}
	buildPaths(&categories)
	sort.Sort(CategoriesByID{categories})
	return categories, nil
}

func apiUrl(query *Query) (string, error) {
	url, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/publisher_categories.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
