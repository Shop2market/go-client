package mapping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	PATH      = "/api/v1/mapping_files.json"
	CACHE_TTL = 1 * time.Hour
)

type creds struct {
	Endpoint string
	Username string
	Password string
}
type cache struct {
	data map[string][][]string
	date *time.Time
}
type Repo struct {
	creds
	cache
}

func New(endpoint, username, password string) (repo *Repo, err error) {
	if !strings.HasSuffix(endpoint, PATH) {
		err = fmt.Errorf("wrong endpoint: `%s`", endpoint)
		return
	}
	creds := creds{Endpoint: endpoint, Username: username, Password: password}
	cache := cache{make(map[string][][]string), nil}
	repo = &Repo{creds, cache}
	return
}

func (repo *Repo) FindAllMappings() (mappings map[string][][]string, err error) {
	if repo.cache.IsValid() {
		mappings = repo.cache.data
		return
	}
	request, err := repo.prepareRequest()
	if err != nil {
		return
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		err = fmt.Errorf("Responded with error: %s", response.Status)
		return
	}
	err = json.NewDecoder(response.Body).Decode(&mappings)
	if err != nil {
		return
	}
	repo.cache.Update(mappings)

	return
}

func (repo *Repo) Find(name string) (mapping [][]string, err error) {
	mappings, err := repo.FindAllMappings()
	if err != nil {
		return
	}
	mapping, ok := mappings[name]
	if ok {
		return
	}
	err = fmt.Errorf("can't find `%s` mapping", name)
	return
}

func (repo *Repo) prepareRequest() (request *http.Request, err error) {
	request, err = http.NewRequest("GET", repo.Endpoint, nil)
	if err != nil {
		return
	}
	request.SetBasicAuth(repo.Username, repo.Password)
	return
}

func (c *cache) Update(mappings map[string][][]string) {
	now := time.Now().UTC()
	c.data = mappings
	c.date = &now
}

func (c *cache) IsValid() bool {
	if c.IsEmpty() {
		return false
	}
	now := time.Now().UTC()
	expiration := c.date.Add(CACHE_TTL)
	return expiration.After(now)
}

func (c *cache) IsEmpty() bool {
	return c.data == nil || c.date == nil
}
