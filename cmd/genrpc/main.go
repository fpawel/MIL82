package main

import (
	"github.com/fpawel/gohelp/winapp"
	"github.com/fpawel/mil82/internal/api"
	"github.com/fpawel/mil82/internal/delphirpc"
	"github.com/powerman/must"
	"os"
	"path/filepath"
	r "reflect"
)

func main() {
	types := []r.Type{
		r.TypeOf((*api.LastPartySvc)(nil)),
		r.TypeOf((*api.ConfigSvc)(nil)),
	}
	m := map[string]string{
		"ProductInfo": "Product",
		"WorkInfo":    "JournalWork",
		"EntryInfo":   "JournalEntry",
	}

	dir := filepath.Join(os.Getenv("DELPHIPATH"),
		"src", "github.com", "fpawel", "mil82gui", "api")
	winapp.EnsuredDirectory(dir)

	createFile := func(fileName string) *os.File {
		return must.Create(filepath.Join(dir, fileName))
	}

	servicesSrc := delphirpc.NewServicesSrc("services", "server_data_types", types, m)

	notifySvcSrc := delphirpc.NewNotifyServicesSrc("notify_services", servicesSrc.DataTypes, []delphirpc.NotifyServiceType{
		{
			"Panic",
			r.TypeOf((*string)(nil)).Elem(),
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
