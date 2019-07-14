package main

import (
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/gohelp/winapp"
	"github.com/fpawel/mil82/internal/api"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/delphirpc"
	"os"
	"path/filepath"
	r "reflect"
)

func main() {
	dir := filepath.Join(os.Getenv("DELPHIPATH"),
		"src", "github.com", "fpawel", "mil82gui", "api")
	winapp.EnsuredDirectory(dir)

	createFile := func(fileName string) *os.File {
		return must.Create(filepath.Join(dir, fileName))
	}

	servicesSrc := delphirpc.NewServicesSrc("services", "server_data_types", []r.Type{
		r.TypeOf((*api.LastPartySvc)(nil)),
		r.TypeOf((*api.ConfigSvc)(nil)),
		r.TypeOf((*api.RunnerSvc)(nil)),
		r.TypeOf((*api.PeerSvc)(nil)),
		r.TypeOf((*api.ChartsSvc)(nil)),
	}, map[string]string{
		"ProductInfo": "Product",
		"WorkInfo":    "JournalWork",
		"EntryInfo":   "JournalEntry",
	})

	notifySvcSrc := delphirpc.NewNotifyServicesSrc("notify_services", servicesSrc.DataTypes, []delphirpc.NotifyServiceType{
		{
			"Panic",
			r.TypeOf((*string)(nil)).Elem(),
		},
		{
			"ReadVar",
			r.TypeOf((*types.AddrVarValue)(nil)).Elem(),
		},
		{
			"AddrError",
			r.TypeOf((*types.AddrError)(nil)).Elem(),
		},
		{
			"WorkStarted",
			r.TypeOf((*string)(nil)).Elem(),
		},
		{
			"WorkComplete",
			r.TypeOf((*types.WorkResultInfo)(nil)).Elem(),
		},
		{
			"Warning",
			r.TypeOf((*string)(nil)).Elem(),
		},
		{
			"Delay",
			r.TypeOf((*types.DelayInfo)(nil)).Elem(),
		},
	})

	file := createFile("services.pas")
	servicesSrc.WriteUnit(file)
	must.Close(file)

	file = createFile("server_data_types.pas")
	servicesSrc.DataTypes.WriteUnit(file)
	must.Close(file)

	file = createFile("notify_services.pas")
	notifySvcSrc.WriteUnit(file)
	must.Close(file)

	dir = filepath.Join(os.Getenv("GOPATH"),
		"src", "github.com", "fpawel", "mil82", "internal", "api", "notify")

	file = must.Create(filepath.Join(dir, "api_generated.go"))
	notifySvcSrc.WriteGoFile(file)
	must.Close(file)

}
