package api

import (
	"github.com/fpawel/mil82/internal/cfg"
)

type ConfigSvc struct {
}

type AppSettings struct {
	ComportProducts,
	ComportTemperature,
	ComportGas string
	Temperature [3]float32
}

func (_ *ConfigSvc) Vars(_ struct{}, vars *[]cfg.Var) error {
	*vars = cfg.Get().Vars
	return nil
}

func (_ *ConfigSvc) SetPlaceChecked(x struct {
	Place   int
	Checked bool
}, _ *struct{}) error {
	c := cfg.Get()
	c.SetPlaceChecked(x.Place, x.Checked)
	cfg.Set(c)
	cfg.Save()

	return nil
}
