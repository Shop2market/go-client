package mapping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Shop2market/go-client/mapping/cache"
)

const PATH = "/api/v1/mapping_files.json"

type creds struct {
	Endpoint string
	Username string
	Password string
}

type Repo struct {
	Creds creds
	Cache *cache.Cache
}

func New(endpoint, username, password string) (repo *Repo, err error) {
	if !strings.HasSuffix(endpoint, PATH) {
		err = fmt.Errorf("wrong endpoint: `%s`", endpoint)
		return
	}
	creds := creds{Endpoint: endpoint, Username: username, Password: password}
	repo = &Repo{creds, cache.New(map[string][][]string{})}
	return
}

func (repo *Repo) FindAllMappings() (mappings map[string][][]string, err error) {
	if repo.Cache.IsValid() {
		mappings, err = repo.Cache.Get()
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
	request, err = http.NewRequest("GET", repo.Creds.Endpoint, nil)
	if err != nil {
		return
	}
	request.SetBasicAuth(repo.Creds.Username, repo.Creds.Password)
	return
}
