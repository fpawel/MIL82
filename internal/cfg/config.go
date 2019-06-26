package cfg

import (
	"github.com/pelletier/go-toml"
	"github.com/powerman/must"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	PlacesUncheck map[int]struct{} `toml:"places_uncheck" comment:"номера мест, с которых снята галочка в таблице"`
}

func Get() (result Config) {
	mu.Lock()
	defer mu.Unlock()
	must.UnmarshalJSON(must.MarshalJSON(&config), &result)
	return
}

func Save() {
	mu.Lock()
	defer mu.Unlock()
	data, err := toml.Marshal(&config)
	if err != nil {
		panic(err)
	}
	must.WriteFile(fileName(), data, 0666)
}

func fileName() string {
	return filepath.Join(filepath.Dir(os.Args[0]), "config.toml")
}

func init() {
	mu.Lock()
	defer mu.Unlock()
	data, err := ioutil.ReadFile(fileName())
	if err == nil {
		_ = toml.Unmarshal(data, &config)
	}
}

var (
	config Config
	mu     sync.Mutex
)
