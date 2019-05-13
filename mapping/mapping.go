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
type Cache struct {
	Data map[string][][]string
	Date *time.Time
}
type Repo struct {
	creds
	Cache
}

func New(endpoint, username, password string) (repo *Repo, err error) {
	if !strings.HasSuffix(endpoint, PATH) {
		err = fmt.Errorf("wrong endpoint: `%s`", endpoint)
		return
	}
	creds := creds{Endpoint: endpoint, Username: username, Password: password}
	Cache := Cache{make(map[string][][]string), nil}
	repo = &Repo{creds, Cache}
	return
}

func (repo *Repo) FindAllMappings() (mappings map[string][][]string, err error) {
	if repo.Cache.IsValid() {
		mappings = repo.Cache.Data
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
	repo.Cache.Update(mappings)

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

func (c *Cache) Update(mappings map[string][][]string) {
	now := time.Now().UTC()
	c.Data = mappings
	c.Date = &now
}

func (c *Cache) IsValid() bool {
	if c.IsEmpty() {
		return false
	}
	now := time.Now().UTC()
	expiration := c.Date.Add(CACHE_TTL)
	return expiration.After(now)
}

func (c *Cache) IsEmpty() bool {
	return c.Data == nil || c.Date == nil
}
