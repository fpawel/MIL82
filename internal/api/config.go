package api

import "github.com/fpawel/mil82/internal/data"

type ConfigSvc struct {
}

func (_ *ConfigSvc) Vars(_ struct{}, vars *[]string) error {
	return data.DB.Select(vars, `SELECT name FROM var ORDER BY var`)
}
