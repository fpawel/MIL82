package internal

import (
	"os"
	"path/filepath"
)

func DataDir() string {
	if os.Getenv("MIL82_DEV_DB") == "true" {
		return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "fpawel", "mil82", "build")
	}
	return filepath.Dir(os.Args[0])
}
