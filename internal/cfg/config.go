package cfg

import (
	"encoding/json"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/mil82/internal/data"
	"github.com/powerman/structlog"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	UserAppSettings
	PlacesUncheck []int
	Vars          []Var
}

type Var struct {
	Code int    `db:"var"`
	Name string `db:"name"`
}

type UserAppSettings struct {
	ComportProducts,
	ComportTemperature,
	ComportGas string
	TemperatureMinus,
	TemperaturePlus float64
	BlowGasMinutes,
	BlowAirMinutes,
	HoldTemperatureMinutes int
}

func Set(v Config) {
	mu.Lock()
	defer mu.Unlock()
	must.UnmarshalJSON(must.MarshalJSON(&v), &config)
	return
}

func Get() (result Config) {
	mu.Lock()
	defer mu.Unlock()
	must.UnmarshalJSON(must.MarshalJSON(&config), &result)
	return
}

func (c Config) PlaceChecked(place int) bool {
	for _, x := range c.PlacesUncheck {
		if x == place {
			return false
		}
	}
	return true
}

func (c *Config) SetPlaceChecked(place int, checked bool) {
	if checked {
		b := c.PlacesUncheck[:0]
		for _, x := range c.PlacesUncheck {
			if x != place {
				b = append(b, x)
			}
		}
		c.PlacesUncheck = b
	} else {
		c.PlacesUncheck = append(c.PlacesUncheck, place)
	}
}

func Save() {
	mu.Lock()
	defer mu.Unlock()
	must.WriteFile(fileName(), must.MarshalIndentJSON(&config, "", "    "), 0666)
}

func Open() {
	mu.Lock()
	defer mu.Unlock()

	if err := data.DB.Select(&config.Vars, `SELECT var, name FROM var ORDER BY var`); err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(fileName())
	if err != nil {
		log.PrintErr(err.Error(), "файл", fileName())
	}
	if err == nil {
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.PrintErr(err.Error(), "файл", fileName())
		}
	}
}

func fileName() string {
	return filepath.Join(filepath.Dir(os.Args[0]), "config.json")
}

var (
	config Config

	mu  sync.Mutex
	log = structlog.New()
)
