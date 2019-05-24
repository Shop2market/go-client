package mapping

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Shop2market/go-client/mapping/cache"
)

const PATH = "/api/v1/mapping_files.json"

type Repo struct {
	Endpoint string
	Username string
	Password string
	Notifier chan string

	Cache *cache.Cache
}

func New(endpoint, username, password string) *Repo {
	return &Repo{endpoint, username, password, make(chan string), cache.New()}
}

func (repo *Repo) FindAllMappings() (mappings map[string][][]string, err error) {
	if repo.Cache.IsValid() {
		mappings, err = repo.Cache.Get()
		if err != nil {
			repo.Notifier <- ".FindAllMappings(): cache not found"
		} else {
			return
		}
	}

	request, err := http.NewRequest("GET", repo.Endpoint, nil)
	if err != nil {
		repo.Notifier <- fmt.Sprintf(".FindAllMappings(): failed: %s | Endpoint: %s", err.Error(), repo.Endpoint)
		return
	}

	request.SetBasicAuth(repo.Username, repo.Password)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		repo.Notifier <- fmt.Sprintf(".FindAllMappings(): failed: %s", err.Error())
		return
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		err = fmt.Errorf("Responded with error: %s", response.Status)
		repo.Notifier <- fmt.Sprintf(".FindAllMappings(): failed: %s", err.Error())
		return
	}
	err = json.NewDecoder(response.Body).Decode(&mappings)
	if err != nil {
		repo.Notifier <- fmt.Sprintf(".FindAllMappings(): failed: %s", err.Error())
		return
	}
	repo.Cache.Update(mappings)
	repo.Notifier <- fmt.Sprintf(".FindAllMappings() successfully executed!")
	return
}

func (repo *Repo) Find(name string) (mapping [][]string, err error) {
	repo.Notifier <- fmt.Sprintf(".Find(\"%s\") called", name)
	mappings, err := repo.FindAllMappings()
	if err != nil {
		return
	}
	mapping, ok := mappings[name]
	if ok {
		repo.Notifier <- fmt.Sprintf("\tmapping \"%s\" found", name)
		return
	}
	err = fmt.Errorf("can't find `%s` mapping", name)
	repo.Notifier <- fmt.Sprintf(".\tFind() failed: %s", err.Error())

	return
}

func PrepareRequest(repo *Repo) (request *http.Request, err error) {
	request, err = http.NewRequest("GET", repo.Endpoint, nil)
	if err != nil {
		return
	}
	request.SetBasicAuth(repo.Username, repo.Password)
	return
}
