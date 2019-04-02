package mapping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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

const PATH = "/api/v1/mapping_files.json"

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
	if repo.hasCache() {
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
	repo.cache.data = mappings
	now := time.Now().UTC()
	repo.cache.date = &now
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

func (repo *Repo) hasCache() bool {
	return repo.cache.date != nil
}
